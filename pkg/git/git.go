package git

import (
	"fmt"
	"gitfw/pkg/ui"
	"os"
	"os/exec"
	"strings"
)

func Run(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao executar: %s (%v)", command, err)
	}
	return nil
}

func IsRepoClean() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) == 0
}

func RemoteExists(name string) bool {
	err := exec.Command("git", "remote", "get-url", name).Run()
	return err == nil
}

func BranchExists(name string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name).Run()
	return err == nil
}

func CheckIfGitInstalled() {
	if err := exec.Command("git", "--version").Run(); err != nil {
		ui.LogError("Git não encontrado no sistema. Por favor, instale o Git.")
		os.Exit(1)
	}
}
