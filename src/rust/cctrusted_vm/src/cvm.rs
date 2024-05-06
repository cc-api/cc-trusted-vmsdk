use crate::tdvm::TdxVM;
use anyhow::*;
use cctrusted_base::api_data::CcReport;
use cctrusted_base::cc_type::*;
use cctrusted_base::tcg::EventLogEntry;
use cctrusted_base::tcg::{TcgAlgorithmRegistry, TcgDigest};
use core::result::Result::Ok;
use sha2::{Digest, Sha512};
use std::{env, fs, path::Path};
use tempfile::tempdir_in;

// the interfaces a CVM should implement
pub trait CVM {
    /***
        retrive ConfigFS-TSM report

        Args:
            nonce (String): against replay attacks
            data (String): user data

        Returns:
            the CcReport or error information
    */
    fn process_tsm_report(
        &mut self,
        nonce: Option<String>,
        data: Option<String>,
    ) -> Result<CcReport, anyhow::Error> {
        let (tmp_dir, tmp_str);
        let tsm_report = match env::var("TSM_REPORT") {
            Ok(v) => {
                tmp_str = v;
                let tsm_dir = Path::new(&tmp_str);
                if !tsm_dir.exists() {
                    return Err(anyhow!(
                        "[process_tsm_report] TSM_REPORT is defined but directory does not exist"
                    ));
                }
                tsm_dir
            }
            Err(_) => {
                let tsm_dir = Path::new(TSM_PREFIX);
                if !tsm_dir.exists() {
                    return Err(anyhow!(
                        "[process_tsm_report] TSM is not supported in the current environment"
                    ));
                }

                tmp_dir = tempdir_in(tsm_dir)?;
                tmp_dir.path()
            }
        };

        // Update the hash value if nonce or data exists
        let mut hasher = Sha512::new();
        if nonce.is_some() {
            match base64::decode(nonce.unwrap()) {
                Ok(v) => hasher.update(v),
                Err(e) => return Err(anyhow!("[process_tsm_report] nonce decode failed: {}", e)),
            }
        }
        if data.is_some() {
            match base64::decode(data.unwrap()) {
                Ok(v) => hasher.update(v),
                Err(e) => return Err(anyhow!("[process_tsm_report] data decode failed: {}", e)),
            }
        }

        let inblob: [u8; 64] = hasher
            .finalize()
            .as_slice()
            .try_into()
            .expect("[process_tsm_report] Wrong length of data");

        let pre_generation = fs::read_to_string(tsm_report.join("generation"))
            .expect("[process_tsm_report] generation read failed")
            .trim()
            .parse::<u32>()
            .expect("[process_tsm_report] generation parse failed");

        // Write hash array to inblob
        fs::write(tsm_report.join("inblob"), inblob)
            .expect("[process_tsm_report] Write to inblob failed");
        // Read outblob
        let outblob =
            fs::read(tsm_report.join("outblob")).expect("[process_tsm_report] outblob read failed");
        // Read provider
        let provider = fs::read_to_string(tsm_report.join("provider"))
            .expect("[process_tsm_report] provider read failed");
        // Read auxblob if exists
        let auxblob = match fs::read(tsm_report.join("auxblob")) {
            Ok(v) => Some(v),
            Err(_) => None,
        };
        // Read generation and check the generation
        let generation = fs::read_to_string(tsm_report.join("generation"))
            .expect("[process_tsm_report] generation read failed")
            .trim()
            .parse::<u32>()
            .expect("[process_tsm_report] generation parse failed");
        if generation - pre_generation > 1 {
            return Err(anyhow!("[process_tsm_report] check write race failed"));
        }
        // Convert provider to TeeType
        let cc_type = match provider.as_str() {
            "tdx_guest\n" => TeeType::TDX,
            "sev_guest\n" => TeeType::SEV,
            &_ => todo!(),
        };

        Ok(CcReport {
            cc_report: outblob,
            cc_type,
            cc_report_generation: Some(generation),
            cc_provider: Some(provider),
            cc_aux_blob: auxblob,
        })
    }

