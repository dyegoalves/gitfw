package main

import (
	"gitfw/cmd"
	"gitfw/pkg/git"
)

func main() {
	// Verifica se o git está instalado antes de qualquer comando
	git.CheckIfGitInstalled()

	// Delega toda a execução para o framework Cobra
	cmd.Execute()
}
