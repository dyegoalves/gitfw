package flow

import (
	"fmt"
	"gitfw/pkg/git"
	"gitfw/pkg/ui"
	"os"
	"os/exec"
	"strings"
)

func HandleInit() {
	ui.LogInfo("📦 Inicializando GitFW e estruturando branches principais...")

	output, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil || strings.TrimSpace(string(output)) != "true" {
		ui.LogInfo("Iniciando novo repositório git...")
		if err := git.Run("git init"); err != nil {
			ui.LogError("Falha ao iniciar repositório git.")
			return
		}
	}

	commitOutput, _ := exec.Command("git", "rev-list", "-n", "1", "--all").Output()
	if len(strings.TrimSpace(string(commitOutput))) == 0 {
		ui.LogInfo("Repositório vazio. Criando commit inicial em 'main'...")
		git.Run("git checkout -b main")
		git.Run("git commit --allow-empty -m 'Initial commit'")
	}

	if !git.BranchExists("main") {
		if err := git.Run("git checkout -b main"); err != nil {
			ui.LogError("Falha ao criar branch main.")
			return
		}
	}

	if !git.BranchExists("develop") {
		if err := git.Run("git checkout -b develop"); err != nil {
			ui.LogError("Falha ao criar branch develop.")
			return
		}
		ui.LogSuccess("Branch 'develop' criado a partir de 'main'.")
	} else {
		ui.LogInfo("Branch 'develop' já existe.")
	}

	git.Run("git checkout develop")
	ui.LogSuccess("Estrutura GitFW configurada com sucesso!")
}

func HandleStatus() {
	ui.LogInfo("📊 Status do fluxo GitFW:")
	branchOutput, _ := exec.Command("git", "branch", "--show-current").Output()
	current := strings.TrimSpace(string(branchOutput))

	fmt.Printf("Branch atual: %s%s%s\n", ui.ColorCyan, current, ui.ColorReset)

	if git.IsRepoClean() {
		fmt.Printf("Estado: %sLimpo%s %s\n", ui.ColorGreen, ui.ColorReset, ui.SymSuccess)
	} else {
		fmt.Printf("Estado: %sCom alterações não commitadas%s %s\n", ui.ColorYellow, ui.ColorReset, ui.SymInfo)
	}

	if git.RemoteExists("origin") {
		fmt.Printf("Remoto: %sConfigurado (origin)%s\n", ui.ColorGreen, ui.ColorReset)
	} else {
		fmt.Printf("Remoto: %sNenhum remoto 'origin' configurado%s\n", ui.ColorYellow, ui.ColorReset)
	}
}

func HandlePrune() {
	ui.LogInfo("🧹 Iniciando limpeza de branches mescladas...")

	if git.RemoteExists("origin") {
		ui.LogInfo("Sincronizando com o remoto origin...")
		git.Run("git fetch --prune origin")
	}

	// Unir branches mescladas em develop e main
	mergedDevelop, _ := exec.Command("git", "branch", "--merged", "develop").Output()
	mergedMain, _ := exec.Command("git", "branch", "--merged", "main").Output()

	mergedStr := string(mergedDevelop) + "\n" + string(mergedMain)
	lines := strings.Split(mergedStr, "\n")

	uniqueBranches := make(map[string]bool)
	prefixes := []string{"feature/", "bugfix/", "release/", "hotfix/", "support/"}

	count := 0
	for _, line := range lines {
		b := strings.TrimSpace(line)
		b = strings.TrimPrefix(b, "* ")
		if b == "" || b == "main" || b == "develop" {
			continue
		}

		isFlowBranch := false
		for _, p := range prefixes {
			if strings.HasPrefix(b, p) {
				isFlowBranch = true
				break
			}
		}

		if isFlowBranch && !uniqueBranches[b] {
			uniqueBranches[b] = true
			if err := git.Run(fmt.Sprintf("git branch -d %s", b)); err == nil {
				ui.LogSuccess("Removida branch mesclada: %s", b)
				count++
			}
		}
	}

	if count == 0 {
		ui.LogInfo("Nenhuma branch mesclada encontrada para limpeza.")
	} else {
		ui.LogSuccess("Limpeza concluída! %d branches removidas.", count)
	}
}

