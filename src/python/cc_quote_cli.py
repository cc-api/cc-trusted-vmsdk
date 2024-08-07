"""
Command line to get quote
"""
import argparse
import base64
import logging
import os
import random
from evidence_api.api import EvidenceApi
from cctrusted_vm.cvm import ConfidentialVM
from cctrusted_vm.sdk import CCTrustedVmSdk

LOG = logging.getLogger(__name__)

logging.basicConfig(level=logging.NOTSET, format="%(name)s %(levelname)-8s %(message)s")

OUT_FORMAT_RAW = "raw"
OUT_FORMAT_HUMAN = "human"

def out_format_validator(out_format):
    """Validator (callback for ArgumentParser) of output format

    Args:
        out_format: User specified output format.

    Returns:
        Validated value of the argument.

    Raises:
        ValueError: An invalid value is given by user.
    """
    if out_format not in (OUT_FORMAT_HUMAN, OUT_FORMAT_RAW):
        raise ValueError
    return out_format

def make_nounce():
    """Make nonce for demo.

    Returns:
        A nonce for demo that is base64 encoded bytes reprensting a 64 bits unsigned integer.
    """
    # Generte a 64 bits unsigned integer randomly (range from 0 to 64 bits max).
    rand_num = random.randrange(0x0, 0xFFFFFFFFFFFFFFFF, 1)
    nonce = base64.b64encode(rand_num.to_bytes(8, "little"))
    return nonce

def make_userdata():
    """Make userdata for demo.

    Returns:
        User data that is base64 encoded bytes for demo.
    """
    userdata = base64.b64encode(bytes("demo user data", "utf-8"))
    return userdata

def main():
    """Example to call get_cc_report and dump the result to stdout."""
    if ConfidentialVM.detect_cc_type() == EvidenceApi.TYPE_CC_NONE:
        LOG.error("This is not a confidential VM!")
        return
    if os.geteuid() != 0:
        LOG.error("Please run as root which is required for this example!")
        return

    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--out-format",
        action="store",
        default=OUT_FORMAT_RAW,
        dest="out_format",
        help="Output format: raw/human. Default raw.",
        type=out_format_validator
    )
    parser.add_argument(
        "--pcr-selection",
        type=str,
        default="sha256:1,2,10",
        help="PCR Selection to generate quote",
        dest="pcr_selection"
    )
    parser.add_argument("--ak-context", type=str, help="Path to ak context", dest="ak_context")
    args = parser.parse_args()

    nonce = make_nounce()
    LOG.info("demo random number in base64: %s", nonce.decode("utf-8"))
    userdata = make_userdata()
    LOG.info("demo user data in base64: %s", userdata.decode("utf-8"))
    extra_args = {}
    extra_args["pcr_selection"] = args.pcr_selection
    extra_args["ak_context"] = args.ak_context

    if ConfidentialVM.detect_cc_type() == EvidenceApi.TYPE_CC_TPM:
        quote = CCTrustedVmSdk.inst().get_cc_report(nonce, userdata, extra_args)
    else:
        quote = CCTrustedVmSdk.inst().get_cc_report(nonce, userdata)
    if quote is not None:
        quote.dump(args.out_format == OUT_FORMAT_RAW)
    else:
        LOG.error("Fail to get Quote!")
        LOG.error("Please double check the log and your config!")

if __name__ == "__main__":
    main()
