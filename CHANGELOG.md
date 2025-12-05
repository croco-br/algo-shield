# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Added
- Sistema completo de Git hooks para testes automatizados
  - Pre-commit: Testes unitários e verificação de formatação
  - Commit-msg: Validação de Conventional Commits
  - Pre-push: Suite completa de testes com coverage
- Script de instalação de hooks (`scripts/install-hooks.sh`)
- Documentação completa de Git hooks (`docs/GIT_HOOKS.md`)
- Arquivo de versões das tecnologias (`VERSIONS.md`)
- Comando `make install-hooks` para instalação fácil

### Changed
- **BREAKING**: Atualizado Go de `1.23` para `1.23.4`
- **BREAKING**: Atualizado PostgreSQL de `16-alpine` para `17-alpine`
- Atualizado Redis de `7-alpine` para `7.4-alpine`
- Atualizado Node.js de `20-alpine` para `22-alpine`
- Melhorado Makefile com instalação automática de dependências
- Atualizado CI/CD para usar versões mais recentes

### Fixed
- Correção na verificação de dependências do Makefile
- Permissões executáveis nos scripts de hooks

## [0.1.0] - 2024-12-05

### Added
- Estrutura inicial do projeto
- API Service com Fiber (Go)
  - Endpoints para transações
  - Endpoints para regras
  - Health checks
- Worker Service para processamento assíncrono
  - Motor de regras customizável
  - Suporte a 4 tipos de regras (amount, velocity, blacklist, pattern)
  - Hot-reload de regras
- UI com SvelteKit
  - Interface de gerenciamento de regras
  - CRUD completo de regras
- Infraestrutura Docker
  - docker-compose.yml completo
  - Dockerfiles otimizados (multi-stage build)
- Database
  - Schema PostgreSQL com índices otimizados
  - Migrations iniciais
  - Regras de exemplo
- Documentação completa
  - README.md detalhado
  - Guia de início rápido (QUICKSTART.md)
  - Exemplos de API (API_EXAMPLES.md)
  - Documentação de arquitetura (ARCHITECTURE.md)
  - Guia de contribuição (CONTRIBUTING.md)
- CI/CD
  - GitHub Actions workflow
  - Testes automatizados
  - Build verification
- Makefile com comandos úteis
- Testes unitários básicos

### Performance
- Compilação com `GOEXPERIMENT=rangefunc`
- Connection pooling otimizado (PostgreSQL e Redis)
- Cache de regras em Redis
- Processamento assíncrono via fila
- Índices de banco de dados para queries <10ms

### Security
- Validação de entrada em todos os endpoints
- Queries parametrizadas (proteção SQL injection)
- CORS configurável
- Health checks para monitoramento

---

## Tipos de Mudanças

- `Added` - Novas funcionalidades
- `Changed` - Mudanças em funcionalidades existentes
- `Deprecated` - Funcionalidades que serão removidas
- `Removed` - Funcionalidades removidas
- `Fixed` - Correções de bugs
- `Security` - Correções de vulnerabilidades
- `Performance` - Melhorias de performance

## Versionamento

Este projeto segue [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Mudanças incompatíveis na API
- **MINOR** (0.X.0): Novas funcionalidades compatíveis
- **PATCH** (0.0.X): Correções de bugs compatíveis

## Links

- [Repositório](https://github.com/yourusername/algo-shield)
- [Issues](https://github.com/yourusername/algo-shield/issues)
- [Pull Requests](https://github.com/yourusername/algo-shield/pulls)

