# GhostStack Status

## v0.4.0 — CLI Production-Ready

Status: **In progress** — Fases P1-P8 entregues, aguardando build/test local e ship.

### Entregue

| Fase | Status | O quê |
|------|--------|-------|
| P1 — Secrets Manager | ✅ | AES-256-GCM + argon2id, `ghost secrets`, integrado com init |
| P2 — TLS Nativo | ✅ | Auto-cert ECDSA P256, `--tls` flag, HTTPS redirect :8443, WSS |
| P3 — SQLite Storage | ✅ | CGO-free (modernc), migrations, persistência runtime/provider/audit |
| P4 — Systemd Service | ✅ | `ghost service install\|uninstall\|enable\|disable\|status` |
| P5 — Kill Switch v2 | ✅ | IPv6/ICMP/DNS leak protection, DNS force, leak test, audit log |
| P6 — Testes + CI | ✅ | Makefile, GitHub Actions (lint/test/build), golangci-lint |
| P7 — Docker | ✅ | Multi-stage Dockerfile, docker-compose, .dockerignore |
| P8 — Audit + Privacy | ✅ | `ghost audit privacy\|security\|all`, 15 checks, score |

### Pendente (local)

- `go mod tidy` para resolver dependências (modernc.org/sqlite + indirects)
- `make ci` para verificar lint + testes + build
- `./ghost version` para smoke test
- Ship via `/ship` após aprovação

### Comandos novos

```
ghost secrets init|set|get|list|delete
ghost start --tls
ghost db status|audit
ghost service install|uninstall|enable|disable|status
ghost killswitch enable|disable|status|leak-test
ghost audit privacy|security|all
ghost version
```