func HandleFlow(prefix, base string, args []string) {
	if len(args) < 1 {
		ui.LogError("Uso: gitfw %s [start|finish|publish|list|track|pull] [nome]", prefix)
		os.Exit(1)
	}

	action := args[0]
	name := ""
	if len(args) >= 2 {
		name = args[1]
	}

	switch action {
	case "start":
		if name == "" {
			ui.LogError("Nome é obrigatório para 'start'.")
			os.Exit(1)
		}
		ui.LogInfo("Iniciando %s: %s (base: %s)", prefix, name, base)
		if err := git.Run(fmt.Sprintf("git checkout %s", base)); err != nil {
			ui.LogError("Erro ao trocar para branch base %s.", base)
			return
		}
		if git.RemoteExists("origin") {
			git.Run(fmt.Sprintf("git pull origin %s", base))
		} else {
			ui.LogWarning("Remoto 'origin' não encontrado. Pulando pull...")
		}
		if err := git.Run(fmt.Sprintf("git checkout -b %s/%s", prefix, name)); err != nil {
			ui.LogError("Falha ao criar branch %s/%s", prefix, name)
			return
		}
		ui.LogSuccess("%s/%s iniciado com sucesso!", prefix, name)

	case "finish":
		if name == "" {
			ui.LogError("Nome é obrigatório para 'finish'.")
			os.Exit(1)
		}
		if !git.IsRepoClean() {
			ui.LogError("O repositório possui alterações não commitadas. Finalize-as antes de concluir o fluxo.")
			return
		}
		if !git.BranchExists(prefix + "/" + name) {
			ui.LogError("A branch %s/%s não foi encontrada localmente.", prefix, name)
			return
		}
		ui.LogInfo("Finalizando %s: %s", prefix, name)
		if err := git.Run(fmt.Sprintf("git checkout %s", base)); err != nil {
			ui.LogError("Erro ao trocar para branch base %s.", base)
			return
		}
		if git.RemoteExists("origin") {
			git.Run(fmt.Sprintf("git pull origin %s", base))
		} else {
			ui.LogWarning("Remoto 'origin' não encontrado. Pulando pull...")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff %s/%s", prefix, name)); err != nil {
			ui.LogError("Conflitos detectados durante o merge. Resolva-os manualmente.")
			return
		}
		if err := git.Run(fmt.Sprintf("git branch -d %s/%s", prefix, name)); err != nil {
			ui.LogWarning("Não foi possível excluir a branch local %s/%s. Exclua manualmente se necessário.", prefix, name)
		}
		ui.LogSuccess("%s/%s finalizado e mesclado em %s.", prefix, name, base)

	case "publish":
		if name == "" {
			ui.LogError("Nome é obrigatório para 'publish'.")
			os.Exit(1)
		}
		ui.LogInfo("Publicando %s: %s no origin", prefix, name)
		if err := git.Run(fmt.Sprintf("git push origin %s/%s", prefix, name)); err != nil {
			ui.LogError("Falha ao fazer push para o servidor.")
			return
		}
		ui.LogSuccess("%s/%s publicado!", prefix, name)

	case "list":
		ui.LogInfo("Listando branches de %s:", prefix)
		cmd := exec.Command("git", "branch", "--list", prefix+"/*")
		cmd.Stdout = os.Stdout
		cmd.Run()

	case "track":
		if name == "" {
			ui.LogError("Nome é obrigatório para 'track'.")
			os.Exit(1)
		}
		ui.LogInfo("Rastreando %s remoto: %s", prefix, name)
		if err := git.Run(fmt.Sprintf("git checkout -b %s/%s origin/%s/%s", prefix, name, prefix, name)); err != nil {
			ui.LogError("Falha ao rastrear branch remota.")
			return
		}
		ui.LogSuccess("Rastreando %s/%s com sucesso!", prefix, name)

	case "pull":
		ui.LogInfo("Atualizando branches de %s do servidor...", prefix)
		if git.RemoteExists("origin") {
			git.Run("git fetch origin")
			ui.LogSuccess("Fetch concluído.")
		} else {
			ui.LogError("Remoto 'origin' não configurado.")
		}

	default:
		ui.LogError("Ação desconhecida: %s", action)
	}
}

