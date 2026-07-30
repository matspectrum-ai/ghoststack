# GhostStack Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-07-30

### Added
- FASE R1 — WireGuard Real: TUN via /dev/net/tun+ioctl, firewall nft/iptables, DNS resolvectl/resolv.conf, wg-quick lifecycle
- FASE R2 — Provider Engine: 5 providers (WireGuard, Tor, sing-box, Unbound, SOCKS5 inline)
- FASE R3 — Plugin Loader + Sandbox: subprocess IPC via Unix socket, seccomp-bpf with syscall whitelist
- FASE R4 — WebSocket + Monitoramento: WSHub com broadcast push, MetricsCollector com /proc, SSE endpoint, dashboard com gráficos
- FASE R5 — Segurança Real: kill switch iptables/nft integrado ao WireGuard, SHA-256 binary verification, secret rotation com audit trail

### Changed
- CLI: `ghost init` com template YAML, `ghost emergency-stop`, `ghost diagnose --force`, `ghost start --provider`
- API: WebSocket em `/api/ws`, SSE em `/api/events`, POST `/api/config`
- Dashboard: ProviderStatus bar, TrafficChart SVG, PeerList, AlertPanel, ConfigScreen com WireGuard form

### Fixed
- seccomp `unsafePtr` retornava NULL → sandbox era no-op
- SOCKS5 address parsing (RSV+ATYP como port) + relay deadlock
- TUN netmask negative shift panic, address offset errado, interface leak
- Firewall command injection via `sh -c`, stdin não pipeado, erros ignorados
- HTTP server shutdown em 2s
- Plugin IPC: deadlock Unload→Disable, socket fechado antes do subprocesso
- Process state race (Stop sobrescrito por wait goroutine)
- Monitoramento coletava /proc real
- Audit flush duplicava entradas em erro
- Kill switch com trigger stub → agora chama firewall real

## [0.2.0] - 2026-07-30

### Added
- FASE 3 — Networking Engine: Gateway, TUN, routing, DNS, firewall com testes
- FASE 4 — Security Engine: Kill switch, integrity checker, audit logger
- FASE 5 — Providers: WireGuard, Tor, sing-box, Unbound, SOCKS5
- FASE 6 — CLI: ghost start/stop/status/diagnose/security/update com Cobra
- FASE 7 — Dashboard: React + Vite + TS com telas de status, monitoramento, config e logs
- FASE 8 — Distribution: goreleaser config para Linux, Windows, macOS
- FASE 9 — Hardening: fuzz tests e benchmarks

### Changed
- Daemon integrado com config, networking, security e providers
- API server com endpoints /api/status, /api/monitoring, /api/logs
- Dashboard serve arquivos estáticos via GHOSTSTACK_DASHBOARD_DIR

### Fixed
- CLI stop trata runtime não iniciado gracefulmente

## [0.1.0] - 2026-07-30

### Added
- Estrutura inicial do repositório
- FASE 0 — Foundation: Go module, Makefile, CI, cmd/ghost/main.go
- FASE 1 — Core Runtime: Ghost Daemon, lifecycle, eventos
- FASE 2 — Configuration Engine: YAML loader, validação, profiles, secrets
- Documentação base: docs/01-42, ADR-001, PLAN.md, AGENTS.md, CLAUDE.md
