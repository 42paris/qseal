package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.42paris.fr/froz/qseal/pkg/qseal"
	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
)

var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Qseal CLI allow to quickly seal and unseal your k8s secrets with kubeseal",
	Long: `
Qseal CLI allow to quickly seal and unseal your k8s secrets with kubeseal
To know how to seal and unseal your secrets, it refer himself to the qsealrc.yaml file
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
}