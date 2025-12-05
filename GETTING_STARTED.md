# 🚀 Getting Started with AlgoShield

## Início Rápido (5 minutos)

### Passo 1: Iniciar com Docker

```bash
# Iniciar todos os serviços
docker-compose up -d

# Aguardar inicialização (cerca de 30 segundos)
docker-compose logs -f
```

### Passo 2: Verificar se está funcionando

```bash
# Testar API
curl http://localhost:8080/health

# Deve retornar: {"status":"ok","postgres":"healthy","redis":"healthy"}
```

### Passo 3: Acessar a UI

Abra seu navegador em: **http://localhost:5173**

Você verá a interface de gerenciamento de regras com algumas regras de exemplo já configuradas.

### Passo 4: Enviar sua primeira transação

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "txn_teste_001",
    "amount": 1000.00,
    "currency": "BRL",
    "from_account": "CONTA001",
    "to_account": "CONTA002",
    "type": "transferencia",
    "metadata": {
      "ip": "192.168.1.1"
    },
    "timestamp": "2024-12-05T10:00:00Z"
  }'
```

### Passo 5: Consultar resultado

```bash
# Listar todas as transações
curl http://localhost:8080/api/v1/transactions

# Você verá sua transação processada com o score de risco!
```

## 🎓 Próximos Passos

### 1. Criar sua primeira regra customizada

**Via UI:**
1. Acesse http://localhost:5173
2. Clique em "Create Rule"
3. Preencha os campos:
   - Nome: "Transações Grandes"
   - Tipo: "amount"
   - Ação: "review"
   - Score: 50
   - Conditions: `{"amount_threshold": 5000}`
4. Clique em "Create"

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Transações Grandes",
    "description": "Revisar transações acima de R$ 5.000",
    "type": "amount",
    "action": "review",
    "priority": 10,
    "enabled": true,
    "conditions": {
      "amount_threshold": 5000
    },
    "score": 50
  }'
```

### 2. Testar sua regra

```bash
# Enviar transação que dispara a regra
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "txn_grande_001",
    "amount": 7500.00,
    "currency": "BRL",
    "from_account": "CONTA001",
    "to_account": "CONTA003",
    "type": "transferencia",
    "metadata": {},
    "timestamp": "2024-12-05T10:00:00Z"
  }'

# Aguardar 1 segundo e consultar
sleep 1
curl http://localhost:8080/api/v1/transactions | jq '.'
```

### 3. Explorar tipos de regras

#### Regra de Velocidade (Anti-fraude)
```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Muitas Transações Rápidas",
    "description": "Detectar mais de 5 transações em 10 minutos",
    "type": "velocity",
    "action": "block",
    "priority": 5,
    "enabled": true,
    "conditions": {
      "transaction_count": 5,
      "time_window_seconds": 600
    },
    "score": 80
  }'
```

#### Regra de Blacklist (PLD)
```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Contas Bloqueadas",
    "description": "Bloquear transações de contas suspeitas",
    "type": "blacklist",
    "action": "block",
    "priority": 1,
    "enabled": true,
    "conditions": {
      "blacklisted_accounts": ["CONTA_SUSPEITA_001", "CONTA_SUSPEITA_002"]
    },
    "score": 100
  }'
```

## 💻 Desenvolvimento Local

Se você quiser desenvolver localmente sem Docker:

### 1. Pré-requisitos
```bash
# Go 1.23+
go version

# Node.js 20+
node --version

# PostgreSQL & Redis via Docker
docker-compose up -d postgres redis
```

### 2. Configurar ambiente
```bash
# Copiar arquivo de configuração
cp env.example .env

# Editar .env e ajustar para localhost
# POSTGRES_HOST=localhost
# REDIS_HOST=localhost
```

### 3. Inicializar banco de dados
```bash
# Executar migrations
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/001_initial_schema.sql
```

### 4. Executar serviços

**Terminal 1 - API:**
```bash
make run-api
```

**Terminal 2 - Worker:**
```bash
make run-worker
```

**Terminal 3 - UI:**
```bash
cd ui
npm install
npm run dev
```

## 📊 Monitoramento

### Ver logs em tempo real
```bash
# Todos os serviços
docker-compose logs -f

# Apenas API
docker-compose logs -f api

# Apenas Worker
docker-compose logs -f worker
```

