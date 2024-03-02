package app

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vm-tool",
	Short: "A demo tool to use vmsdk of golang ",
	Long: `vm-tool
A sample implementation to use vmsdk to get evidence, such as
* cc quote
* event log
* integrated measuremt registers`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	FlagFormat string
)

const (
	FlagFormatHuman = "human"
	FlagFormatRaw   = "raw"
)

func init() {
	// Persistent Flags
	rootCmd.PersistentFlags().StringVarP(
		&FlagFormat, "out-format", "f", FlagFormatHuman,
		"the format of the output: human/raw")

	// sub-commands
	// 1. report command
	rootCmd.AddCommand(reportCmd)
	// 2. imr command
	rootCmd.AddCommand(imrCmd)
	// 3. eventlog command
	rootCmd.AddCommand(eventLogCmd)

}
