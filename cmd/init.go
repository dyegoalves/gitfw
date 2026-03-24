package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa o repositório com as branches main e develop",
	Run: func(cmd *cobra.Command, args []string) {
		flow.HandleInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
