# GhostStack v0.4.0 — Roadmap

Versão: 0.4.0
Status: Em andamento
Árvore: FASEs 16-20

## FASE 16 — Performance & Optimization

Entregas:
- Perfil de memória com pprof.
- Startup time optimization.
- Benchmark suite contínua.
- Hot path profiling.
- Memory pool reutilização.

Critérios de aceite:
- `ghost start` < 500ms.
- Memória idle < 50MB.
- Benchmarks passando.

## FASE 17 — Security Hardening

Entregas:
- Audit logging estruturado (JSON).
- Plugin sandboxing.
- Integrity checking com hash.
- Secure boot do daemon.
- Secret rotation.

Critérios de aceite:
- Todos os eventos auditados.
- Plugins isolados por processo.
- Integrity check em boot.

## FASE 18 — Cross-Platform Real

Entregas:
- Windows WFP firewall.
- macOS pf + NetworkExtension.
- Linux nftables + tun.
- Build matrix GitHub Actions.
- Platform abstraction layer.

Critérios de aceite:
- Builds nativos por OS.
- Firewall funcional em cada plataforma.
- TUN device abstraction.

## FASE 19 — UX & DX Improvement

Entregas:
- Onboarding interativo.
- Documentação interativa.
- DX scorecard.
- Error messages melhorados.
- CLI help system.

Critérios de aceite:
- Novo usuário operacional em < 5min.
- Docs cobrem 100% das APIs públicas.
- DX score > 8/10.

## FASE 20 — Integration & Ecosystem

Entregas:
- REST API pública.
- TypeScript SDK.
- Webhook system.
- CLI scripting.
- Plugin marketplace preparação.

Critérios de aceite:
- API documentada e testada.
- SDK TypeScript funcional.
- Webhooks configuráveis.
