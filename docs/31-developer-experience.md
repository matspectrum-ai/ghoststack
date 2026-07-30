# GhostStack — Developer Experience Specification

Versão: 1.0.0
Status: Normativo
Tipo: Engineering Workflow Specification
Dependências: System Architecture, Repository & Contribution Specification, Build System & CI/CD Specification, Testing Strategy & Quality Assurance Specification, API Governance & Versioning Specification, Plugin System & Extension Architecture Specification

## 1. Objetivo

Definir o ambiente oficial de desenvolvimento do GhostStack.

## 2. Filosofia

O ambiente deve ser: simples de iniciar, consistente, automatizado, reproduzível, amigável para humanos e IA.

## 3. Developer Workflow

Fluxo oficial: Clone Repository -> Install Dependencies -> Run Environment Check -> Create Feature Branch -> Implement -> Test -> Review -> Merge.

## 4. Requisitos de Ambiente

Backend: Go >= 1.23
Frontend: Node.js LTS, npm/pnpm
Ferramentas: Git, Make, Docker/Podman (opcional), Editor compatível.

## 5. Primeiro Setup

Comando: git clone https://github.com/ghoststack/ghoststack && cd ghoststack && make setup. Executa: Check Environment -> Install Dependencies -> Generate Files -> Run Validation.

## 6. Estrutura para Desenvolvedor

ghoststack/ com cmd/, internal/, pkg/, plugins/, tests/, docs/, scripts/, tools/, Makefile.

## 7. Comandos Oficiais

make build, make test, make lint, make format, make release.

## 8. Code Style

Princípios: código explícito, nomes descritivos, baixo acoplamento, funções pequenas, erros tratados.

## 9. Regras Go

Obrigatório: go fmt, go vet, tratamento explícito de erros. Evitar panic() sem justificativa.

## 10. Regras TypeScript

Obrigatório: TypeScript strict mode, tipos explícitos em contratos, evitar any.

## 11. Arquitetura de Código

Separação: Domain -> Application -> Infrastructure -> Interface. O Core não deve conhecer: CLI, Dashboard, providers externos.

## 12. Desenvolvimento Orientado a Contratos

Toda nova feature deve possuir: Specification, Contract, Tests, Implementation, Review. Fluxo: SPEC -> API Contract -> Tests -> Code -> Refactor.

## 13. Documentação Obrigatória

Toda feature deve atualizar: docs, contratos, exemplos, changelog.

## 14. ADR (Architecture Decision Record)

Decisões importantes devem possuir ADR. Formato: Context, Decision, Alternatives, Consequences.

## 15. Ambiente Local

O projeto deve possuir .env.example. Nunca .env no Git.

## 16. Configuração Local

make dev inicializa: Ghost Core, API, Dashboard, Development Tools.

## 17. Hot Reload

Suporte: Backend rebuild automático, Frontend Vite HMR.

## 18. Debug

Modo debug: ghost --debug. Logs adicionais: logging.level: debug.

## 19. Observabilidade para Desenvolvimento

Disponível: logs estruturados, métricas, tracing local.

## 20. Desenvolvimento com IA

O projeto deve ser AI-friendly. Requisitos: documentação próxima do código, contratos explícitos, schemas, testes claros.

## 21. Context Files

O repositório deve possuir .ai/ com architecture.md, conventions.md, contracts/, decisions/.

## 22. Prompt Engineering para Agentes

Cada agente deve receber: objetivo, restrições, arquivos relevantes, testes esperados. Modelo: Task, Context, Constraints, Acceptance Criteria, Tests.

## 23. Branch Naming

Padrão: feature/nome, bugfix/nome, security/nome, refator/nome.

## 24. Commit Convention

Formato: TYPE: description. Exemplos: feat: add tor provider, fix: repair dns routing, security: validate certificates.

## 25. Code Review

Todo PR deve verificar: arquitetura, segurança, testes, performance, documentação.

## 26. Dependências Externas

Toda dependência nova deve justificar: motivo, alternativa, impacto.

## 27. Ferramentas de Automação

Scripts em scripts/ com setup, test, release, security-check.

## 28. Ambiente Reproduzível

Um novo desenvolvedor deve conseguir: make setup, make test, make build sem configuração manual adicional.

## 29. Suporte a Sistemas

Ambiente deve funcionar em: Linux, Windows, macOS.

## 30. Definition of Done

Este documento será considerado concluído quando:
- um novo desenvolvedor conseguir iniciar o projeto rapidamente;
- agentes de IA conseguirem navegar pelo código;
- padrões estiverem documentados;
- builds e testes forem reproduzíveis;
- mudanças seguirem contratos definidos.
