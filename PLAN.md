# GhostStack v0.5.0 — Remote Coordination

**Branch:** `release/v0.5.0` (a criar)
**Base:** `master` (v0.4.0)
**Ship:** único (F1+F2+F3)

## Objetivo

Transformar o daemon local em um sistema multi-máquina com API autenticada. Um controller coordena N agents remotos via REST + API Key.

## Arquitetura

```
┌─────────────────────┐     HTTPS/API Key      ┌──────────────────────┐
│  Controller          │◄──────────────────────│  Agent (node-1)      │
│  (ghost start        │                       │  (ghost agent)       │
│   --public)          │── POST /api/v2/agents │                      │
│                      │── POST /heartbeat     │  providers locais    │
│  /api/v2/agents/*    │── GET  /command       │  SQLite local        │
│  SQLite (global)     │                       └──────────────────────┘
│  Dashboard           │     HTTPS/API Key      ┌──────────────────────┐
└─────────────────────┘◄──────────────────────│  Agent (node-N)      │
        ▲                                      └──────────────────────┘
        │ CLI (ghost remote)
   ┌────┴────┐
   │  dev    │
   └─────────┘
```

---

## Fase 1 — Auth + Public Bind

### Arquivos novos

| Arquivo | O que |
|---------|-------|
| `internal/api/auth.go` | APIKeyStore (SQLite), geração/validação de keys, middleware |
| `internal/cli/apikey.go` | `ghost apikey create\|list\|revoke` |

### Arquivos modificados

| Arquivo | O que muda |
|---------|------------|
| `internal/api/server.go` | Adiciona middleware de auth nas rotas existentes |
| `internal/api/tls.go` | (se necessário) suporte a bind `0.0.0.0` |
| `internal/cli/start.go` | Flag `--public` para bind externo |
| `internal/storage/provider.go` | `APIKey` struct + `SaveAPIKey`/`LoadAPIKey`/`DeleteAPIKey`/`ListAPIKeys` |
| `internal/storage/sqlite.go` | Migração `api_keys` + implementação dos métodos |
| `internal/cli/root.go` | `newAPIKeyCommand()` adicionado |

### Tasks

1. **Extender StorageProvider** — structs + interface para api_keys
2. **SQLite migrations** — tabela `api_keys` (id, key_hash, name, created_at, last_used_at, revoked)
3. **auth.go** — `APIKeyStore` struct, `NewAPIKeyStore(db)`, `Generate()`, `Validate()`, `Revoke()`, `List()`
4. **Middleware** — `AuthMiddleware(store)` que lê `Authorization: Bearer <key>` e valida
5. **server.go** — aplicar middleware nas rotas atuais
6. **start.go** — flag `--public` (bind `0.0.0.0:<port>`), `--api-key` opcional (gera se nao passar)
7. **CLI `ghost apikey`** — create (retorna key uma vez), list, revoke
8. **Testes** — auth middleware, key gen/val, CLI

### Critérios de aceitação

- `ghost start --public` bind em `0.0.0.0` com auth
- Request sem `Authorization: Bearer <key>` retorna 401
- `ghost apikey create` gera key e persiste no SQLite
- `ghost apikey revoke <name>` revoga key
- `make ci` passa

---

## Fase 2 — Agent Protocol

### Arquivos novos

| Arquivo | O que |
|---------|-------|
| `internal/agent/agent.go` | Agent runtime (registro, heartbeat, comando) |
| `internal/api/v2/server.go` | Router v2 para agent endpoints |
| `internal/api/v2/agents.go` | Handlers CRUD agents |
| `internal/api/v2/commands.go` | Handlers execução remota |
| `internal/cli/agent.go` | `ghost agent --controller <url> --api-key <key>` |

### Arquivos modificados

| Arquivo | O que muda |
|---------|------------|
| `internal/storage/provider.go` | `Agent` struct + `SaveAgent`/`LoadAgent`/`ListAgents`/`DeleteAgent` + `SaveCommand`/`PendingCommands` |
| `internal/storage/sqlite.go` | Migração `agents` + `commands` + implementação |
| `internal/api/server.go` | Monta router v2 em `/api/v2/` |

### Tasks

