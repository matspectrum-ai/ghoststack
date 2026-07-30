# GhostStack — Plugin System & Extension Architecture Specification

Versão: 1.0.0
Status: Normativo
Tipo: Extension Architecture Specification
Dependências: PRD, System Architecture, Domain Model, Security Operations Specification, API Governance & Versioning Specification, Data Model & Storage Specification, Cross Platform Architecture Specification, Observability & Monitoring Specification

## 1. Objetivo

Definir o sistema oficial de plugins e extensões do GhostStack.

## 2. Filosofia

Plugins devem permitir expansão de funcionalidades, integração com serviços externos e criação de módulos personalizados. Sem permitir alteração do Core, acesso irrestrito ao sistema ou quebra do modelo de segurança.

## 3. Princípios

O sistema segue: Least Privilege, Explicit Permissions, Versioned Contracts, Fail Closed, Isolation By Default.

## 4. Arquitetura Geral

Ghost Core -> Plugin Runtime -> Plugin A, Plugin B, Plugin C

## 5. Tipos de Plugins

Três categorias:
- Core Plugins: oficiais, mantidos pelo projeto, maior nível de confiança (ex: vpn-wireguard, tor-controller, dns-unbound, gateway-singbox)
- Community Plugins: criados por terceiros, permissões restritas, assinatura recomendada
- Local Plugins: uso pessoal, desenvolvimento, testes, automações privadas

## 6. Estrutura de Plugin

Todo plugin deve possuir:
```
plugin/
├── manifest.yaml
├── binary/
├── config/
├── assets/
└── README.md
```

## 7. Manifesto

Arquivo obrigatório manifest.yaml com name, version, type, ghoststack.min_version, permissions.

## 8. Identidade do Plugin

Todo plugin possui: ID único, nome, versão, autor, assinatura. Exemplo: plugin_id: ghoststack.plugin.example

## 9. Ciclo de Vida

Estados: Installed -> Loaded -> Initialized -> Running -> Stopped -> Removed

## 10. Eventos do Plugin

Plugins recebem eventos: RUNTIME_STARTED, VPN_CONNECTED, TOR_READY, CONFIG_CHANGED, NETWORK_CHANGED.

## 11. Comunicação com Core

Plugins nunca acessam o Core diretamente. Comunicação: Plugin -> Plugin API -> Core.

## 12. Plugin SDK

SDK oficial: ghoststack-sdk. Linguagens iniciais: Go, TypeScript.

## 13. API do Plugin

PluginContext com Logger(), Config(), Events(), Storage().

## 14. Permissões

Modelo: permissions com network, filesystem, secrets, system. Padrão: deny all.

## 15. Permissões Disponíveis

- Network: criar conexões, usar APIs externas
- Filesystem: leitura/escrita limitada
- Secrets: acessar secrets autorizados
- System: operações privilegiadas, requer aprovação explícita

## 16. Sandbox

Plugins devem executar isolados. Métodos: Linux (namespaces, seccomp, usuário separado), Windows (restricted token, job objects), macOS (sandbox profiles).

## 17. Processo Separado

Modo recomendado: Ghost Core -> Plugin Runtime -> Plugin Process. Falha do plugin não deve derrubar o Core.

## 18. Instalação

Comando: ghost plugin install <name>. Fluxo: Download -> Verify -> Validate Manifest -> Check Permissions -> Install -> Load.

## 19. Remoção

Comando: ghost plugin remove <name>. Fluxo: Stop -> Unload -> Remove Files -> Cleanup Data.

## 20. Atualização

Comando: ghost plugin update. Fluxo: Check Version -> Backup -> Install Update -> Migration -> Restart Plugin.

## 21. Assinatura de Plugins

Plugins oficiais: signature obrigatória. Plugins comunitários: signature recomendada.

## 22. Compatibilidade

Manifesto: compatibility.ghoststack.min: 1.0, max: 2.x. Plugins incompatíveis não iniciam.

## 23. Versionamento

Plugins usam Semantic Versioning. Exemplo: 1.2.0.

## 24. Storage de Plugins

Cada plugin possui isolamento. plugins/example/data/. Não pode acessar dados de outro plugin.

## 25. Logs

Cada plugin possui namespace plugin.example. Logs seguem o padrão com plugin e event.

## 26. Métricas

Plugins podem registrar: métricas próprias, health checks, eventos. Exemplo: plugin_status, plugin_errors.

## 27. Dashboard Integration

Plugins podem fornecer páginas, widgets, informações. Nunca podem modificar autenticação ou substituir componentes críticos.

## 28. Marketplace Futuro

Arquitetura preparada para Ghost Marketplace com catálogo, avaliações, assinatura, distribuição.

## 29. Segurança

Antes de carregar: Manifest Validation -> Permission Review -> Signature Check -> Sandbox Start.

## 30. Falha de Plugin

Se plugin falhar: Detect -> Disable -> Notify -> Keep Core Running.

## 31. Testes

Obrigatórios: carregamento, permissões, isolamento, atualização, remoção, falha controlada.

## 32. Definition of Done

Este documento será considerado concluído quando:
- plugins puderem ser adicionados sem alterar o Core;
- permissões forem explícitas;
- plugins forem isolados;
- compatibilidade for controlada;
- falhas não comprometerem o sistema.
