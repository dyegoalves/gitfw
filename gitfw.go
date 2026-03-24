package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

// Symbols
const (
	SymInfo    = "ℹ"
	SymSuccess = "✔"
	SymError   = "✘"
	SymWait    = "➜"
	SymPackage = "📦"
)

func logInfo(format string, a ...interface{}) {
	fmt.Printf(ColorBlue+SymWait+" "+format+ColorReset+"\n", a...)
}

func logSuccess(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+SymSuccess+" "+format+ColorReset+"\n", a...)
}

func logError(format string, a ...interface{}) {
	fmt.Printf(ColorRed+SymError+" "+format+ColorReset+"\n", a...)
}

func logWarning(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+SymInfo+" "+format+ColorReset+"\n", a...)
}

func run(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao executar: %s (%v)", command, err)
	}
	return nil
}

func isRepoClean() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) == 0
}

func remoteExists(name string) bool {
	err := exec.Command("git", "remote", "get-url", name).Run()
	return err == nil
}

func checkIfGitInstalled() {
	if err := exec.Command("git", "--version").Run(); err != nil {
		logError("Git não encontrado no sistema. Por favor, instale o Git.")
		os.Exit(1)
	}
}

func main() {
	checkIfGitInstalled()

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	module := os.Args[1]

	switch module {
	case "init":
		handleInit()
	case "status":
		handleStatus()
	case "feature":
		handleFlow("feature", "develop")
	case "bugfix":
		handleFlow("bugfix", "develop")
	case "release":
		handleRelease()
	case "hotfix":
		handleHotfix()
	case "support":
		handleFlow("support", "main")
	case "version", "-v", "--version":
		fmt.Println("GitFW CLI v1.2.2")
	case "help", "-h", "--help":
		showHelp()
	default:
		logError("Módulo desconhecido: %s", module)
		showHelp()
		os.Exit(1)
	}
}

func handleInit() {
	logInfo("📦 Inicializando GitFW e estruturando branches principais...")

	// Check if we are inside a git repo
	output, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil || strings.TrimSpace(string(output)) != "true" {
		logInfo("Iniciando novo repositório git...")
		if err := run("git init"); err != nil {
			logError("Falha ao iniciar repositório git.")
			return
		}
	}

	// Forçar criação do branch main caso não existam commits
	commitOutput, _ := exec.Command("git", "rev-list", "-n", "1", "--all").Output()
	if len(strings.TrimSpace(string(commitOutput))) == 0 {
		logInfo("Repositório vazio. Criando commit inicial em 'main'...")
		run("git checkout -b main")
		run("git commit --allow-empty -m 'Initial commit'")
	}

	// Garantir que main existe
	if !branchExists("main") {
		if err := run("git checkout -b main"); err != nil {
			logError("Falha ao criar branch main.")
			return
		}
	}

	// Check if develop already exists
	if !branchExists("develop") {
		if err := run("git checkout -b develop"); err != nil {
			logError("Falha ao criar branch develop.")
			return
		}
		logSuccess("Branch 'develop' criado a partir de 'main'.")
	} else {
		logInfo("Branch 'develop' já existe.")
	}

	run("git checkout develop")
	logSuccess("Estrutura GitFW configurada com sucesso!")
}

func handleStatus() {
	logInfo("📊 Status do fluxo GitFW:")
	branchOutput, _ := exec.Command("git", "branch", "--show-current").Output()
	current := strings.TrimSpace(string(branchOutput))

	fmt.Printf("Branch atual: %s%s%s\n", ColorCyan, current, ColorReset)

	if isRepoClean() {
		fmt.Printf("Estado: %sLimpo%s %s\n", ColorGreen, ColorReset, SymSuccess)
	} else {
		fmt.Printf("Estado: %sCom alterações não commitadas%s %s\n", ColorYellow, ColorReset, SymInfo)
	}

	if remoteExists("origin") {
		fmt.Printf("Remoto: %sConfigurado (origin)%s\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("Remoto: %sNenhum remoto 'origin' configurado%s\n", ColorYellow, ColorReset)
	}
}

func branchExists(name string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name).Run()
	return err == nil
}

