package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Limpa branches locais já mescladas nas branches principais",
	Run: func(cmd *cobra.Command, args []string) {
		flow.HandlePrune()
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
