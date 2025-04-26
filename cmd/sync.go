package cmd

import (
	"github.com/42paris/qseal/pkg/qseal"
	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Will sync all secrets defined in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		qsealRc, err := qsealrc.Load()
		if err != nil {
			cmd.PrintErrln("error loading configuration:", err)
			return
		}

		err = qseal.Sync(*qsealRc)
		if err != nil {
			cmd.PrintErrln("error sealing secrets:", err)
			return
		}
		cmd.Println("all secrets synced successfully")
	},
}