func handleFlow(prefix, base string) {
	if len(os.Args) < 3 {
		logError("Uso: gitfw %s [start|finish|publish|list|track|pull] [nome]", prefix)
		os.Exit(1)
	}

	action := os.Args[2]
	name := ""
	if len(os.Args) >= 4 {
		name = os.Args[3]
	}

	switch action {
	case "start":
		if name == "" {
			logError("Nome é obrigatório para 'start'.")
			os.Exit(1)
		}
		logInfo("Iniciando %s: %s (base: %s)", prefix, name, base)
		if err := run(fmt.Sprintf("git checkout %s", base)); err != nil {
			logError("Erro ao trocar para branch base %s.", base)
			return
		}
		if remoteExists("origin") {
			run(fmt.Sprintf("git pull origin %s", base))
		} else {
			logWarning("Remoto 'origin' não encontrado. Pulando pull...")
		}
		if err := run(fmt.Sprintf("git checkout -b %s/%s", prefix, name)); err != nil {
			logError("Falha ao criar branch %s/%s", prefix, name)
			return
		}
		logSuccess("%s/%s iniciado com sucesso!", prefix, name)

	case "finish":
		if name == "" {
			logError("Nome é obrigatório para 'finish'.")
			os.Exit(1)
		}
		if !isRepoClean() {
			logError("O repositório possui alterações não commitadas. Finalize-as antes de concluir o fluxo.")
			return
		}
		if !branchExists(prefix + "/" + name) {
			logError("A branch %s/%s não foi encontrada localmente.", prefix, name)
			return
		}
		logInfo("Finalizando %s: %s", prefix, name)
		if err := run(fmt.Sprintf("git checkout %s", base)); err != nil {
			logError("Erro ao trocar para branch base %s.", base)
			return
		}
		if remoteExists("origin") {
			run(fmt.Sprintf("git pull origin %s", base))
		} else {
			logWarning("Remoto 'origin' não encontrado. Pulando pull...")
		}
		if err := run(fmt.Sprintf("git merge --no-ff %s/%s", prefix, name)); err != nil {
			logError("Conflitos detectados durante o merge. Resolva-os manualmente.")
			return
		}
		if err := run(fmt.Sprintf("git branch -d %s/%s", prefix, name)); err != nil {
			logWarning("Não foi possível excluir a branch local %s/%s. Exclua manualmente se necessário.", prefix, name)
		}
		logSuccess("%s/%s finalizado e mesclado em %s.", prefix, name, base)

	case "publish":
		if name == "" {
			logError("Nome é obrigatório para 'publish'.")
			os.Exit(1)
		}
		logInfo("Publicando %s: %s no origin", prefix, name)
		if err := run(fmt.Sprintf("git push origin %s/%s", prefix, name)); err != nil {
			logError("Falha ao fazer push para o servidor.")
			return
		}
		logSuccess("%s/%s publicado!", prefix, name)

	case "list":
		logInfo("Listando branches de %s:", prefix)
		cmd := exec.Command("git", "branch", "--list", prefix+"/*")
		cmd.Stdout = os.Stdout
		cmd.Run()

	case "track":
		if name == "" {
			logError("Nome é obrigatório para 'track'.")
			os.Exit(1)
		}
		logInfo("Rastreando %s remoto: %s", prefix, name)
		if err := run(fmt.Sprintf("git checkout -b %s/%s origin/%s/%s", prefix, name, prefix, name)); err != nil {
			logError("Falha ao rastrear branch remota.")
			return
		}
		logSuccess("Rastreando %s/%s com sucesso!", prefix, name)

	case "pull":
		logInfo("Atualizando branches de %s do servidor...", prefix)
		run("git fetch origin")
		logSuccess("Fetch concluído.")

	default:
		logError("Ação desconhecida: %s", action)
	}
}

func handleRelease() {
	if len(os.Args) < 4 {
		logError("Uso: gitfw release [start|finish|publish|list|track] [versão]")
		os.Exit(1)
	}

	action := os.Args[2]
	version := os.Args[3]

	if action == "finish" {
		if !isRepoClean() {
			logError("O repositório possui alterações não commitadas. Finalize-as antes de concluir a release.")
			return
		}
		if !branchExists("release/" + version) {
			logError("A branch release/%s não foi encontrada localmente.", version)
			return
		}
		logInfo("Finalizando release v%s", version)

		// Merge para main
		if err := run("git checkout main"); err != nil {
			logError("Erro ao trocar para branch main.")
			return
		}
		if remoteExists("origin") {
			run("git pull origin main")
		}
		if err := run(fmt.Sprintf("git merge --no-ff release/%s", version)); err != nil {
			logError("Erro ao mesclar release em main.")
			return
		}
		run(fmt.Sprintf("git tag -a v%s -m 'Release %s'", version, version))

		// Merge de volta para develop
		if err := run("git checkout develop"); err != nil {
			logError("Erro ao trocar para branch develop.")
			return
		}
		if remoteExists("origin") {
			run("git pull origin develop")
		}
		if err := run(fmt.Sprintf("git merge --no-ff release/%s", version)); err != nil {
			logError("Erro ao mesclar release de volta em develop.")
			return
		}

		run(fmt.Sprintf("git branch -d release/%s", version))
		logSuccess("Release v%s finalizada, tagueada e mesclada!", version)
		logInfo("Lembre-se de fazer 'git push origin main develop --tags'")
	} else {
		handleFlow("release", "develop")
	}
}

