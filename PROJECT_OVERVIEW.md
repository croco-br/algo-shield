# AlgoShield - Visão Geral do Projeto

## 🎯 Objetivo

AlgoShield é uma ferramenta open-source de análise de transações para prevenção de fraude e lavagem de dinheiro (AML), projetada para processar cada transação com latência inferior a 50ms.

## 📁 Estrutura do Projeto

```
algo-shield/
├── api/                          # Serviço de API (Fiber)
│   ├── cmd/
│   │   └── main.go              # Entry point da API
│   └── internal/
│       ├── handlers/            # Handlers HTTP
│       │   ├── health.go        # Health checks
│       │   ├── rule.go          # Gerenciamento de regras
│       │   └── transaction.go   # Processamento de transações
│       ├── middleware/          # Middleware HTTP
│       │   ├── cors.go          # CORS configuration
│       │   └── logger.go        # Request logging
│       └── routes/
│           └── routes.go        # Definição de rotas
│
├── workers/                      # Serviço de Workers
│   ├── cmd/
│   │   └── main.go              # Entry point do Worker
│   └── internal/
│       ├── processor/
│       │   └── processor.go     # Processador de transações
│       └── rules/
│           └── engine.go        # Motor de regras
│
├── pkg/                          # Pacotes compartilhados
│   ├── config/
│   │   ├── config.go            # Configuração da aplicação
│   │   └── config_test.go       # Testes de configuração
│   ├── database/
│   │   ├── postgres.go          # Cliente PostgreSQL
│   │   └── redis.go             # Cliente Redis
│   ├── models/
│   │   ├── rule.go              # Modelo de regras
│   │   ├── transaction.go       # Modelo de transações
│   │   └── transaction_test.go  # Testes de modelos
│   └── utils/
│       └── logger.go            # Utilitários de logging
│
├── ui/                           # Interface Web (SvelteKit)
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +layout.svelte   # Layout principal
│   │   │   └── +page.svelte     # Página de gerenciamento de regras
│   │   ├── app.css              # Estilos globais
│   │   └── app.html             # Template HTML
│   ├── package.json
│   ├── svelte.config.js
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── scripts/                      # Scripts de inicialização
│   ├── migrations/
│   │   └── 001_initial_schema.sql  # Schema inicial do banco
│   └── init-db.sh               # Script de inicialização do DB
│
├── docs/                         # Documentação
│   ├── API_EXAMPLES.md          # Exemplos de uso da API
│   ├── ARCHITECTURE.md          # Documentação da arquitetura
│   └── QUICKSTART.md            # Guia de início rápido
│
├── .github/
│   └── workflows/
│       └── ci.yml               # CI/CD pipeline
│
├── Dockerfile.api               # Dockerfile para API
├── Dockerfile.worker            # Dockerfile para Worker
├── Dockerfile.ui                # Dockerfile para UI
├── docker-compose.yml           # Orquestração de containers
├── Makefile                     # Comandos de build e desenvolvimento
├── go.mod                       # Dependências Go
├── go.sum                       # Checksums das dependências
├── env.example                  # Exemplo de configuração
├── CONTRIBUTING.md              # Guia de contribuição
└── README.md                    # Documentação principal
```

## 🔧 Tecnologias Utilizadas

### Backend (Go 1.23)
- **Fiber v2**: Framework web de alta performance
- **pgx v5**: Driver PostgreSQL otimizado
- **go-redis v9**: Cliente Redis
- **uuid**: Geração de UUIDs
- **godotenv**: Gerenciamento de variáveis de ambiente

### Frontend (SvelteKit)
- **Svelte 5**: Framework reativo moderno
- **SvelteKit 2**: Framework full-stack
- **TypeScript**: Type safety
- **Vite**: Build tool rápido

### Infraestrutura
- **PostgreSQL 16**: Banco de dados principal
- **Redis 7**: Cache e fila de mensagens
- **Docker**: Containerização
- **Docker Compose**: Orquestração local

## 🚀 Componentes Principais

### 1. API Service (api/)
- **Responsabilidade**: Receber transações, gerenciar regras, fornecer interface de consulta
- **Porta**: 8080
- **Endpoints principais**:
  - `POST /api/v1/transactions` - Submeter transação
  - `GET /api/v1/transactions` - Listar transações
  - `POST /api/v1/rules` - Criar regra
  - `GET /api/v1/rules` - Listar regras
  - `GET /health` - Health check

### 2. Worker Service (workers/)
- **Responsabilidade**: Processar transações através do motor de regras
- **Concorrência**: Configurável via `WORKER_CONCURRENCY` (padrão: 10)
- **Hot-reload**: Recarrega regras a cada 10 segundos automaticamente
- **Processamento**: Assíncrono via fila Redis

### 3. Rules Engine (workers/internal/rules/)
- **Tipos de regras suportadas**:
  - **Amount**: Threshold de valor
  - **Velocity**: Frequência de transações
  - **Blacklist**: Lista de contas bloqueadas
  - **Pattern**: Matching de padrões
  - **Custom**: Extensível para novas regras

