package app

import (
	"encoding/hex"
	"log"

	"github.com/cc-api/evidence-api/common/golang/evidence_api"

	"github.com/spf13/cobra"
)

var imrCmd = &cobra.Command{
	Use:   "imr",
	Short: "Retrieve the IMR of the cc vm",
	RunE: func(cmd *cobra.Command, args []string) error {
		sdk, err := GetSDK()
		if err != nil {
			return err
		}
		report, err := sdk.GetCCReport("", "", nil)
		if err != nil {
			return err
		}

		group := report.IMRGroup()
		l := log.Default()
		l.Printf("Measurement Count: %d\n", group.MaxIndex+1)
		alg := evidence_api.GetDefaultTPMAlg()
		for index, digest := range group.Group {
			l.Printf("Index: %v\n", index)
			l.Printf("Algorithms: %v\n", alg)
			l.Printf("HASH: %v\n", hex.EncodeToString(digest.Hash))
		}
		return nil
	},
}