func HandleBugfix(args []string) {
	if len(args) < 1 {
		ui.LogError("Uso: gitfw bugfix [start|finish|publish|list|track|pull] [nome]")
		os.Exit(1)
	}

	action := args[0]
	name := ""
	if len(args) >= 2 {
		name = args[1]
	}

	switch action {
	case "start":
		if name == "" {
			ui.LogError("Nome do bug é obrigatório para 'start'.")
			os.Exit(1)
		}
		ui.LogInfo("🔧 Corrigindo bugfix: %s (base: develop)", name)
		if err := git.Run("git checkout develop"); err != nil {
			ui.LogError("Erro ao trocar para branch develop.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin develop")
		}
		if err := git.Run(fmt.Sprintf("git checkout -b bugfix/%s", name)); err != nil {
			ui.LogError("Falha ao criar branch bugfix/%s", name)
			return
		}
		ui.LogSuccess("Branch bugfix/%s criado com sucesso!", name)

	case "finish":
		if name == "" {
			ui.LogError("Nome do bug é obrigatório para 'finish'.")
			os.Exit(1)
		}
		if !git.IsRepoClean() {
			ui.LogError("O repositório possui alterações não commitadas. Finalize-as antes de concluir o bugfix.")
			return
		}
		if !git.BranchExists("bugfix/" + name) {
			ui.LogError("A branch bugfix/%s não foi encontrada.", name)
			return
		}

		ui.LogInfo("🏁 Concluindo correção de bugfix: %s", name)
		if err := git.Run("git checkout develop"); err != nil {
			ui.LogError("Erro ao trocar para branch develop.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin develop")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff bugfix/%s", name)); err != nil {
			ui.LogError("Conflitos detectados durante o merge do bugfix.")
			return
		}
		git.Run(fmt.Sprintf("git branch -d bugfix/%s", name))
		ui.LogSuccess("Bugfix %s mesclado em develop e branch removida.", name)

	default:
		HandleFlow("bugfix", "develop", args)
	}
}

func HandleRelease(args []string) {
	if len(args) < 2 {
		ui.LogError("Uso: gitfw release [start|finish|publish|list|track] [versão]")
		os.Exit(1)
	}

	action := args[0]
	version := args[1]

	if action == "finish" {
		if !git.IsRepoClean() {
			ui.LogError("O repositório possui alterações não commitadas. Finalize-as antes de concluir a release.")
			return
		}
		if !git.BranchExists("release/" + version) {
			ui.LogError("A branch release/%s não foi encontrada localmente.", version)
			return
		}
		ui.LogInfo("Finalizando release v%s", version)

		if err := git.Run("git checkout main"); err != nil {
			ui.LogError("Erro ao trocar para branch main.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin main")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff release/%s", version)); err != nil {
			ui.LogError("Erro ao mesclar release em main.")
			return
		}
		git.Run(fmt.Sprintf("git tag -a v%s -m 'Release %s'", version, version))

		if err := git.Run("git checkout develop"); err != nil {
			ui.LogError("Erro ao trocar para branch develop.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin develop")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff release/%s", version)); err != nil {
			ui.LogError("Erro ao mesclar release de volta em develop.")
			return
		}

		git.Run(fmt.Sprintf("git branch -d release/%s", version))
		ui.LogSuccess("Release v%s finalizada, tagueada e mesclada!", version)
		ui.LogInfo("Lembre-se de fazer 'git push origin main develop --tags'")
	} else {
		HandleFlow("release", "develop", args)
	}
}

func HandleHotfix(args []string) {
	if len(args) < 2 {
		ui.LogError("Uso: gitfw hotfix [start|finish|publish|list|track] [nome]")
		os.Exit(1)
	}

	action := args[0]
	name := args[1]

	if action == "finish" {
		if !git.IsRepoClean() {
			ui.LogError("O repositório possui alterações não commitadas. Finalize-as antes de concluir o hotfix.")
			return
		}
		if !git.BranchExists("hotfix/" + name) {
			ui.LogError("A branch hotfix/%s não foi encontrada localmente.", name)
			return
		}
		ui.LogInfo("Finalizando hotfix: %s", name)

		if err := git.Run("git checkout main"); err != nil {
			ui.LogError("Erro ao trocar para branch main.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin main")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff hotfix/%s", name)); err != nil {
			ui.LogError("Erro ao mesclar hotfix em main.")
			return
		}
		git.Run(fmt.Sprintf("git tag -a v%s-hotfix -m 'Hotfix %s'", name, name))

		if err := git.Run("git checkout develop"); err != nil {
			ui.LogError("Erro ao trocar para branch develop.")
			return
		}
		if git.RemoteExists("origin") {
			git.Run("git pull origin develop")
		}
		if err := git.Run(fmt.Sprintf("git merge --no-ff hotfix/%s", name)); err != nil {
			ui.LogError("Erro ao mesclar hotfix de volta em develop.")
			return
		}

		git.Run(fmt.Sprintf("git branch -d hotfix/%s", name))
		ui.LogSuccess("Hotfix %s finalizado, tagueada e mesclado!", name)
		ui.LogInfo("Lembre-se de fazer 'git push origin main develop --tags'")
	} else {
		HandleFlow("hotfix", "main", args)
	}
}
