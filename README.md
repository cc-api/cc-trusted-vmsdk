# cc-trusted-vmsdk


## 1. Overview

The `cc-trusted-vmsdk` is a software development kit (SDK) that provides a set of tools and libraries for building an Intel TDX-compatible Confidential Virtual Machine (CVM) image from an off-the-shelf regular VM image. This SDK simplifies the process of creating secure and trusted virtual machines in a cloud computing environment, it offers developers a seamless experience in building secure and reliable applications.


## 2. Features

- Support Attestation through Integrity Measurement Architecture (IMA): Ensure the integrity of Confidential Virtual Machine (CVM) instances through robust attestation mechanisms leveraging Integrity Measurement Architecture (IMA).
  
- Support `cloud-init` for seamless initial state setting for CVMs: Utilize `cloud-init` for effortless setup of initial states for Confidential Virtual Machines (CVMs), ensuring a smooth and consistent bootstrapping process.

- Support `Terraform`-alike deployment: Facilitate easy and efficient deployment of Confidential Virtual Machines (CVMs) with support for Terraform-like infrastructure provisioning.
  
- Support seamless Transformation of Ubuntu and Debian Images into CVM Images: Effortlessly convert regular Ubuntu and Debian images into secure and trusted Confidential Virtual Machine (CVM) images, ensuring compatibility and reliability.

- Support Rust and Python modes
  - Python Mode for Fast and Lightweight Deployment: Leverage the Python mode for quick and lightweight deployment scenarios. Python provides agility and ease of use, making it an ideal choice for rapid application development and deployment.
  - Rust Mode for Enhanced Safety and Reliability: Opt for the Rust mode when prioritizing safety and reliability. Rust's memory safety features and strong emphasis on preventing common programming errors make it a robust choice for building secure and high-performance applications.


## 3. Getting Started

VMSDK is supposed to provide trusted primitives (measurement, eventlog, quote) of CVM.
All below steps are supposed to run in a CVM, such as Intel® TD.

### Installation

`VMSDK` package is already available in PyPI. You can install the SDK simply by:

```
$ pip install cctrusted-vm
```

If you would like to run from source code. Try:

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk/src
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
$ cd cc-trusted-vmsdk/src
$ sudo su
$ source setupenv.sh
$ python3 vmsdk/python/cc_imr_cli.py
```
_NOTE: The CLI tool needs to run via root user._

Below is example output of `cc_imr_cli.py`.

![](/docs/imr-cli-output.png)

### Run Tests

It provides test cases for Python VMSDK. Run tests with the below commands.

```
$ git clone https://github.com/cc-api/cc-trusted-vmsdk.git
$ cd cc-trusted-vmsdk/src
$ sudo su
$ source setupenv.sh
$ python3 -m pip install pytest
$ python3 -m pytest -v ./vmsdk/python/tests/test_sdk.py
```

_NOTE: The tests need to run via root user._

### Test the CVM image 

```
$ ./qemu-test.sh -i /path-to-your-cvm-qcow2/td.qcow2 -k /path-to-your-td-guest-os/vmlinuz -r /dev/vda1 
```


## 4. License
This project is licensed under the Apache 2.0 License.

## 5. Contact
For any inquiries or support, please contact us at XXX.


