[![VMSDK Python Test](https://github.com/cc-api/cc-trusted-vmsdk/actions/workflows/vmsdk-test-python.yaml/badge.svg)](https://github.com/cc-api/cc-trusted-vmsdk/actions/workflows/vmsdk-test-python.yaml)
[![VMSDK Rust Test](https://github.com/cc-api/cc-trusted-vmsdk/actions/workflows/vmsdk-test-rust.yaml/badge.svg)](https://github.com/cc-api/cc-trusted-vmsdk/actions/workflows/vmsdk-test-rust.yaml)

# cc-trusted-vmsdk


## 1. Overview

The `cc-trusted-vmsdk` is a software development kit (SDK) that provides a set of tools and libraries for building an Intel TDX-compatible Confidential Virtual Machine (CVM) image from an off-the-shelf regular VM image, and provides trusted primitives (measurement, eventlog, quote) of CVM. All below steps are supposed to run in a CVM, such as Intel® TD.
This SDK simplifies the process of creating secure and trusted virtual machines in a cloud computing environment, it offers developers a seamless experience in building secure and reliable applications.


## 2. Features

- Support Attestation through Integrity Measurement Architecture (IMA): Ensure the integrity of Confidential Virtual Machine (CVM) instances through robust attestation mechanisms leveraging Integrity Measurement Architecture (IMA). It provides trusted primitives (measurement, eventlog, quote) of CVM. All below steps are supposed to run in a CVM, such as Intel® TD.

- Support `cloud-init` for seamless initial state setting for CVMs: Utilize `cloud-init` for effortless setup of initial states for Confidential Virtual Machines (CVMs), ensuring a smooth and consistent bootstrapping process.

- Support `Terraform`-alike deployment: Facilitate easy and efficient deployment of Confidential Virtual Machines (CVMs) with support for Terraform-like infrastructure provisioning.

- Support seamless Transformation of Ubuntu and Debian Images into CVM Images: Effortlessly convert regular Ubuntu and Debian images into secure and trusted Confidential Virtual Machine (CVM) images, ensuring compatibility and reliability.

- Support Rust and Python modes
  - Python Mode for Fast and Lightweight Deployment: Leverage the Python mode for quick and lightweight deployment scenarios. Python provides agility and ease of use, making it an ideal choice for rapid application development and deployment.
  - Rust Mode for Enhanced Safety and Reliability: Opt for the Rust mode when prioritizing safety and reliability. Rust's memory safety features and strong emphasis on preventing common programming errors make it a robust choice for building secure and high-performance applications.


## 3. Getting Started

VMSDK is supposed to provide VM image rewrite to CVM image, and provide trusted primitives (measurement, eventlog, quote)
of CVM.
All below steps are supposed to run in a CVM, such as Intel® TD with native CCEL and RTMR as trusted foundation.

### Installation

`VMSDK` package is already available in PyPI. You can install the SDK simply by:

```
$ pip install cctrusted-vm
```

If you would like to run from source code. Try:

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk
$ source setupenv.sh
```

### Run CLI tool

It provides 3 CLI tools for quick usage of Python VMSDK.

- [cc_event_log_cli.py](./src/python/cc_event_log_cli.py): Print event log of CVM.
- [cc_imr_cli.py](./src/python/cc_imr_cli.py): Print algorithm and hash od Integrity Measurement Registers (IMR).
- [cc_quote_cli.py](./src/python/cc_quote_cli.py): Print quote of CVM.


How to run the CLI tool:

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk
$ sudo su
$ source setupenv.sh
$ python3 ./src/python/cc_imr_cli.py
```
_NOTE: The CLI tool needs to run via root user._

Below is example output of `cc_imr_cli.py`.

![](/docs/imr-cli-output.png)

### Run Tests

It provides test cases for Python VMSDK. Run tests with the below commands.

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk
$ sudo su
$ source setupenv.sh
$ python3 -m pytest -v ./src/python/tests/test_sdk.py
```

_NOTE: The tests need to run via root user._


## 4. Run in Google Cloud TDX VM environment with vTPM

Google Cloud Platform (GCP) [TDX Preview](https://cloud.google.com/confidential-computing/confidential-vm/docs/create-a-confidential-vm-instance#intel-tdx) does not support CCEL and RTMR yet, but it supports vTPM.
The SDK will get event log and integrated measurement register from vTPM for GCP TDs.

Refer to [How to create GCP TD](https://github.com/cc-api/confidential-cluster/blob/main/deployment/single_node_gcp.md) to create a GCP TD.

Run the following steps in the GCP TD:

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk
$ sudo su
$ source setupenv.sh

# Get PCRs of vTPM
$ python3 ./src/python/cc_imr_cli.py

# Get vTPM event logs in TCG compliant format
$ python3 ./src/python/cc_event_log_cli.py
```

Extra steps are needed before one trying to get a TPM quote.
User need to generate their AK themselves and save the context someplace on the machine. Sample commands using [tpm2_tools](https://github.com/tpm2-software/tpm2-tools) are listed here:

```
# Generate EK (optional if you already have one)
$ tpm2_createek -c <EK_HANDLE, e.g. 0x8101000A> -G rsa  -u ekpub.pem -f pem

# Generate AK that will be used to sign the TPM quote and save the ak context, public pems, etc.
# User could change the algorithm according to their need.
$ tpm2_createak -C <YOUR_EK_HANDLE> -c <PATH_TO_AK_CTX> -G rsa -g sha256 -s rsassa -u akpub.pem -f pem -n akpub.name
```

After having the ak generated, user could use the command below to generate a TPM quote.

```
# Specify the pcr_selection you would like to include for the quote and the path to the ak context while running the command
$ python3 ./src/python/cc_quote_cli.py --pcr-selection <PCR_SELECTION, e.g. "sha256:1,2,10"> --ak-context <PATH_TO_AK_CTX>
```

- The example output of PCRs (IMR) in a GCP TD as follows:
![](/docs/google_tdx_tpm_dump_imr.png)

- The example output of the TPM event log in a GCP TD as follows:
![](/docs/google_tdx_tpm_dump_eventlog.png)

- The example output of the TPM quote in a GCP TD as follows:
![](/docs/google_tdx_tpm_dump_quote.png)

## 5. License
This project is licensed under the Apache 2.0 License.

## 6. Contact
For any inquiries or support, please open an issue or contact us at [Slack](https://cc-api.slack.com/archives/C070P10A0DR).


