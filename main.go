package main

import (
	"fmt"
	"os"

	"gitlab.42paris.fr/froz/qseal/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
