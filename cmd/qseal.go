package cmd

import (
	"os"

	"github.com/42paris/qseal/pkg/qseal"
	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "qseal is a command-line tool that makes it easy to seal and unseal your Kubernetes secrets using kubeseal.",
	Long: `qseal is a command-line tool that makes it easy to seal and unseal your Kubernetes secrets using kubeseal.
It relies on the qsealrc.yaml file to determine how secrets should be sealed or unsealed.
`,
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

func init() {
	// Add subcommands to the root command
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(sealAllCmd)
	RootCmd.AddCommand(unsealAllCmd)
	RootCmd.AddCommand(syncCmd)
	RootCmd.AddCommand(statusCmd)

	initCmd.PersistentFlags().BoolP("ignore-parents", "i", false, "ignore existing qsealrc.yaml files in parent directories")
}