    /***
        retrive CVM signed report

        Args:
            nonce (String): against replay attacks
            data (String): user data

        Returns:
            the cc report byte array or error information
    */
    fn process_cc_report(
        &mut self,
        nonce: Option<String>,
        data: Option<String>,
    ) -> Result<CcReport, anyhow::Error>;

    /***
        retrive CVM max number of measurement registers

        Args:
            None

        Returns:
            max index of register of CVM
    */
    fn get_max_index(&self) -> u8;

    /***
        retrive CVM measurement registers, e.g.: RTMRs, vTPM PCRs, etc.

        Args:
            index (u8): the index of measurement register,
            algo_id (u8): the alrogithms ID

        Returns:
            TcgDigest struct
    */
    fn process_cc_measurement(
        &mut self,
        index: u8,
        algo_id: u16,
    ) -> Result<TcgDigest, anyhow::Error>;

    /***
        retrive CVM eventlogs

        Args:
            start and count of eventlogs

        Returns:
            array of eventlogs
    */
    fn process_cc_eventlog(
        &self,
        start: Option<u32>,
        count: Option<u32>,
    ) -> Result<Vec<EventLogEntry>, anyhow::Error>;

    /***
        retrive CVM type

        Args:
            None

        Returns:
            CcType of CVM
    */
    fn get_cc_type(&self) -> CcType;

    //Dump confidential CVM information
    fn dump(&self);
}

// used for return of Boxed trait object in build_cvm()
// this composed trait includes functions in both trait CVM and trait TcgAlgorithmRegistry
pub trait BuildCVM: CVM + TcgAlgorithmRegistry {}

// holds the device node info
pub struct DeviceNode {
    pub device_path: String,
}

/***
 instance a specific  object containers specific CVM methods
 and desired trait functions specified by "dyn BuildCVM"
*/
pub fn build_cvm() -> Result<Box<dyn BuildCVM>, anyhow::Error> {
    // instance a CVM according to detected TEE type
    match get_cvm_type().tee_type {
        TeeType::TDX => Ok(Box::new(TdxVM::new())),
        TeeType::SEV => todo!(),
        TeeType::CCA => todo!(),
        TeeType::TPM => todo!(),
        TeeType::PLAIN => Err(anyhow!("[build_cvm] Error: not in any TEE!")),
    }
}

// detect CVM type
pub fn get_cvm_type() -> CcType {
    let mut tee_type = TeeType::PLAIN;

    if Path::new(TSM_PREFIX).exists() || env::var("TSM_REPORT").is_ok() {
        let provider = match env::var("TSM_REPORT") {
            Ok(v) => {
                let tsm_dir = Path::new(&v);
                fs::read_to_string(tsm_dir.join("provider"))
                    .expect("[get_cvm_type] provider read failed")
            }
            Err(_) => {
                let tsm_dir = Path::new(TSM_PREFIX);
                fs::read_to_string(
                    tempdir_in(tsm_dir)
                        .expect("[get_cvm_type] create temp dir failed")
                        .path()
                        .join("provider"),
                )
                .expect("[get_cvm_type] provider read failed")
            }
        };
        tee_type = match provider.as_str() {
            "tdx_guest\n" => TeeType::TDX,
            "sev_guest\n" => TeeType::SEV,
            &_ => todo!(),
        };
    } else if Path::new(TEE_TPM_PATH).exists() {
        tee_type = TeeType::TPM;
    } else if Path::new(TEE_TDX_1_0_PATH).exists() || Path::new(TEE_TDX_1_5_PATH).exists() {
        tee_type = TeeType::TDX;
    } else if Path::new(TEE_SEV_PATH).exists() {
        tee_type = TeeType::SEV;
    } else {
        // TODO add support for CCA and etc.
    }

    CcType {
        tee_type: tee_type.clone(),
    }
}
