package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Gerencia o fluxo de release",
}

func init() {
	releaseCmd.AddCommand(&cobra.Command{
		Use:   "start [versão]",
		Short: "Inicia um novo fluxo de release",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleRelease(append([]string{"start"}, args...))
		},
	})

	releaseCmd.AddCommand(&cobra.Command{
		Use:   "finish [versão]",
		Short: "Finaliza e mescla o fluxo de release",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleRelease(append([]string{"finish"}, args...))
		},
	})

	// Reutilizar outros subcomandos de flow.go? Sim, mas chamando HandleFlow diretamente se compatível
	subcommands := []struct {
		use   string
		short string
	}{
		{"publish", "Envia a release para o servidor remoto"},
		{"list", "Lista releases ativas"},
		{"track", "Rastreia uma release remota"},
		{"pull", "Sincroniza releases com o servidor"},
	}

	for _, sc := range subcommands {
		sCmd := sc // capture variable
		releaseCmd.AddCommand(&cobra.Command{
			Use:   sCmd.use,
			Short: sCmd.short,
			Run: func(cmd *cobra.Command, args []string) {
				flow.HandleFlow("release", "develop", append([]string{sCmd.use}, args...))
			},
		})
	}

	rootCmd.AddCommand(releaseCmd)
}
