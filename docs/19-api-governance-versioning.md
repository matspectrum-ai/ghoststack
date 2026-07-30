# GhostStack — API Governance & Versioning Specification

Versão: 1.0.0
Status: Normativo
Tipo: API Architecture Specification
Dependências: PRD, System Architecture, Domain Model, Configuration Specification, Contracts Specification, OpenAPI Specification, Event Specification, Plugin SDK Specification, Observability & Monitoring Specification

## 1. Objetivo

Definir as regras oficiais para criação, evolução e manutenção das APIs do GhostStack.

## 2. Princípios

As APIs do GhostStack devem ser: estáveis, previsíveis, documentadas, versionadas, retrocompatíveis quando possível.

## 3. Tipos de API

Quatro categorias:
- Internal API
- Core API
- Public API
- Plugin API

## 4. Internal API

Uso: comunicação interna entre módulos. Características: não possui garantia pública, pode mudar entre versões, deve possuir testes de contrato.

## 5. Core API

Comunicação entre Core, Runtime, Providers, Plugins oficiais. Possui controle rígido.

## 6. Public API

Consumidores: Dashboard, CLI, ferramentas externas, automações. Obrigatório: documentação OpenAPI, versionamento, compatibilidade.

## 7. Plugin API

Interface oficial para extensões. Regras: compatibilidade por versão do SDK, mudanças controladas, deprecated policy.

## 8. Versionamento

Padrão: Semantic Versioning (MAJOR.MINOR.PATCH).

## 9. Breaking Changes

Mudanças incompatíveis exigem: nova versão major, ADR, migração documentada. Permitido: adicionar campo opcional. Não permitido: remover campo existente.

## 10. API REST

Padrão: HTTP REST, OpenAPI 3.1, Base: /api/v1/

## 11. Recursos Principais

Endpoints oficiais:
- Runtime: GET /api/v1/runtime, POST /api/v1/runtime/start, POST /api/v1/runtime/stop, POST /api/v1/runtime/reload
- Configuration: GET /api/v1/config, PUT /api/v1/config, POST /api/v1/config/validate
- VPN: GET /api/v1/vpn/status, POST /api/v1/vpn/connect, POST /api/v1/vpn/disconnect
- Tor: GET /api/v1/tor/status, POST /api/v1/tor/new-circuit
- Plugins: GET /api/v1/plugins, POST /api/v1/plugins/install, DELETE /api/v1/plugins/{id}

## 12. Resposta Padronizada

{
 "success":true,
 "data":{},
 "error":null,
 "metadata":{}
}

## 13. Erros

{
 "success":false,
 "error":{
    "code":"VPN-001",
    "message":"Connection failed"
 }
}

## 14. Códigos de Erro

Formato: MODULE-NUMBER. Exemplo: CFG-001, VPN-002, TOR-003.

## 15. Autenticação

Métodos suportados: Local Authentication (inicial). Futuro: OAuth, API Keys, certificados.

## 16. Autorização

Modelo: Permission Based Access Control. Exemplo: runtime.read, runtime.write, plugin.install.

## 17. Rate Limiting

APIs públicas devem possuir limite. Exemplo: api.rate_limit.enabled: true, requests: 100/min.

## 18. Idempotência

Operações críticas devem aceitar Idempotency-Key.

## 19. OpenAPI

Toda API pública deve possuir openapi.yaml em /docs/api/.

## 20. SDKs

SDKs oficiais futuros: Go SDK, TypeScript SDK, Python SDK.

## 21. CLI

A CLI utiliza a mesma API pública. ghost CLI -> API Client -> Ghost Daemon.

## 22. Dashboard

O Dashboard nunca acessa módulos internos. Fluxo: Dashboard -> REST API -> Core.

## 23. Eventos

Eventos seguem versionamento próprio: { "type":"VPNConnected", "version":"1.0" }.

## 24. Depreciação

Processo: Announce -> Deprecated -> Migration Period -> Remove.

## 25. Compatibilidade

Objetivo: clientes antigos devem continuar funcionando dentro da mesma versão major.

## 26. Testes de API

Obrigatórios: schema validation, contract tests, compatibility tests, authentication tests, authorization tests.

## 27. Documentação

Toda API deve possuir: descrição, exemplos, erros possíveis, permissões, versão.

## 28. Mudanças Arquiteturais

Alterações em APIs públicas exigem: ADR, atualização OpenAPI, atualização dos testes, changelog.

## 29. Segurança

A API nunca deve expor: secrets, chaves privadas, dados sensíveis, informações internas desnecessárias.

## 30. Definition of Done

Este documento será considerado concluído quando:
- todas as APIs possuírem governança;
- versionamento estiver definido;
- breaking changes forem controladas;
- CLI e Dashboard puderem evoluir sem quebrar o Core;
- integrações externas tiverem contratos estáveis.
