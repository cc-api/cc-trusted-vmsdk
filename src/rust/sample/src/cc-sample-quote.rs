use evidence_api::api::*;
use evidence_api::api_data::*;
use evidence_api::cc_type::TeeType;
use evidence_api::tdx::quote::TdxQuote;
use cctrusted_vm::sdk::API;

use log::*;
use rand::Rng;

fn main() {
    // set log level
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    /***
     * Note: in real user case, the nonce should come from attestation server
     * side to prevent replay attack and the data should be generate by API caller
     * according to user define spec
     */
    let nonce = base64::encode(rand::thread_rng().gen::<[u8; 32]>());
    let data = base64::encode(rand::thread_rng().gen::<[u8; 32]>());

    // retrieve cc report with API "get_cc_report"
    info!("call Evidence API [get_cc_report] to retrieve cc report!");
    let report = match API::get_cc_report(Some(nonce), Some(data), ExtraArgs {}) {
        Ok(q) => q,
        Err(e) => {
            error!("error getting TDX report: {:?}", e);
            return;
        }
    };

    info!("cc report size {}", report.cc_report.len());

    // dump the cc report with API "dump_cc_report"
    // info!("call Evidence API [dump_cc_report] to dump cc report!");
    // API::dump_cc_report(&report.cc_report);

    // parse the cc report with API "parse_cc_report"
    if report.cc_type == TeeType::TDX {
        let tdx_quote: TdxQuote = match CcReport::parse_cc_report(report.cc_report) {
            Ok(q) => q,
            Err(e) => {
                error!("error parse tdx quote: {:?}", e);
                return;
            }
        };
        info!(
            "version = {}, report_data = {}",
            tdx_quote.header.version,
            base64::encode(tdx_quote.body.report_data)
        );

        // show data of the struct TdxQuoteHeader
        info!("call struct show function to show data of the struct TdxQuoteHeader!");
        tdx_quote.header.show();
    }
}
