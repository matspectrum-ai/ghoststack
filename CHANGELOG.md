# GhostStack Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