1. **Extender StorageProvider** — structs + interface p/ agents e commands
2. **SQLite migrations** — tabela `agents` (id, name, url, api_key, status, last_heartbeat, version, metadata), tabela `commands` (id, agent_id, action, payload, status, created_at, executed_at, result)
3. **Agent runtime** → loop: register → heartbeat (30s) → poll commands → execute → report
4. **CLI `ghost agent`** — `ghost agent start --controller <url> --api-key <key>`
5. **API v2 controller side** — register, heartbeat, commands/pending, commands/result, list agents
6. **Testes** — agent register/heartbeat, poll+exec, v2 handlers

### Critérios de aceitação

- `ghost agent start --controller http://c:8080 --api-key xyz` registra e inicia heartbeat
- Controller persiste agent no SQLite
- Heartbeat atualiza status visível via GET /api/v2/agents
- Comandos pendentes são entregues no poll
- `make ci` passa

---

## Fase 3 — Remote CLI + WebBridge

### Arquivos novos

| Arquivo | O que |
|---------|-------|
| `internal/cli/remote.go` | `ghost remote add|list|exec|logs|rm` |

### Arquivos modificados

| Arquivo | O que muda |
|---------|------------|
| `internal/cli/root.go` | `newRemoteCommand()` adicionado |
| `internal/api/server.go` | WebBridge: eventos de agent propagados via WS |
| `internal/api/ws.go` | `BroadcastAgentEvent()` |

### Tasks

1. **CLI `ghost remote`** — add/list/rm/exec/logs com persistência em `~/.ghoststack/remotes.json`
2. **WebBridge** — eventos de agent (register, heartbeat, command result) propagados para WS clients
3. **Testes** — remote CLI, WebBridge

### Critérios de aceitação

- `ghost remote exec node1 status` mostra estado remoto
- `ghost remote exec node1 stop` envia comando via controller
- Dashboard recebe eventos de agent via WebSocket
- `make ci` passa

---

## Dependências

```
F1 (auth+bind) ──→ F2 (agent protocol) ──→ F3 (remote CLI)
       │                    │
       ↓                    ↓
   Testes F1            Testes F2 + F3
```

F1 é standalone. F2 depende de F1 (auth + key store). F3 depende de F2.

## Estimativa

| Fase | Arquivos | Esforço |
|------|----------|---------|
| F1 | ~6 new, ~5 mod | 1-2h |
| F2 | ~5 new, ~5 mod | 2-3h |
| F3 | ~1 new, ~3 mod | 1h |
| Testes | ~vários | 30min |
| **Total** | ~12 new, ~13 mod | **4-6h** |

## Verificação

```bash
make ci
make test-race
go build -o /dev/null ./cmd/ghost
```

## GSTACK REVIEW REPORT

### Revisão CEO

- **Decisão:** Escopo aprovado — Remote Coordination desbloqueia multi-máquina com mínimo de complexidade
- **Observação:** API key + multi-máquina é o menor passo que entrega valor real para coordenação remota
- **Risco:** Agent protocol adiciona estado distribuído — heartbeat timeout e reconnect precisam ser robustos

### Revisão Engineering

- **Decisão:** Arquitetura aprovada — F1-F3 sequenciais sem refatoração necessária
- **Observação:** StorageProvider precisa extensão, não refatoração — baixo risco
- **Risco:** Agent HTTP polling pode escalar mal com 50+ agents; usar backoff exponencial no poll

### Revisão Design

- **Decisão:** API v2 endpoints seguem padrão REST — sem quebra de compatibilidade com v1
- **Observação:** WebBridge reusa WSHub existente — zero nova infra
- **Risco:** CLI `ghost remote exec` precisa tratar timeout de rede

### Revisão DX

- **Decisão:** CLI commands seguem padrão cobra existente — sem quebra de DX
- **Observação:** `ghost agent start` espelha `ghost start` — consistência mantida
- **Risco:** `--controller --api-key` na linha de comando pode expor key no history

### Decisões de Gosto

1. API key via `Authorization: Bearer` — não inventar protocolo custom
2. Agent polling com backoff exponencial (1s, 2s, 4s, max 60s) — não WebSocket bidirecional
3. `ghost remote` persiste em `~/.ghoststack/remotes.json` — não no SQLite (pode existir sem daemon)
4. Chaves do tipo `gs_` + 40 chars hex — prefixo identificável, comprimento seguro

### Aprovação

- **Status:** Aprovado para execução
- **Branch:** `release/v0.5.0`
- **Próximo:** Execução das 3 fases + `/ship`
