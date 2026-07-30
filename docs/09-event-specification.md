# GhostStack — Event Specification

Versão: 1.0.0
Status: Normativo
Dependências: PRD, System Architecture, Threat Model, Runtime Specification, ADR-001 Runtime Abstraction, Domain Model, Configuration Specification, Contracts Specification, OpenAPI Specification

## 1. Objetivo

Definir o sistema oficial de eventos do GhostStack. Nenhum módulo poderá emitir eventos fora desta especificação.

## 2. Filosofia

O GhostStack utiliza arquitetura orientada a eventos. Os módulos nunca devem depender diretamente entre si quando um evento puder representar a comunicação.

## 3. Componentes

Todos publicam eventos: Gateway, VPN, Tor, DNS, Firewall, Runtime, Plugins, Dashboard, CLI.

## 4. Event Bus

Existe apenas um Event Bus. Responsabilidades: Publish, Subscribe, Dispatch, Retry, Persistência (quando configurada).

## 5. Fluxo

Module -> Publish -> EventBus -> Subscribers

## 6. Estrutura do Evento

{
  "id": "uuid-v7",
  "type": "RuntimeStarted",
  "timestamp": "2026-08-01T12:00:00Z",
  "source": "runtime",
  "version": "1.0",
  "payload": {},
  "metadata": {}
}

## 7. Campos Obrigatórios

id, type, timestamp, source, version, payload, metadata.

## 8. Metadata

Pode conter: requestId, correlationId, traceId, user, profile.

## 9. Categorias

Runtime, Gateway, VPN, Tor, DNS, Firewall, Configuration, Plugin, Security, System.

## 10. Runtime Events

RuntimeInitializing, RuntimeStarted, RuntimeReady, RuntimeReloaded, RuntimeStopping, RuntimeStopped, RuntimeFailed

## 11. Gateway Events

GatewayInitialized, GatewayReloaded, GatewayUpdated, RouteChanged, OutboundChanged, InboundChanged

## 12. VPN Events

VPNConnecting, VPNConnected, VPNDisconnected, VPNReconnecting, VPNFailed

## 13. Tor Events

TorBootstrapping, TorReady, TorCircuitRotated, TorStopped, TorFailed

## 14. DNS Events

DNSReloaded, DNSCacheFlushed, DNSLeakDetected

## 15. Firewall Events

FirewallApplied, FirewallReloaded, KillSwitchEnabled, KillSwitchDisabled, KillSwitchActivated

## 16. Configuration Events

ConfigurationLoaded, ConfigurationValidated, ConfigurationReloaded, ConfigurationMigrated, ProfileApplied

## 17. Plugin Events

PluginInstalled, PluginEnabled, PluginDisabled, PluginRemoved, PluginFailed

## 18. Security Events

SecurityWarning, SecurityViolation, SecretAccessed, CertificateExpired

## 19. System Events

ShutdownRequested, RestartRequested, UpdateAvailable

## 20. Health Events

HealthHealthy, HealthDegraded, HealthFailed

## 21. Severidade

Trace, Debug, Info, Warning, Error, Critical.

## 22. Garantia de Entrega

Por padrão: At Least Once. Eventos críticos deverão possuir confirmação de processamento.

## 23. Ordenação

A ordem é garantida por origem. Não há garantia de ordenação entre módulos distintos.

## 24. Persistência

Configurável:
events:
  persistence: true
  retention: 7d

## 25. Replay

Suportado: ghost events replay. Pode ser utilizado para depuração e auditoria.

## 26. Assinaturas

Todo módulo pode assinar eventos. Exemplo: Gateway subscribe VPNConnected.

## 27. Filtros

Permitidos: origem, tipo, severidade, intervalo de tempo, profile, requestId.

## 28. SSE

Eventos públicos: GET /events
Formato:
event: RuntimeStarted
data: {...}

## 29. WebSocket

Opcional: /ws/events. Mesmo modelo de eventos.

## 30. Plugins

Plugins podem publicar e consumir eventos. Nunca poderão sobrescrever eventos oficiais.

## 31. Auditoria

Eventos persistidos para auditoria nunca poderão ser alterados. Exclusão por políticas de retenção deve ocorrer apenas de forma controlada e auditável.

## 32. Idempotência

Consumidores devem tratar eventos repetidos sem efeitos colaterais. O identificador (id) deve ser utilizado para deduplicação.

## 33. Retentativas

Quando um consumidor falhar: Publish -> Retry -> Retry -> Dead Letter Queue. Após o número configurado de tentativas, o evento deve ser encaminhado para uma fila de erros (DLQ), quando habilitada.

## 34. Dead Letter Queue

Eventos não processados poderão ser armazenados. A DLQ permite: inspeção, reprocessamento, auditoria.

## 35. Versionamento

Eventos possuem versão própria. Mudanças incompatíveis exigem nova versão do evento.

## 36. Compatibilidade

Novos campos podem ser adicionados. Campos existentes nunca poderão mudar de significado dentro da mesma versão major.

## 37. Segurança

Eventos nunca poderão conter: senhas, chaves privadas, tokens de autenticação, conteúdo do tráfego de rede, segredos em texto puro. Quando necessário, utilizar referências ou identificadores.

## 38. Observabilidade

Cada evento deverá registrar: timestamp, duração do processamento (quando aplicável), origem, destino(s), resultado do processamento.

## 39. Testes

Cada evento oficial deverá possuir: teste de serialização, teste de desserialização, teste de compatibilidade, teste de idempotência.

## Definition of Done

Este documento será considerado concluído quando:
- Todos os eventos oficiais estiverem especificados.
- O formato dos eventos for único e consistente.
- As regras de entrega, ordenação e persistência estiverem definidas.
- Os contratos de eventos puderem ser implementados sem ambiguidades.
- O Event Bus puder ser desenvolvido exclusivamente com base nesta especificação.
