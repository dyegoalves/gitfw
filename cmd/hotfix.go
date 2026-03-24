package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "Gerencia o fluxo de hotfix",
}

func init() {
	hotfixCmd.AddCommand(&cobra.Command{
		Use:   "start [nome]",
		Short: "Inicia um novo fluxo de hotfix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleHotfix(append([]string{"start"}, args...))
		},
	})

	hotfixCmd.AddCommand(&cobra.Command{
		Use:   "finish [nome]",
		Short: "Finaliza e mescla o fluxo de hotfix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleHotfix(append([]string{"finish"}, args...))
		},
	})

	subcommands := []struct {
		use   string
		short string
	}{
		{"publish", "Envia o hotfix para o servidor remoto"},
		{"list", "Lista hotfixes ativos"},
		{"track", "Rastreia um hotfix remoto"},
		{"pull", "Sincroniza hotfixes com o servidor"},
	}

	for _, sc := range subcommands {
		sCmd := sc
		hotfixCmd.AddCommand(&cobra.Command{
			Use:   sCmd.use,
			Short: sCmd.short,
			Run: func(cmd *cobra.Command, args []string) {
				flow.HandleFlow("hotfix", "main", append([]string{sCmd.use}, args...))
			},
		})
	}

	rootCmd.AddCommand(hotfixCmd)
}
