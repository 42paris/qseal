package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.42paris.fr/froz/qseal/pkg/qseal"
	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
)

var unsealAllCmd = &cobra.Command{
	Use:   "unseal-all",
	Short: "Unseal all secrets defined in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		qsealRc, err := qsealrc.Load()
		if err != nil {
			cmd.PrintErrln("error loading configuration:", err)
			return
		}

		err = qseal.UnsealAll(*qsealRc)
		if err != nil {
			cmd.PrintErrln("error sealing secrets:", err)
			return
		}
		cmd.Println("all secrets unsealed successfully")
	},
}