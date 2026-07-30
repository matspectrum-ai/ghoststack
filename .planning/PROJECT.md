# GhostStack — Project Context

Versão do milestone: 0.4.0
Status: Active
Branch: master
Última revisão: gsd-new-milestone (2026-07-30)

## 1. Missão

Consolidar o GhostStack como plataforma cross-platform de privacidade production-ready, com desempenho, segurança e experiência de usuário de nível comercial.

## 2. Contexto

O milestone v0.3.0 entregou extensibilidade, observabilidade e CLI madura. A arquitetura está estável, com plugin system, event bus, providers hardened e dashboard responsivo. O próximo passo é levar o GhostStack de "plataforma funcional" para "produto production-ready" com foco em performance, segurança e experiência.

## 3. Objetivos do Milestone

1. Performance e otimização — reduzir latência e footprint.
2. Segurança avançada — audit, sandboxing e hardening.
3. Cross-platform real — suporte completo a Windows, macOS e Linux.
4. Experiência de usuário — onboarding, docs e DX.
5. Integração e ecossistema — SDKs, APIs e extensibilidade.

## 4. Fora do Escopo (v0.4.0)

- Novos providers VPN além dos já suportados.
- Container runtime nativo.
- Modo headless-only.
- Telemetria externa.

## 5. Stack Oficial (mantida)

- Core: Go
- Dashboard: TypeScript + React + Vite
- Configuração: YAML 1.2 + JSON Schema
- Plugins: Go SDK + TypeScript SDK
- Observabilidade: OpenTelemetry + Prometheus

## 6. Próximos Passos Imediatos

- Iniciar FASE 16 — Performance & Optimization
- Revisar docs/16-performance-optimization.md
- Abrir ADR para sandboxing de plugins (FASE 17)
