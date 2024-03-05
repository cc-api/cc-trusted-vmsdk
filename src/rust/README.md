# Rust SDK for CC Trusted API in Confidential VM

This is the Rust version of our VM SDK to help you using the CC Trusted API in your Rust programs. The sub folder "cctrusted_vm" include all the source code for the VMSDK. The sub folder "sample" includes some commandline examples for your reference.

# Run CLI Samples

We can try the CLI samples like this:

```bash
cd sample

# get measurement
cargo run --bin cc-sample-measurement

# get event log
cargo run --bin cc-sample-eventlog

# get quote
cargo run --bin cc-sample-quote
```

Or, after build successfully, we can also run the CLIs directly:

```bash
cd sample

# build the release version
cargo build --release

# get measurement
target/release/cc-sample-measurement

# get event log
target/release/cc-sample-eventlog

# get quote
target/release/cc-sample-quote
```

# Run Test

The test is simple:

```bash
cd cctrusted_vm
cargo test
```