### Health Checks
```bash
# Status geral
curl http://localhost:8080/health

# Readiness
curl http://localhost:8080/ready
```

### Métricas de performance
Cada transação inclui tempo de processamento:
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{"external_id":"test","amount":100,"currency":"BRL","from_account":"A","to_account":"B","type":"transfer","timestamp":"2024-12-05T10:00:00Z"}' \
  | jq '.processing_time'
```

## 🧪 Teste de Carga

### Teste simples com curl
```bash
# Enviar 100 transações
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/transactions \
    -H "Content-Type: application/json" \
    -d "{
      \"external_id\": \"load_test_$i\",
      \"amount\": $((RANDOM % 10000 + 100)),
      \"currency\": \"BRL\",
      \"from_account\": \"CONTA001\",
      \"to_account\": \"CONTA002\",
      \"type\": \"transferencia\",
      \"timestamp\": \"2024-12-05T10:00:00Z\"
    }" &
done
wait

# Verificar quantas foram processadas
sleep 2
curl http://localhost:8080/api/v1/transactions | jq '.transactions | length'
```

## 🔍 Troubleshooting

### Problema: API não responde
```bash
# Verificar se está rodando
docker-compose ps api

# Ver logs
docker-compose logs api

# Reiniciar
docker-compose restart api
```

### Problema: Transações não são processadas
```bash
# Verificar workers
docker-compose ps worker

# Ver logs do worker
docker-compose logs worker

# Verificar fila do Redis
docker-compose exec redis redis-cli llen transaction:queue
```

### Problema: UI não carrega
```bash
# Verificar UI
docker-compose ps ui

# Recompilar UI
docker-compose up -d --build ui
```

### Problema: Erro de conexão com banco
```bash
# Verificar PostgreSQL
docker-compose ps postgres

# Testar conexão
docker-compose exec postgres psql -U algoshield -d algoshield -c "SELECT 1;"
```

## 📚 Documentação Adicional

- [README.md](./README.md) - Documentação principal
- [QUICKSTART.md](./docs/QUICKSTART.md) - Guia rápido
- [API_EXAMPLES.md](./docs/API_EXAMPLES.md) - Exemplos de API
- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Arquitetura
- [PROJECT_OVERVIEW.md](./PROJECT_OVERVIEW.md) - Visão geral técnica

## 🎯 Casos de Uso Comuns

### 1. E-commerce - Prevenção de Fraude
```bash
# Regra: Bloquear compras muito altas
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Compra Suspeita Alto Valor",
    "type": "amount",
    "action": "block",
    "priority": 5,
    "enabled": true,
    "conditions": {"amount_threshold": 15000},
    "score": 100
  }'

# Regra: Múltiplas compras em curto período
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Múltiplas Compras Rápidas",
    "type": "velocity",
    "action": "review",
    "priority": 10,
    "enabled": true,
    "conditions": {
      "transaction_count": 3,
      "time_window_seconds": 300
    },
    "score": 60
  }'
```

### 2. Banco - Anti-Lavagem de Dinheiro (PLD)
```bash
# Regra: Transferências internacionais altas
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Transferência Internacional Alta",
    "type": "pattern",
    "action": "review",
    "priority": 15,
    "enabled": true,
    "conditions": {"pattern": "internacional"},
    "score": 70
  }'
```

### 3. Fintech - Onboarding
```bash
# Regra: Primeira transação alta
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Primeira Transação Alta",
    "type": "amount",
    "action": "review",
    "priority": 20,
    "enabled": true,
    "conditions": {"amount_threshold": 3000},
    "score": 50
  }'
```

## ✅ Checklist de Produção

Antes de colocar em produção:

- [ ] Alterar senhas padrão no `.env`
- [ ] Configurar backup do PostgreSQL
- [ ] Configurar monitoramento (Prometheus/Grafana)
- [ ] Configurar logs centralizados (ELK/Loki)
- [ ] Implementar autenticação na API
- [ ] Configurar rate limiting
- [ ] Setup SSL/TLS
- [ ] Configurar alertas
- [ ] Documentar regras de negócio
- [ ] Treinar equipe operacional

## 🚀 Pronto para começar!

Agora você tem tudo para começar a usar o AlgoShield. Boa proteção contra fraudes! 🛡️

**Dúvidas?** Consulte a documentação ou abra uma issue no GitHub.

