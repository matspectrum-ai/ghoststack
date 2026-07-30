# GhostStack — Contracts Specification

Versão: 1.0.0
Status: Normativo
Dependências: PRD, System Architecture, Threat Model, Runtime Specification, ADR-001 Runtime Abstraction, Domain Model, Configuration Specification

## 1. Objetivo

Este documento define todos os contratos internos do GhostStack. Os contratos representam a única forma oficial de comunicação entre módulos. Nenhum módulo poderá acessar diretamente outro módulo sem utilizar um contrato definido neste documento.

## 2. Princípios

Todos os contratos deverão ser: fortemente tipados, determinísticos, síncronos ou assíncronos explicitamente, independentes da implementação, compatíveis com testes automatizados.

## 3. Convenções

Interfaces: I<Name>
DTOs: <Name>Request, <Name>Response
Eventos: <Name>Event
Providers: <Name>Provider

## 4. Contrato IRuntimeProvider

Responsável pela execução dos componentes.

interface IRuntimeProvider {
  initialize(): Promise<void>
  start(): Promise<void>
  stop(): Promise<void>
  restart(): Promise<void>
  reload(): Promise<void>
  status(): Promise<RuntimeStatus>
  health(): Promise<HealthStatus>
  version(): Promise<string>
}

## 5. Contrato IGateway

interface IGateway {
  initialize(): Promise<void>
  apply(config: GatewayConfig): Promise<void>
  reload(): Promise<void>
  shutdown(): Promise<void>
  status(): Promise<GatewayStatus>
  routes(): Promise<Route[]>
}

## 6. Contrato IVPNProvider

interface IVPNProvider {
  connect(profile: string): Promise<void>
  disconnect(): Promise<void>
  reconnect(): Promise<void>
  status(): Promise<VpnStatus>
  latency(): Promise<number>
}

## 7. Contrato ITorProvider

interface ITorProvider {
  start(): Promise<void>
  stop(): Promise<void>
  rotateCircuit(): Promise<void>
  bootstrap(): Promise<number>
  status(): Promise<TorStatus>
}

## 8. Contrato IDNSProvider

interface IDNSProvider {
  initialize(): Promise<void>
  reload(): Promise<void>
  flushCache(): Promise<void>
  status(): Promise<DnsStatus>
}

## 9. Contrato IFirewallProvider

interface IFirewallProvider {
  apply(): Promise<void>
  clear(): Promise<void>
  enableKillSwitch(): Promise<void>
  disableKillSwitch(): Promise<void>
  status(): Promise<FirewallStatus>
}

## 10. Contrato IProxyProvider

interface IProxyProvider {
  connect(): Promise<void>
  disconnect(): Promise<void>
  health(): Promise<HealthStatus>
}

## 11. Contrato IConfigurationProvider

interface IConfigurationProvider {
  load(): Promise<Configuration>
  validate(): Promise<ValidationResult>
  reload(): Promise<void>
  save(): Promise<void>
}

## 12. Contrato IPlugin

Todo plugin implementará:

interface IPlugin {
  initialize(): Promise<void>
  enable(): Promise<void>
  disable(): Promise<void>
  unload(): Promise<void>
  manifest(): PluginManifest
}

## 13. Contrato IHealthCheck

interface IHealthCheck {
  check(): Promise<HealthStatus>
}

## 14. Contrato ISecretsProvider

interface ISecretsProvider {
  get(key: string): Promise<string>
  exists(key: string): Promise<boolean>
}

## 15. Contrato IEventBus

interface IEventBus {
  publish(event: DomainEvent): Promise<void>
  subscribe(event: string, handler: EventHandler): Promise<void>
}

## 16. Contrato ILogger

interface ILogger {
  trace()
  debug()
  info()
  warn()
  error()
  fatal()
}

## 17. DTO RuntimeStatus

interface RuntimeStatus {
  state: RuntimeState
  provider: RuntimeProvider
  uptime: number
}

## 18. DTO GatewayStatus

interface GatewayStatus {
  backend: string
  mode: RuntimeMode
  ready: boolean
}

## 19. DTO HealthStatus

interface HealthStatus {
  status: "Healthy" | "Degraded" | "Failed" | "Unknown"
  message: string
}

## 20. DTO ValidationResult

interface ValidationResult {
  valid: boolean
  errors: ValidationError[]
}

## 21. DTO ValidationError

interface ValidationError {
  code: string
  field: string
  message: string
}

## 22. Plugin Manifest

interface PluginManifest {
  id: string
  version: string
  author: string
  description: string
  capabilities: string[]
  permissions: string[]
}

## 23. Eventos de Domínio

Eventos oficiais:
- RuntimeStarted
- RuntimeStopped
- GatewayReloaded
- ModeChanged
- VPNConnected
- VPNDisconnected
- TorStarted
- TorStopped
- ProxyConnected
- ConfigurationReloaded
- PluginInstalled
- PluginRemoved
- HealthChanged

## 24. Runtime Provider

Implementações oficiais:
- NativeRuntimeProvider
- DockerRuntimeProvider
- PodmanRuntimeProvider

Implementações futuras:
- ContainerdRuntimeProvider
- KubernetesRuntimeProvider
- LXCRuntimeProvider

## 25. Gateway Provider

Implementação oficial:
- SingBoxGateway

Implementações futuras:
- XrayGateway
- HysteriaGateway

## 26. VPN Provider

Implementações oficiais:
- WireGuardProvider
- OpenVPNProvider

Implementações futuras:
- TailscaleProvider
- AmneziaProvider

## 27. DNS Provider

Implementações oficiais:
- UnboundProvider

Futuros:
- CoreDNSProvider

## 28. Firewall Provider

Implementação oficial:
- NftablesProvider

Futuro:
- IPTablesProvider

## 29. Regras

Todo Provider deverá:
- implementar sua interface integralmente;
- retornar erros tipados;
- ser substituível sem alterar consumidores;
- ser testável isoladamente.

## 30. Contrato de Erros

Todos os erros deverão possuir:

interface GhostError {
  code: string
  module: string
  message: string
  cause?: string
  recommendation?: string
}

## 31. Contrato de Versionamento

Todos os providers deverão expor:

interface Versioned {
  version(): Promise<string>
}

## 32. Contrato de Observabilidade

Todo módulo deverá expor: versão, health, status, métricas, eventos.

## 33. Compatibilidade

Interfaces públicas nunca poderão sofrer alterações incompatíveis dentro da mesma versão major. Mudanças incompatíveis exigem incremento de versão major.

## 34. Requisitos Normativos

Todos os contratos deverão:
- ser independentes da implementação;
- utilizar DTOs imutáveis;
- retornar Promise;
- ser documentados;
- possuir testes de contrato (contract tests).

## Definition of Done

Este documento será considerado concluído quando:
- todos os contratos estiverem formalmente definidos;
- cada módulo possuir uma interface oficial;
- DTOs estiverem especificados;
- eventos estiverem identificados;
- providers puderem ser implementados sem ambiguidade.
