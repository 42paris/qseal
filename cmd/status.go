package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.42paris.fr/froz/qseal/pkg/qseal"
	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Will show the status of all secrets defined in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		qsealRc, err := qsealrc.Load()
		if err != nil {
			cmd.PrintErrln("error loading configuration:", err)
			return
		}

		err = qseal.Status(*qsealRc)
		if err != nil {
			cmd.PrintErrln("error sealing secrets:", err)
			return
		}
	},
}
