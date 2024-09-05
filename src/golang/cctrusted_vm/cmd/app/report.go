package app

import (
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/cc-api/evidence-api/common/golang/evidence_api"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Retrieve the cc report and dump it",
	RunE: func(cmd *cobra.Command, args []string) error {
		sdk, err := GetSDK()
		if err != nil {
			return err
		}

		nonce := makeNonce()
		userData := makeUserData()
		report, err := sdk.GetCCReport(nonce, userData, nil)
		if err != nil {
			return err
		}
		report.Dump(evidence_api.QuoteDumpFormat(FlagFormat))
		return nil
	},
}

func makeNonce() []byte {
	num := uint64(rand.Int63n(math.MaxInt64))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

func makeUserData() []byte {
	b := []byte("demo user data")
	return b
}
