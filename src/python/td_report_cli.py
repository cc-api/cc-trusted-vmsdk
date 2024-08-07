
"""
Command line to dump the integrated measurement register
"""
import logging
import os
from evidence_api.api import EvidenceApi
from cctrusted_vm.cvm import ConfidentialVM
from cctrusted_vm.tdx import CCTrustedTdvmSdk

LOG = logging.getLogger(__name__)

logging.basicConfig(level=logging.NOTSET, format="%(name)s %(levelname)-8s %(message)s")

def main():
    """Example to call get_tdreport and dump the result to stdout."""
    if ConfidentialVM.detect_cc_type() != EvidenceApi.TYPE_CC_TDX:
        LOG.error("This is not a TD VM!")
        return
    if os.geteuid() != 0:
        LOG.error("Please run as root which is required for this example!")
        return

    tdreport = CCTrustedTdvmSdk.inst().get_tdreport()
    tdreport.dump()

if __name__ == "__main__":
    main()
