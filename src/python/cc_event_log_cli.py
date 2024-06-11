"""
Command line to dump the cc event logs
"""
import logging
import argparse
import os
from cctrusted_base.api import CCTrustedApi
from cctrusted_base.eventlog import TcgEventLog
from cctrusted_base.tcgcel import TcgTpmsCelEvent
from cctrusted_base.tcg import TcgAlgorithmRegistry
from cctrusted_vm.cvm import ConfidentialVM
from cctrusted_vm.sdk import CCTrustedVmSdk


LOG = logging.getLogger(__name__)

logging.basicConfig(level=logging.NOTSET, format='%(name)s %(levelname)-8s %(message)s')

def main():
    """Example cc event log fetching utility."""
    if ConfidentialVM.detect_cc_type() == CCTrustedApi.TYPE_CC_NONE:
        LOG.error("This is not a confidential VM!")
        return
    if os.geteuid() != 0:
        LOG.error("Please run as root which is required for this example!")
        return

    parser = argparse.ArgumentParser(
        description="The example utility to fetch CC event logs")
    parser.add_argument('-s', type=int,
                        help='index of first event log to fetch', dest='start')
    parser.add_argument("-c", type=int, help="number of event logs to fetch",
                        dest="count")
    parser.add_argument("-f", type=str, help="enable canonical tlv format", default="false",
                        dest="cel_format")
    args = parser.parse_args()

    event_logs = CCTrustedVmSdk.inst().get_cc_eventlog(args.start, args.count)
    if event_logs is None:
        LOG.error("No event log fetched. Check debug log for issues.")
        return
    LOG.info("Total %d of event logs fetched.", len(event_logs))

    res = CCTrustedApi.replay_cc_eventlog(event_logs)
    # pylint: disable-next=C0301
    LOG.info("Note: If the underlying platform is TDX, the IMR index showing is cc measurement register instead of TDX measurement register.")
    # pylint: disable-next=C0301
    LOG.info("Please refer to the spec https://www.intel.com/content/www/us/en/content-details/726790/guest-host-communication-interface-ghci-for-intel-trust-domain-extensions-intel-tdx.html")
    LOG.info("Replayed result of collected event logs:")
    # pylint: disable-next=C0201
    for k in sorted(res.keys()):
        LOG.info("IMR[%d]: ", k)
        for alg, h in res.get(k).items():
            LOG.info("   %s: ", TcgAlgorithmRegistry.get_algorithm_string(alg))
            LOG.info("      %s", h.hex())

    LOG.info("Dump collected event logs:")
    for event in event_logs:
        if isinstance(event, TcgTpmsCelEvent):
            if args.cel_format.lower() == 'true':
                TcgTpmsCelEvent.encode(event, TcgEventLog.TCG_FORMAT_CEL_TLV).dump()
            else:
                event.to_pcclient_format().dump()
        else:
            event.dump()

if __name__ == "__main__":
    main()
