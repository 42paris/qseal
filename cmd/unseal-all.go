package cmd

import (
	"github.com/42paris/qseal/pkg/qseal"
	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/spf13/cobra"
)

var unsealAllCmd = &cobra.Command{
	Use:   "unseal-all",
	Short: "Unseal all secrets defined in the config file (not recommended, use `qseal sync` or `qseal` most of the time)",
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
