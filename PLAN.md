# GhostStack — Implementation Plan

Versão: 0.4.0
Status: Draft
Branch: master
Última revisão: autoplan (2026-07-30)

## 1. Contexto

Plano de implementação do GhostStack, derivado das especificações normativas `docs/01..42` e da ADR-001.
Revisão automática executada via `/autoplan` após conclusão das FASEs 0-9 (v0.2.0).

## 2. Escopo

- FASEs 0-9: ✅ Concluídas (v0.2.0)
- FASEs 10-15: ✅ Concluídas (v0.3.0)
- FASEs 16-20: ✅ Concluídas (v0.4.0)

## 3. Fases

### FASE 16 — Performance & Optimization
- Perfil de memória com pprof
- Startup time optimization
- Benchmark suite contínua
- Hot path profiling
- Memory pool reutilização

### FASE 17 — Security Hardening
- Audit logging estruturado (JSON)
- Plugin sandboxing
- Integrity checking com hash
- Secure boot do daemon
- Secret rotation

### FASE 18 — Cross-Platform Real
- Windows WFP firewall
- macOS pf + NetworkExtension
- Linux nftables + tun
- Build matrix GitHub Actions
- Platform abstraction layer

### FASE 19 — UX & DX Improvement
- Onboarding interativo
- Documentação interativa
- DX scorecard
- Error messages melhorados
- CLI help system

### FASE 20 — Integration & Ecosystem
- REST API pública
- TypeScript SDK
- Webhook system
- CLI scripting
- Plugin marketplace preparação


## 4. Critérios de Aceitação

- make build PASS
- Testes passando
- Documentação atualizada
- Segurança validada

## 5. Stack

- Core: Go
- Dashboard: TypeScript + React + Vite
- Configuração: YAML 1.2 + JSON Schema
- Plugins: Go SDK + TypeScript SDK
- Observability: OpenTelemetry + Prometheus

## 6. Próximos Passos

- Iniciar FASE 16 — Performance & Optimization
- Continuar execução do milestone v0.4.0

## GSTACK REVIEW REPORT

### Revisão CEO
- **Decisão:** Manter escopo atual, expandir gradualmente
- **Observação:** Plugin system é diferencial competitivo, priorizar FASE 10
- **Risco:** Complexidade de sandbox pode atrasar FASE 10

### Revisão Engineering
- **Decisão:** FASEs 10-15 são viáveis com arquitetura atual
- **Observação:** Event bus precisa refatoração para suportar plugins
- **Risco:** Firewall real requer trabalho específico por plataforma

### Revisão Design
- **Decisão:** Dashboard atual é base suficiente para FASE 15
- **Observação:** Mobile responsive requer redesign das telas existentes
- **Risco:** WebSocket integration pode quebrar telas atuais

### Revisão DX
- **Decisão:** CLI expansion melhora significativamente DX
- **Observação:** Plugin SDK precisa documentação extensa
- **Risco:** Curva de aprendizado para desenvolvedores de plugins

### Decisões de Gosto
1. Plugin System como próxima fase (não pulá-la)
2. Event bus refatoração antes de FASE 13
3. Firewall real como stub até FASE 12
4. Mobile responsive como última prioridade

### Aprovação
- **Status:** Aprovado para execução
- **Próximo milestone:** v0.4.0 (FASEs 16+)
- **Comando:** `/gsd-new-milestone v0.4.0`
