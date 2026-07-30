# GhostStack v0.3.0 — Requirements

Versão: 0.3.0
Status: Em andamento
Base: PLAN.md FASEs 10-15

## REQ-01 — Plugin System (FASE 10)

- Plugin Manager com descoberta em diretórios oficiais.
- Manifest parser com validação de campos obrigatórios.
- Permissions model: network, filesystem, secrets, system, events.
- Plugin SDK base em Go com ciclo de vida completo.
- Sandbox básico por processo separado.
- Assinatura digital opcional de plugins.

## REQ-02 — Update System & Migration (FASE 11)

- Update Manager com canais stable, beta, nightly.
- Version Resolver compatível com SemVer.
- Migration Engine para configuração e dados.
- Rollback automático em falha de migração.
- Update Manifest com checksum e assinatura.

## REQ-03 — Provider Hardening (FASE 12)

- OpenVPN provider compatível com arquitetura cross-platform.
- DoH/DoT DNS providers com fallback automático.
- Firewall real: nftables, iptables, pf por plataforma.
- Health checks em todos os providers.
- Fallback automático entre providers.

## REQ-04 — Event Bus & Observability (FASE 13)

- Event Bus real pub/sub com correlation IDs.
- OpenTelemetry básico integrado.
- Métricas Prometheus padronizadas por módulo.
- Tracing distribuído opcional.
- Health endpoints: /health, /ready, /alive.

## REQ-05 — CLI Expansion (FASE 14)

- ghost plugin install/remove/update.
- ghost config validate/reload.
- ghost provider list/select.
- ghost update --dry-run/rollback.
- ghost doctor e ghost backup/restore.

## REQ-06 — Dashboard Hardening (FASE 15)

- WebSocket para eventos real-time.
- Telas funcionais de configuração e logs.
- Mobile responsive com redesign das telas.
- PWA support offline.
- Autenticação local obrigatória.

## Requisitos Transversais

- Cobertura mínima: Core 95%, Plugins 90%, CLI 85%.
- Todo módulo novo deve expor health checks.
- Sanitização de logs obrigatória.
- CI multi-plataforma: Linux amd64, Windows amd64, macOS amd64/arm64.
- Documentação atualizada e exemplos de plugin.
