# GhostStack — Project Context

Versão do milestone: 0.3.0
Status: Active
Branch: master
Última revisão: gsd-new-milestone (2026-07-30)

## 1. Missão

Evoluir o GhostStack de runtime base para plataforma extensível, observável e pronta para produção, mantendo privacidade, segurança e multi-plataforma como pilares centrais.

## 2. Contexto

O milestone v0.2.0 entregou a base funcional: Core Runtime, Networking Engine, Security Engine, Providers iniciais (WireGuard, Tor, sing-box, SOCKS5 Proxy), CLI básico e Dashboard inicial. A arquitetura estabilizou-se em torno de Go para o Core, TypeScript + React + Vite para o Dashboard, e YAML/JSON Schema para configuração.

Com v0.3.0, o foco desloca-se de "funciona" para "escala": extensibilidade via plugins, atualização segura, hardening de providers, observabilidade real e experiência de usuário madura.

## 3. Objetivos do Milestone

1. Tornar o GhostStack extensível sem modificar o Core — Plugin System completo.
2. Permitir atualização e migração seguras — Update System & Migration.
3. Aumentar a superfície de providers e endurecer o firewall — Provider Hardening.
4. Dar observabilidade de produção — Event Bus & Observability.
5. Ampliar o poder da CLI — CLI Expansion.
6. Transformar o Dashboard em produto — Dashboard Hardening.

## 4. Fora do Escopo (v0.3.0)

- Novos providers além dos listados em FASE 12.
- WASI / WebAssembly plugins.
- Container runtime nativo.
- Modo headless-only.
- Coleta de telemetria externa.

## 5. Stack Oficial (mantida)

- Core: Go
- Dashboard: TypeScript + React + Vite
- Configuração: YAML 1.2 + JSON Schema
- Plugins: Go SDK + TypeScript SDK
- Observabilidade: OpenTelemetry + Prometheus

## 6. Próximos Passos Imediatos

- Iniciar FASE 10 — Plugin System
- Revisar `docs/10-plugin-sdk-specification.md`
- Abrir ADR para Event Bus refatoração (pré-requisito FASE 13)
