package cmd

import (
	"gitfw/pkg/flow"
	"strings"
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
	rootCmd.AddCommand(newBugfixCmd())
	rootCmd.AddCommand(newFlowCmd("support", "main"))
}

func newBugfixCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bugfix",
		Short: "Gerencia o fluxo de bugfix (desenvolvimento)",
	}

	subcommands := []struct {
		use   string
		short string
	}{
		{"start [nome]", "Inicia uma correção de bug"},
		{"finish [nome]", "Finaliza e mescla a correção em develop"},
		{"publish [nome]", "Envia o bugfix para o servidor"},
		{"list", "Lista bugfixes ativos"},
		{"track [nome]", "Rastreia um bugfix remoto"},
		{"pull", "Sincroniza bugfixes com o servidor"},
	}

	for _, sc := range subcommands {
		sCmd := sc
		cmd.AddCommand(&cobra.Command{
			Use:   sCmd.use,
			Short: sCmd.short,
			Run: func(cmd *cobra.Command, args []string) {
				flow.HandleBugfix(append([]string{strings.Split(sCmd.use, " ")[0]}, args...))
			},
		})
	}

	return cmd
}
