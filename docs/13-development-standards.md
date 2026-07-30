# GhostStack — Development Standards Specification

Versão: 1.0.0
Status: Normativo
Dependências: PRD, System Architecture, Threat Model & Security Architecture, Runtime Specification, ADR-001 Runtime Abstraction, Domain Model, Configuration Specification, Contracts Specification, OpenAPI Specification, Event Specification, Plugin SDK Specification, Testing & Verification Specification, Implementation Roadmap

## 1. Objetivo

Definir os padrões obrigatórios de desenvolvimento do GhostStack. Nenhum código poderá ser escrito fora destas regras.

## 2. Princípios

O desenvolvimento seguirá: Specification Driven Development (SDD), Test Driven Development (TDD), Domain Driven Design (DDD), Clean Architecture, Contract First, Security By Design.

## 3. Stack Oficial

- Core: Go
- Interface Web: TypeScript, React
- CLI: Go
- Configuração: YAML 1.2, JSON Schema
- API: REST, OpenAPI 3.1
- Banco Local: SQLite

## 4. Arquitetura de Código

Estrutura oficial:
```
ghoststack/
├── cmd/
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── runtime/
│   ├── networking/
│   ├── security/
│   ├── plugins/
│   └── events/
├── api/
├── sdk/
├── web/
├── tests/
├── docs/
└── scripts/
```

## 5. Clean Architecture

Dependências obrigatórias: Presentation -> Application -> Domain -> Infrastructure.

## 6. Regras de Dependência

O domínio nunca pode importar: banco, HTTP, filesystem, Docker, APIs externas. Nunca permitido: Domain -> Infrastructure.

## 7. Organização dos Módulos

Cada módulo deve possuir: module/ com domain/, application/, infrastructure/, interfaces/ e tests/.

## 8. Interfaces

Interfaces devem estar próximas do consumidor, possuir nomes explícitos e evitar abstrações genéricas.

## 9. Nomenclatura

- Arquivos: snake_case (ex: runtime_manager.go)
- Interfaces: prefixo I (ex: IRuntimeProvider)
- Eventos: past tense (ex: VPNConnected, não ConnectVPN)

## 10. Erros

Erros devem ser tipados. Nunca utilizar somente errors.New("failed").

## 11. Logs

Logs devem possuir timestamp, módulo, nível e código.

## 12. Configuração

Nenhum módulo lê ghost.yaml diretamente. Fluxo obrigatório: Config Provider -> Validated Config -> Module.

## 13. Testes

Todo código novo deve seguir: Red -> Green -> Refactor. Antes da implementação: obrigatório criar contrato, teste e caso de erro.

## 14. Cobertura

- Core: 95%
- Infraestrutura: 90%
- API: 90%
- Plugins: 90%

## 15. Git Workflow

Branches: main, develop, feature/*, bugfix/*, release/*.

## 16. Commits

Formato obrigatório: type(scope): description. Tipos: feat, fix, docs, test, refactor, security, build, ci.

## 17. Pull Requests

Toda PR deve conter: objetivo, documento relacionado, testes, impacto, riscos.

## 18. ADR obrigatório

Criar ADR quando houver mudança: arquitetura, banco, protocolo, segurança, contrato público.

## 19. CI/CD

Pipeline: Commit -> Lint -> Unit Tests -> Contract Tests -> Security Scan -> Build -> Integration Tests -> Release.

## 20. Static Analysis

Obrigatório: análise de código, análise de dependências, verificação de vulnerabilidades.

## 21. Segurança de Dependências

Toda dependência externa deve possuir: licença compatível, manutenção ativa, análise de vulnerabilidades.

## 22. Documentação

Todo módulo deve possuir: README, arquitetura, contratos, testes, exemplos.

## 23. Código Gerado por IA

Código gerado por agentes deverá: seguir contratos existentes, possuir testes, passar CI, não criar novas abstrações sem ADR.

## 24. Regras para Agentes IA

Agentes deverão, antes de modificar código: ler documentação relacionada, identificar contratos existentes, criar ou atualizar testes, implementar alteração mínima, executar validações.

## 25. Proibido

Não permitido: código sem teste, alteração de contrato sem ADR, segredo no código, dependência não documentada, bypass de segurança.

## 26. Release

Toda release deve conter: changelog, versão, documentação atualizada, migrações necessárias, notas de segurança.

## 27. Definition of Done

Este documento estará concluído quando: todos os padrões de desenvolvimento estiverem definidos; novos colaboradores puderem contribuir sem decisões implícitas; agentes de IA puderem seguir regras claras; o código permanecer consistente entre módulos.
