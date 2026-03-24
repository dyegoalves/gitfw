# GitFW CLI 🚀

**GitFW** é uma ferramenta de linha de comando premium e robusta para implementação da metodologia Git Flow. Desenvolvida em Go, ela oferece performance, segurança e uma experiência de usuário simplificada para gerenciar o ciclo de vida do seu software.

---

## ✨ Funcionalidades Principais

-   📦 **`init`**: Configura automaticamente as branches `main` e `develop`.
-   📊 **`status`**: Visualização rápida do estado do seu fluxo atual.
-   🛡️ **Segurança**: Bloqueia finalizações se houver alterações não commitadas.
-   🌐 **Resiliência**: Funciona perfeitamente em repositórios locais (ignora erros de `origin` se não configurado).
-   🚀 **Automatização de Release**: Faz merges duplos (`main` e `develop`) e cria tags de versão com um único comando.

---

## 🛠️ Instalação (macOS)

Certifique-se de ter o [Go](https://golang.org/dl/) instalado.

1. Clone o repositório.
2. No diretório do projeto, execute:
   ```bash
   make install
   ```
   *Isso compilará o binário e o instalará em `/usr/local/bin/gitfw`.*

---

## 📖 Guia de Uso

### 1. Inicializando
```bash
gitfw init
```

### 2. Fluxo de Funcionalidade (Feature)
```bash
# Iniciar uma nova feature
gitfw feature start billing-system

# Finalizar e mesclar em develop
gitfw feature finish billing-system
```

### 3. Fluxo de Lançamento (Release)
```bash
# Iniciar uma nova versão
gitfw release start 1.3.0

# Finalizar (mescla em main e develop + cria tag v1.3.0)
gitfw release finish 1.3.0
```

### 4. Correções Críticas (Hotfix)
```bash
gitfw hotfix start security-patch
gitfw hotfix finish security-patch
```

---

## 🎨 Outros Comandos

-   `gitfw status`: Mostra branch atual, estado do repo e remoto.
-   `gitfw prune`: Limpa branches locais que já foram mescladas nas branches principais.
-   `gitfw feature list`: Lista todas as features ativas.
-   `gitfw version`: Exibe a versão atual do CLI.

---

## 🔧 Estrutura de Branches (Gitflow)
-   `main`: Onde o código de produção reside (tagueado).
-   `develop`: Onde o desenvolvimento ocorre e onde as features são integradas.
-   `feature/*`: Branches temporárias para novas funcionalidades.
-   `release/*`: Preparação da próxima versão.
-   `hotfix/*`: Correções urgentes para a `main`.

---

**Desenvolvido com foco em Clean Code e Robustez.** 🛠️
