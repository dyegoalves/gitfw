package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o status atual do fluxo GitFW",
	Run: func(cmd *cobra.Command, args []string) {
		flow.HandleStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
