package cmd

import (
	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the qseal configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		ignoreParents, err := cmd.PersistentFlags().GetBool("ignore-parents")
		if err != nil {
			cmd.PrintErrln("error reading command arguments:", err)
			return
		}

		err = qsealrc.Init(ignoreParents)
		if err != nil {
			cmd.PrintErrln("error initializing configuration:", err)
			return
		}
	},
}
