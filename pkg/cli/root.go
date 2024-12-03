package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "delta",
	Short: "delta is a key-value store",
	Long:  "",
}
