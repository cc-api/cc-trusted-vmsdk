package app

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"

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
		report.Dump(cctrusted_base.QuoteDumpFormat(FlagFormat))
		return nil
	},
}

func makeNonce() string {
	num := uint64(rand.Int63n(math.MaxInt64))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return base64.StdEncoding.EncodeToString(b)
}

func makeUserData() string {
	b := []byte("demo user data")
	return base64.StdEncoding.EncodeToString(b)
}
