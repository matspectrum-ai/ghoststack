# GhostStack v0.3.0 — Roadmap

Versão: 0.3.0
Status: Em andamento
Árvore: FASEs 10-15

## FASE 10 — Plugin System

Entregas:
- Plugin Manager com discovery em `/usr/lib/ghoststack/plugins`, `~/.local/share/ghoststack/plugins`, `./plugins`.
- Manifest parser e validação de compatibilidade SDK/ABI.
- Permissions model declarativa: network, filesystem, secrets, system, runtime, events.
- Plugin SDK base em Go.
- Sandbox básico por processo separado.
- Hot reload e rollback de plugins.

Critérios de aceite:
- Plugin externo instalável sem modificar o Core.
- Manifest inválido impede carregamento.
- Permissões não declaradas não são concedidas.

## FASE 11 — Update System & Migration

Entregas:
- Update Manager com canais stable, beta, nightly.
- Version Resolver SemVer.
- Migration Engine para configuração e dados.
- Rollback automático em falha.
- Update Manifest com checksum e assinatura.

Critérios de aceite:
- `ghost update` executa migração segura.
- Falha em migração reverte para versão anterior.
- Rollback manual via `ghost rollback` funciona.

## FASE 12 — Provider Hardening

Entregas:
- OpenVPN provider cross-platform.
- DoH/DoT DNS providers.
- Firewall real por SO: Linux nftables, Windows WFP, macOS pf.
- Health checks por provider.
- Fallback automático entre providers.

Critérios de aceite:
- Providers possuem Health() e Status().
- Falha de VPN ativa fallback automático.
- Firewall bloqueia tráfego fora da política.

## FASE 13 — Event Bus & Observability

Entregas:
- Event Bus real pub/sub.
- OpenTelemetry básico.
- Métricas Prometheus por módulo.
- Tracing distribuído.
- Correlação por correlationId.

Critérios de aceite:
- Plugins publicam/consomem eventos via contratos.
- Dashboard consome métricas via API oficial.
- Logs são sanitizados e padronizados em JSON.

## FASE 14 — CLI Expansion

Entregas:
- ghost plugin install/remove/update.
- ghost config validate/reload.
- ghost provider list/select.
- ghost update --dry-run/rollback.
- ghost doctor.

Critérios de aceite:
- Todo fluxo de gerenciamento via CLI.
- `ghost doctor` valida dependências e permissões.

## FASE 15 — Dashboard Hardening

Entregas:
- WebSocket para eventos real-time.
- Telas funcionais de configuração e logs.
- Mobile responsive.
- PWA support.
- Autenticação local obrigatória.

Critérios de aceite:
- Dashboard disponível em localhost:9090.
- Eventos refletem em tempo real sem reload.
- Interface utilizável em telas pequenas.
