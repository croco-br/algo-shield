# Git Hooks - AlgoShield

Este documento descreve os Git hooks configurados para o projeto AlgoShield.

## 📋 Índice

- [Instalação](#instalação)
- [Hooks Disponíveis](#hooks-disponíveis)
- [Como Funciona](#como-funciona)
- [Bypass de Hooks](#bypass-de-hooks)

## Instalação

Para instalar os Git hooks no seu ambiente local:

```bash
./scripts/install-hooks.sh
```

Isso irá:
1. Configurar o Git para usar o diretório `.githooks`
2. Tornar todos os hooks executáveis
3. Ativar os hooks automaticamente

## Hooks Disponíveis

### 1. Pre-Commit Hook

**Arquivo**: `.githooks/pre-commit`

**Executado**: Antes de cada commit

**Verificações**:
- ✅ Executa testes unitários Go (`go test -short ./...`)
- ✅ Verifica formatação do código Go (`gofmt`)
- ✅ Verifica consistência de `go.mod` e `go.sum`
- ✅ Adiciona automaticamente mudanças em go.mod/go.sum se necessário

**Exemplo de saída**:
```
🔍 Running pre-commit checks...
→ Running Go unit tests...
✓ Go tests passed
→ Checking Go formatting...
✓ Go formatting is correct
→ Checking Go modules...
✓ All pre-commit checks passed!
```

**Quando falha**:
```
✗ Go tests failed
Commit aborted. Fix the tests and try again.
```

### 2. Commit-Msg Hook

**Arquivo**: `.githooks/commit-msg`

**Executado**: Após escrever a mensagem de commit, antes de confirmar

**Verificações**:
- ✅ Valida formato de Conventional Commits
- ✅ Verifica comprimento da mensagem (aviso se >100 caracteres)

**Formato esperado**:
```
<type>(<scope>): <subject>
```

**Tipos válidos**:
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Mudanças em documentação
- `style`: Mudanças de formatação (não afetam o código)
- `refactor`: Refatoração de código
- `perf`: Melhorias de performance
- `test`: Adição ou atualização de testes
- `chore`: Mudanças em ferramentas, configurações
- `build`: Mudanças no sistema de build
- `ci`: Mudanças em CI/CD

**Exemplos válidos**:
```bash
git commit -m "feat: add transaction velocity check"
git commit -m "fix(api): correct rule evaluation logic"
git commit -m "docs: update API documentation"
git commit -m "test(rules): add unit tests for blacklist rules"
git commit -m "refactor(worker): improve processor performance"
```

**Exemplos inválidos**:
```bash
git commit -m "Added new feature"  # ✗ Não segue o formato
git commit -m "Fix bug"            # ✗ Falta o escopo e os dois pontos
git commit -m "update"             # ✗ Tipo inválido
```

### 3. Pre-Push Hook

**Arquivo**: `.githooks/pre-push`

**Executado**: Antes de fazer push para o repositório remoto

**Verificações**:
- ✅ Executa **todos** os testes (incluindo race detector)
- ✅ Gera relatório de coverage
- ✅ Verifica novos TODOs/FIXMEs (apenas aviso)

**Exemplo de saída**:
```
🚀 Running pre-push checks...
→ Running full test suite...
✓ All tests passed
✓ Test coverage: 67.8%
→ Checking for TODOs and FIXMEs...
⚠ Found 2 new TODO/FIXME comments
  Consider resolving them before pushing
✓ All pre-push checks passed!
```

## Como Funciona

### Fluxo de Trabalho

```
1. Você faz mudanças no código
2. git add <arquivos>
3. git commit -m "feat: nova funcionalidade"
   ↓
   [pre-commit hook]
   - Roda testes unitários
   - Verifica formatação
   - Verifica go.mod/go.sum
   ↓
   [commit-msg hook]
   - Valida formato da mensagem
   ↓
   Commit criado ✓

4. git push
   ↓
   [pre-push hook]
   - Roda suite completa de testes
   - Gera coverage report
   - Verifica TODOs
   ↓
   Push realizado ✓
```

### Configuração no Git

Os hooks são instalados configurando:
```bash
git config core.hooksPath .githooks
```

Isso faz o Git usar `.githooks/` ao invés de `.git/hooks/`

## Bypass de Hooks

### ⚠️ Importante
Bypassing hooks deve ser feito apenas em casos excepcionais!

### Para commits
```bash
git commit --no-verify -m "mensagem"
# ou
git commit -n -m "mensagem"
```

### Para push
```bash
git push --no-verify
# ou
git push -n
```

### Quando é aceitável fazer bypass?

✅ **Casos válidos**:
- Commits de documentação urgentes
- Hotfixes críticos em produção (mas rode os testes depois!)
- Quando os testes estão falhando por problemas de infraestrutura

❌ **Evite**:
- Para "economizar tempo"
- Quando os testes estão falhando por bugs no seu código
- Para pushar código quebrado

## Desinstalação

Para desabilitar os hooks:

```bash
git config --unset core.hooksPath
```

Ou apague o diretório `.githooks/`

## Troubleshooting

### Hook não está executando

Verifique se os hooks têm permissão de execução:
```bash
chmod +x .githooks/*
```

Verifique a configuração do Git:
```bash
git config --get core.hooksPath
# Deve retornar: .githooks
```

### Testes falhando no hook mas passando localmente

Certifique-se de que está na mesma versão do Go:
```bash
go version  # Deve ser go1.23.4 ou superior
```

Limpe o cache e rode novamente:
```bash
go clean -testcache
go test ./...
```

### Erro de permissão no Windows

No Windows, pode ser necessário usar Git Bash ou WSL para executar os hooks.

Alternativamente, você pode criar versões `.bat` dos hooks:
```bash
# Exemplo: .githooks/pre-commit.bat
@echo off
bash .githooks/pre-commit
```

## Customização

Para modificar os hooks, edite os arquivos em `.githooks/`:

- `.githooks/pre-commit` - Adicione mais verificações antes do commit
- `.githooks/commit-msg` - Mude o formato de mensagem requerido
- `.githooks/pre-push` - Adicione verificações antes do push

Após modificar, não esqueça de:
```bash
chmod +x .githooks/<hook-modificado>
```

## Boas Práticas

1. **Rode os testes localmente** antes de commitar:
   ```bash
   make test
   ```

2. **Use mensagens de commit descritivas** seguindo Conventional Commits

3. **Mantenha os commits pequenos** e focados em uma única mudança

4. **Não commite código comentado** ou debug logs

5. **Revise suas mudanças** antes de commitar:
   ```bash
   git diff --staged
   ```

## CI/CD Integration

Os mesmos testes rodados nos hooks locais são executados no CI (GitHub Actions).

Veja `.github/workflows/ci.yml` para mais detalhes.

## Referências

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Hooks Documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)
- [Go Testing](https://golang.org/pkg/testing/)

---

**AlgoShield** - Commits limpos, código testado! 🛡️

