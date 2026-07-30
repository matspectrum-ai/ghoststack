# GhostStack — Repository & Contribution Specification

Versão: 1.0.0
Status: Normativo
Tipo: Software Governance Specification
Dependências: Developer Experience Specification, Build System & CI/CD Specification, Testing Strategy & Quality Assurance Specification, API Governance & Versioning Specification, System Architecture, Security Architecture

## 1. Objetivo

Definir as regras oficiais de organização, manutenção e contribuição do repositório GhostStack.

## 2. Filosofia

O repositório deve ser: organizado, auditável, previsível, seguro, preparado para colaboração humana e agentes de IA.

## 3. Estrutura Oficial do Repositório

ghoststack/ com cmd/, internal/, pkg/, plugins/, dashboard/, sdk/, schemas/, configs/, docs/, examples/, scripts/, tests/, tools/, .github/, Makefile, README.md, CHANGELOG.md, LICENSE, SECURITY.md.

## 4. Responsabilidade dos Diretórios

cmd/: executáveis ghost e ghost-daemon
internal/: código privado do Core, não pode ser importado externamente
pkg/: bibliotecas públicas (SDK, contratos, APIs)
plugins/: extensões oficiais
dashboard/: frontend
sdk/: ferramentas para criação de plugins
schemas/: contratos (JSON Schema, OpenAPI, Configuration Schema)
docs/: documentação técnica

## 5. Arquivos Obrigatórios

Todo repositório deve possuir: README.md, LICENSE, CHANGELOG.md, SECURITY.md, CONTRIBUTING.md.

## 6. Branch Strategy

Modelo: main -> develop -> feature/*.

## 7. Branch Main

Sempre compilável, sempre testada, pronta para release.

## 8. Branch Develop

Uso: integração de funcionalidades, próxima versão.

## 9. Feature Branch

Formato: feature/<nome>. Exemplo: feature/tor-provider.

## 10. Bugfix Branch

Formato: bugfix/<nome>.

## 11. Security Branch

Formato: security/<nome>. Uso: correções críticas, vulnerabilidades.

## 12. Commit Convention

Padrão: TYPE: description. Tipos: feat, fix, security, refactor, docs, test, build, ci.

## 13. Pull Request

Todo PR deve conter: Descrição (problema, solução, impacto), Alterações (Added, Changed, Removed), Testes (Tests Added, Tests Passed), Risco (Low, Medium, High, Critical).

## 14. Pull Request Template

## Summary
## Motivation
## Changes
## Tests
## Security Impact
## Breaking Changes

## 15. Code Review

Revisores devem avaliar: arquitetura, segurança, testes, performance, compatibilidade.

## 16. Aprovação

Regras: mudanças comuns 1 reviewer, mudanças críticas 2 reviewers, mudanças de segurança security review obrigatório.

## 17. CODEOWNERS

Áreas possuem responsáveis. Exemplo: /internal/security @security-team, /dashboard @frontend-team.

## 18. Issues

Tipos: Feature, Bug, Security, Documentation, Performance, Question.

## 19. Feature Request

Obrigatório conter: problema, motivação, proposta, impacto.

## 20. Bug Report

Obrigatório: sistema operacional, versão, passos para reproduzir, logs, comportamento esperado.

## 21. Security Report

Nunca abrir publicamente. Fluxo: Researcher -> Security Channel -> Validation -> Patch -> Disclosure.

## 22. Documentação

Mudanças de comportamento exigem atualização de: docs, schemas, exemplos.

## 23. Versionamento

O projeto utiliza Semantic Versioning. Formato: MAJOR.MINOR.PATCH.

## 24. Breaking Changes

Mudanças incompatíveis exigem: novo MAJOR, migration guide, ADR.

## 25. Release Management

Fluxo: Develop -> Release Branch -> Testing -> Approval -> Main -> Tag.

## 26. Tags

Formato: v1.0.0.

## 27. Changelog

Toda release deve registrar: Added, Changed, Fixed, Security, Deprecated.

## 28. Dependabot / Dependency Management

Dependências devem possuir: atualização automática, análise de vulnerabilidade, revisão.

## 29. Regras para Agentes de IA

Agentes trabalhando no repositório devem, antes de alterar: ler documentação relacionada, entender contratos, executar testes existentes, criar novos testes quando necessário.

## 30. Proibição de Alterações Diretas

Agentes e colaboradores não devem: modificar contratos sem ADR, remover testes para passar CI, alterar segurança sem revisão.

## 31. Auditoria

Todas mudanças importantes devem possuir: autor, commit, revisão, motivação.

## 32. Definition of Done

Este documento será considerado concluído quando:
- estrutura do repositório estiver definida;
- contribuição estiver padronizada;
- revisões forem obrigatórias;
- histórico permanecer auditável;
- humanos e agentes puderem colaborar com segurança.
