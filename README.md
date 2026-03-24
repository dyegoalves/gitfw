# GitFW CLI 🚀 (v0.1.8-beta)

**GitFW** é uma ferramenta de linha de comando premium e robusta para implementação da metodologia Git Flow. Agora utilizando o framework **Cobra CLI**, oferece autocompletação profissional, performance e segurança para gerenciar seu ciclo de vida de software.

---

## ✨ Funcionalidades Principais

-   📦 **`init`**: Configura automaticamente as branches `main` e `develop`.
-   📊 **`status`**: Visualização rápida do estado do seu fluxo atual.
-   🧹 **`prune`**: Limpa branches locais já mescladas nas branches principais.
-   ⌨️ **Autocomplete**: Suporte nativo para Bash, Zsh e Fish.
-   🛡️ **Segurança**: Bloqueia finalizações se houver alterações não commitadas.
-   🚀 **Release Automatizada**: Merges duplos e criação de tags via CLI.

---

## 🛠️ Instalação (macOS)

1. Clone o repositório.
2. Execute a instalação:
   ```bash
   make install
   ```

### ⚡ Habilitar Autocomplete (Zsh/Mac)
Para ativar as sugestões de comandos via TAB, adicione ao seu `.zshrc`:
```bash
source <(gitfw completion zsh)
```

---

## 📖 Comandos Comuns

| Comando | Descrição |
| :--- | :--- |
| `gitfw feature start [nome]` | Inicia uma nova funcionalidade |
| `gitfw feature finish [nome]` | Finaliza e mescla em develop |
| `gitfw release start [v]` | Inicia preparação de versão |
| `gitfw release finish [v]` | Finaliza v, tagueia e mescla em main/develop |
| `gitfw bugfix start [nome]` | Inicia correção de bug durante desenvolvimento |
| `gitfw bugfix finish [nome]`| Finaliza e mescla correção em develop |
| `gitfw hotfix start [nome]` | Inicia correção crítica em produção |
| `gitfw prune` | Limpa branches locais obsoletas |
| `gitfw version` / `-v` | Exibe a versão atual |

---

## 🔧 Estrutura Gitflow
-   `main`: Produção (apenas código tagueado).
-   `develop`: Branch oficial de integração.
-   `feature/*`: Novas funcionalidades.
-   `bugfix/*`: Correções de bugs durante o desenvolvimento.
-   `release/*`: Preparação de lançamentos.
-   `hotfix/*`: Correções urgentes em produção.

---

**Desenvolvido com foco em Clean Code e Robustez.** 🛠️
