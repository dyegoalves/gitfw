package cmd

import (
	"fmt"
	"gitfw/pkg/ui"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gitfw",
	Short:   "GitFW - Premium Git Flow Implementation",
	Version: "0.1.9-beta",
	Long: ui.ColorYellow + `GitFW CLI - A robust and premium implementation of the Gitflow workflow.
Designed for efficiency, safety, and visual clarity.` + ui.ColorReset,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.SetVersionTemplate("gitfw version {{.Version}}\n")
	rootCmd.Flags().BoolP("version", "v", false, "display version")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Mostra a versão do GitFW CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gitfw version %s\n", rootCmd.Version)
		},
	})
}
