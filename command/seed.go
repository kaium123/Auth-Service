package command

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run sync",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
