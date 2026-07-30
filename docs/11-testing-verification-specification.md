# GhostStack — Testing & Verification Specification

Versão: 1.0.0
Status: Normativo

Dependências: PRD, Architecture, Threat Model, Runtime Specification, ADR-001, Domain Model, Configuration Specification, Contracts Specification, OpenAPI Specification, Event Specification, Plugin SDK Specification

## 1. Objetivo

Definir a estratégia oficial de testes do GhostStack. Nenhum componente poderá ser considerado concluído sem atender aos critérios definidos neste documento.

## 2. Filosofia

O GhostStack adota Specification-Driven Development (SDD), Test-Driven Development (TDD), Contract-First Development e Continuous Verification.

Ordem obrigatória:
Specification -> Contracts -> Tests -> Implementation -> Verification -> Refactor.

Implementações sem testes correspondentes não são aceitas.

## 3. Pirâmide de Testes

| Tipo | Percentual |
|------|------------|
| Unit | 70% |
| Contract | 15% |
| Integration | 10% |
| End-to-End | 5% |

## 4. Categorias

Todo módulo deverá possuir: Unit Tests, Contract Tests, Integration Tests, End-to-End Tests, Security Tests, Performance Tests, Regression Tests.

## 5. Unit Tests

Devem validar regras de negócio, validações, funções puras, transformação de dados e serialização. Nunca dependem de rede, banco, arquivos ou processos.

## 6. Contract Tests

Todos os contratos definidos no Documento 07 deverão possuir testes. Verificações mínimas: tipos, erros, compatibilidade, comportamento esperado.

## 7. Integration Tests

Devem validar integração entre módulos oficiais. Exemplos: Runtime -> Gateway -> Firewall e Configuration -> Runtime -> DNS.

## 8. End-to-End

Fluxos completos. Exemplo: Load Config -> Validate -> Start Runtime -> VPN -> Gateway -> Health.

## 9. Security Tests

Devem validar vazamentos DNS, Kill Switch, isolamento, sandbox, permissões, plugins, autenticação e autorização.

## 10. Performance Tests

Devem medir tempo de inicialização, tempo de reload, uso de memória, uso de CPU, latência da API e throughput do Event Bus.

## 11. Stress Tests

Executar milhares de eventos, centenas de reloads, reconexões contínuas e mudanças frequentes de perfil. Objetivo: garantir estabilidade sob carga.

## 12. Chaos Tests

Simulações obrigatórias: VPN indisponível, Tor indisponível, DNS falhando, Gateway reiniciado, Plugin com falha, disco cheio e falta de memória. O sistema deverá entrar em estado degradado de forma controlada.

## 13. Regression Tests

Todo bug corrigido deverá gerar um teste de regressão permanente.

## 14. Testes de Configuração

Cobertura mínima: parser YAML, merge, includes, profiles, herança, variáveis de ambiente, secrets e migração.

## 15. Testes de Eventos

Todos os eventos deverão validar serialização, desserialização, ordem, idempotência e compatibilidade.

## 16. Testes da API

Validar códigos HTTP, payloads, autenticação, paginação, filtros, erros e schemas OpenAPI.

## 17. Testes de Plugins

Todo plugin oficial deverá possuir testes de: contrato, permissões, manifesto, sandbox, atualização e carregamento.

## 18. Cobertura

| Categoria | Mínimo |
|-----------|--------|
| Core | 95% |
| Runtime | 95% |
| Gateway | 90% |
| API | 90% |
| Plugins | 90% |
| CLI | 85% |

Cobertura não substitui qualidade dos testes.

## 19. Testes Automatizados

Toda Pull Request deverá executar: Lint -> Static Analysis -> Unit -> Contract -> Integration -> Security -> Performance Smoke -> Build. Qualquer falha bloqueia o merge.

## 20. Benchmarks

Benchmarks oficiais: startup, reload, troca de perfil, troca de VPN, rotação Tor e processamento de eventos. Resultados devem ser comparáveis entre versões.

## 21. Compatibilidade

Toda versão deverá validar: configuração anterior, API anterior, eventos anteriores e plugins compatíveis.

## 22. Ferramentas

Ferramentas recomendadas: Go testing, Testcontainers, Mock Servers, OpenAPI Validator, JSON Schema Validator, Benchmark Harness.

## 23. Critérios de Aprovação

Um módulo somente poderá ser considerado concluído quando: todos os testes passarem, cobertura mínima for atingida, análise estática não encontrar erros críticos, testes de contrato forem aprovados e regressões inexistentes.

## 24. Pipeline de Qualidade

Fluxo obrigatório:
Specification -> Contracts -> Tests -> Implementation -> Static Analysis -> Security Scan -> Performance -> Release Candidate.

## 25. Requisitos Normativos

Todo componente deverá: possuir testes automatizados, ser reproduzível, ser determinístico e executar em CI sem intervenção manual.

## 26. Definition of Done

Este documento será considerado concluído quando:
- todas as categorias de testes estiverem formalmente definidas;
- a cobertura mínima estiver especificada;
- o pipeline de validação estiver documentado;
- qualquer implementação puder ser verificada exclusivamente com base nesta especificação.