func handleHotfix() {
	if len(os.Args) < 4 {
		logError("Uso: gitfw hotfix [start|finish|publish|list|track] [nome]")
		os.Exit(1)
	}

	action := os.Args[2]
	name := os.Args[3]

	if action == "finish" {
		if !isRepoClean() {
			logError("O repositório possui alterações não commitadas. Finalize-as antes de concluir o hotfix.")
			return
		}
		if !branchExists("hotfix/" + name) {
			logError("A branch hotfix/%s não foi encontrada localmente.", name)
			return
		}
		logInfo("Finalizando hotfix: %s", name)

		// Merge para main
		if err := run("git checkout main"); err != nil {
			logError("Erro ao trocar para branch main.")
			return
		}
		if remoteExists("origin") {
			run("git pull origin main")
		}
		if err := run(fmt.Sprintf("git merge --no-ff hotfix/%s", name)); err != nil {
			logError("Erro ao mesclar hotfix em main.")
			return
		}
		run(fmt.Sprintf("git tag -a v%s-hotfix -m 'Hotfix %s'", name, name))

		// Merge de volta para develop
		if err := run("git checkout develop"); err != nil {
			logError("Erro ao trocar para branch develop.")
			return
		}
		if remoteExists("origin") {
			run("git pull origin develop")
		}
		if err := run(fmt.Sprintf("git merge --no-ff hotfix/%s", name)); err != nil {
			logError("Erro ao mesclar hotfix de volta em develop.")
			return
		}

		run(fmt.Sprintf("git branch -d hotfix/%s", name))
		logSuccess("Hotfix %s finalizado, tagueada e mesclado!", name)
		logInfo("Lembre-se de fazer 'git push origin main develop --tags'")
	} else {
		handleFlow("hotfix", "main")
	}
}

func showHelp() {
	fmt.Println(ColorYellow + "GitFW CLI - Premium Git Flow Implementation" + ColorReset)
	fmt.Println("Uso: gitfw <módulo> <ação> [argumentos]")
	fmt.Println("\n" + ColorCyan + "Módulos disponíveis:" + ColorReset)
	fmt.Printf("   %-9s Inicializa o repositório com main/develop\n", "init")
	fmt.Printf("   %-9s Mostra o status do fluxo atual\n", "status")
	fmt.Printf("   %-9s Gerencia funcionalidades (base: develop)\n", "feature")
	fmt.Printf("   %-9s Gerencia correções de bugs (base: develop)\n", "bugfix")
	fmt.Printf("   %-9s Gerencia versões de lançamento (base: develop)\n", "release")
	fmt.Printf("   %-9s Gerencia correções críticas (base: main)\n", "hotfix")
	fmt.Printf("   %-9s Gerencia branches de suporte (base: main)\n", "support")
	fmt.Printf("   %-9s Mostra a versão do CLI\n", "version")

	fmt.Println("\n" + ColorCyan + "Ações (para feature, bugfix, release, hotfix, support):" + ColorReset)
	fmt.Printf("  %-12s Inicia um novo fluxo\n", "start [nome]")
	fmt.Printf("  %-12s Finaliza e mescla o fluxo\n", "finish [nome]")
	fmt.Printf("  %-12s Envia o fluxo para o servidor remoto\n", "publish [nome]")
	fmt.Printf("  %-12s Lista branches ativos do módulo\n", "list")
	fmt.Printf("  %-12s Rastreia uma branch remota\n", "track [nome]")
	fmt.Printf("  %-12s Sincroniza branches com o servidor\n", "pull")

	fmt.Println("\n" + ColorCyan + "Exemplos:" + ColorReset)
	fmt.Printf("  gitfw feature start billing-system\n")
	fmt.Printf("  gitfw bugfix start checkout-error\n")
	fmt.Printf("  gitfw hotfix finish security-patch\n")
	fmt.Printf("  gitfw support list\n")
}

// comentario aqui
