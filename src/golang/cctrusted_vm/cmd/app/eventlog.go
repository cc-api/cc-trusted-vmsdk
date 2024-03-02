package app

import (
	"encoding/hex"
	"log"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"

	"github.com/spf13/cobra"
)

var (
	start *int
	count *int
)

var eventLogCmd = &cobra.Command{
	Use:   "eventlog",
	Short: "Handle event log, with sub commands",
	Long:  `Get event log of boot-time (and run-time if necessary)`,
}

var eventLogDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump the retrieved eventlog",
	RunE: func(cmd *cobra.Command, args []string) error {
		el, err := filterEventLog()
		if err != nil {
			return err
		}

		el.Dump(cctrusted_base.QuoteDumpFormat(FlagFormat))
		return nil
	},
}

var eventLogReplayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Replay the retrieved eventlog, printing the result",
	RunE: func(cmd *cobra.Command, args []string) error {
		sdk, err := GetSDK()
		if err != nil {
			return err
		}

		l := log.Default()
		el, err := filterEventLog()
		if err != nil {
			return err
		}

		replay := sdk.ReplayCCEventLog(el.EventLog())
		// Or direct `replay := el.Replay()`
		for idx, elem := range replay {
			for alg, v := range elem {
				l.Printf("Index: %v\n", idx)
				l.Printf("Algorithms: %v\n", alg)
				l.Printf("HASH: %v\n", hex.EncodeToString(v))
			}
		}
		return nil
	},
}

func filterEventLog() (*cctrusted_base.EventLogger, error) {
	sdk, err := GetSDK()
	if err != nil {
		return nil, err
	}
	el, err := sdk.GetCCEventLog(int32(*start), int32(*count))
	if err != nil {
		return nil, err
	}
	return el, nil
}

func init() {
	start = eventLogCmd.Flags().IntP("start", "s", 0, "the start index of the event log")
	count = eventLogCmd.Flags().IntP("count", "c", 0, "the count of the event log")
	eventLogCmd.MarkFlagsRequiredTogether("start", "count")

	eventLogCmd.AddCommand(eventLogDumpCmd)
	eventLogCmd.AddCommand(eventLogReplayCmd)
}