### 4. UI (ui/)
- **Responsabilidade**: Interface para gerenciamento de regras
- **Porta**: 5173
- **Features**:
  - Criar/editar/deletar regras
  - Ativar/desativar regras em tempo real
  - Visualizar configurações de regras
  - UI moderna e responsiva

## 📊 Fluxo de Dados

```
1. Cliente envia transação → API
2. API valida e coloca na fila Redis
3. Worker pega transação da fila
4. Worker carrega regras (cache Redis)
5. Worker avalia transação contra regras
6. Worker calcula score e status de risco
7. Worker salva resultado no PostgreSQL
8. Cliente consulta resultado via API
```

## 🎯 Características de Performance

### Otimizações Implementadas
- ✅ Compilação com Go 1.23 + flag `GOEXPERIMENT=rangefunc`
- ✅ Connection pooling PostgreSQL (50 max connections)
- ✅ Connection pooling Redis (50 connections)
- ✅ Cache de regras em Redis (TTL: 5 minutos)
- ✅ Processamento assíncrono via fila
- ✅ Workers concorrentes configuráveis
- ✅ Índices otimizados no banco de dados
- ✅ Zero-allocation routing (Fiber)

### Metas de Performance
- **Latência da API**: <5ms (operação de fila)
- **Processamento de transação**: <50ms (fim-a-fim)
- **Avaliação de regras**: <30ms (média)
- **Throughput**: 10,000+ req/s por instância de API

## 🔐 Segurança

### Implementado
- Validação de entrada em todos os endpoints
- Queries parametrizadas (proteção contra SQL injection)
- CORS configurável
- Health checks para monitoramento

### Roadmap
- [ ] Autenticação JWT/OAuth2
- [ ] Rate limiting por cliente
- [ ] Criptografia em repouso
- [ ] Audit logging
- [ ] RBAC (Role-Based Access Control)

## 🧪 Testing

### Testes Unitários
```bash
make test
```

### Testes incluídos:
- ✅ Models (Transaction, Rule)
- ✅ Configuration loading
- ✅ Database DSN generation

### Roadmap de Testes
- [ ] Integration tests
- [ ] API endpoint tests
- [ ] Rules engine tests
- [ ] Performance benchmarks

## 📈 Escalabilidade

### Horizontal
- **API**: Múltiplas instâncias atrás de load balancer
- **Workers**: Escalar replicas baseado em profundidade da fila
- **Database**: Read replicas para consultas
- **Redis**: Redis Cluster para HA

### Vertical
- Aumentar `WORKER_CONCURRENCY`
- Aumentar pool de conexões do banco
- Otimizar queries com índices adicionais

## 🛠️ Comandos Úteis

### Docker
```bash
make docker-build    # Build de todas as imagens
make docker-up       # Iniciar todos os serviços
make docker-down     # Parar todos os serviços
make docker-logs     # Ver logs
```

### Desenvolvimento Local
```bash
make deps           # Instalar dependências
make build          # Build dos binários
make run-api        # Executar API
make run-worker     # Executar Worker
make dev-ui         # Executar UI em modo dev
make test           # Executar testes
```

## 📚 Documentação Adicional

- [README.md](./README.md) - Documentação principal
- [QUICKSTART.md](./docs/QUICKSTART.md) - Guia de início rápido
- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Arquitetura detalhada
- [API_EXAMPLES.md](./docs/API_EXAMPLES.md) - Exemplos de uso da API
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Como contribuir

## 🎯 Próximos Passos (Roadmap)

### Fase 1 - Core (✅ Completo)
- [x] API REST com Fiber
- [x] Worker com processamento assíncrono
- [x] Motor de regras básico
- [x] UI para gerenciamento de regras
- [x] Docker & docker-compose
- [x] Documentação completa

### Fase 2 - Produção
- [ ] Autenticação e autorização
- [ ] Rate limiting
- [ ] Métricas e monitoramento (Prometheus/Grafana)
- [ ] Logs estruturados (ELK/Loki)
- [ ] Distributed tracing (Jaeger/OpenTelemetry)
- [ ] Helm charts para Kubernetes

### Fase 3 - Features Avançadas
- [ ] Geração de dados sintéticos
- [ ] Machine Learning integration
- [ ] Dashboard de analytics
- [ ] Sistema de notificações
- [ ] API de webhooks
- [ ] Multi-tenancy

### Fase 4 - Otimizações
- [ ] GraphQL API
- [ ] gRPC para comunicação interna
- [ ] Event sourcing
- [ ] CQRS pattern
- [ ] Read-through cache
- [ ] Sharding do banco de dados

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, leia [CONTRIBUTING.md](./CONTRIBUTING.md) para detalhes sobre o processo de contribuição.

## 📝 Licença

Este projeto é licenciado sob a licença MIT - veja o arquivo [LICENSE](./LICENSE) para detalhes.

## 👥 Equipe

- Desenvolvido com ❤️ para a comunidade de prevenção a fraudes

## 📞 Suporte

- 🐛 [Reportar Bug](https://github.com/yourusername/algo-shield/issues)
- 💡 [Sugerir Feature](https://github.com/yourusername/algo-shield/issues)
- 💬 [Discussões](https://github.com/yourusername/algo-shield/discussions)

---

**AlgoShield** - Proteção inteligente contra fraude e lavagem de dinheiro 🛡️

