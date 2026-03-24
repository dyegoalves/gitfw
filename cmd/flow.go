package cmd

import (
	"gitfw/pkg/flow"
	"github.com/spf13/cobra"
)

func newFlowCmd(prefix, base string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   prefix,
		Short: "Gerencia o fluxo de " + prefix,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "start [nome]",
		Short: "Inicia um novo fluxo de " + prefix,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, append([]string{"start"}, args...))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "finish [nome]",
		Short: "Finaliza e mescla o fluxo de " + prefix,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, append([]string{"finish"}, args...))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "publish [nome]",
		Short: "Envia o fluxo de " + prefix + " para o servidor remoto",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, append([]string{"publish"}, args...))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Lista branches ativos de " + prefix,
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, []string{"list"})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "track [nome]",
		Short: "Rastreia uma branch remota de " + prefix,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, append([]string{"track"}, args...))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "pull",
		Short: "Sincroniza branches de " + prefix + " com o servidor",
		Run: func(cmd *cobra.Command, args []string) {
			flow.HandleFlow(prefix, base, []string{"pull"})
		},
	})

	return cmd
}

func init() {
	rootCmd.AddCommand(newFlowCmd("feature", "develop"))
	rootCmd.AddCommand(newFlowCmd("bugfix", "develop"))
	rootCmd.AddCommand(newFlowCmd("support", "main"))
}
