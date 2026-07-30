# GhostStack — Domain Model Specification

Versão: 0.1.0
Status: Normativo
Dependências: PRD, System Architecture, Runtime Specification, Threat Model, ADR-001 Runtime Abstraction

## 1. Objetivo

Definir o modelo de domínio oficial do GhostStack. Nenhum contrato poderá introduzir conceitos fora deste modelo.

## 2. Filosofia

O modelo de domínio representa o negócio, não a implementação. Ele independe de Docker, SQLite, REST, gRPC, linguagem ou banco.

## 3. Bounded Contexts

O domínio é dividido em oito contextos: Configuration, Runtime, Networking, Security, Profiles, Plugins, Observability, System. Cada contexto possui seus próprios modelos.

## 4. Configuration

Responsável por: ghost.yaml, merge, validação, profiles, defaults. Entidade raiz: Configuration.

## 5. Runtime

Responsável por: estado, lifecycle, runtime provider, health. Entidade raiz: Runtime.

## 6. Networking

Responsável por: gateway, dns, vpn, tor, proxy, routing. Entidade raiz: Gateway.

## 7. Security

Responsável por: firewall, kill switch, secrets, certificates. Entidade raiz: SecurityPolicy.

## 8. Plugins

Responsável por: plugins, versões, dependências. Entidade raiz: Plugin.

## 9. Observability

Responsável por: logs, métricas, eventos, health. Entidade raiz: Telemetry.

## 10. Entidades

As entidades possuem identidade própria: Runtime (RuntimeId), Gateway (GatewayId), VPN (VpnProfileId), Tor (TorInstanceId), Proxy (ProxyId), Profile (ProfileId), Plugin (PluginId).

## 11. Value Objects

Não possuem identidade: Endpoint (host, port), Duration (5s, 30s, 10m), Version (major, minor, patch), CIDR (10.0.0.0/24), SecretReference (env, file, keyring).

## 12. Agregados

Configuration -> Profiles -> Runtime Defaults
Gateway -> Routes -> Outbounds -> Inbounds
Plugin -> Manifest -> Permissions -> Dependencies
Runtime -> Components -> Health -> Metrics

## 13. Relacionamentos

Configuration -> Runtime -> Gateway -> Networking
Configuration -> Profiles
Runtime -> Health
Plugin -> Capabilities

## 14. Runtime

Possui: Current Mode, Provider, Status, Components.

## 15. Gateway

Possui: Rules, Outbounds, Inbounds, DNS, Policies.

## 16. VPN

Possui: Provider, Profile, Status, Latency.

## 17. Tor

Possui: Circuit, Bootstrap, Status.

## 18. Proxy

Possui: Endpoint, Protocol, Authentication, Priority.

## 19. Firewall

Possui: Backend, Rules, Kill Switch.

## 20. DNS

Possui: Resolver, Cache, DNSSEC, DoH, DoT.

## 21. Runtime Provider

Representa a implementação responsável pela execução dos componentes. Tipos oficiais: Native, Docker, Podman. Futuros: Containerd, Kubernetes, LXC.

## 22. Plugin

Todo plugin deverá possuir: Id, Version, Author, Manifest, Capabilities, Permissions.

## 23. Capabilities

Exemplos: vpn.provider, proxy.provider, dns.provider, gateway.backend, dashboard.widget.

## 24. Permissions

Plugins deverão declarar permissões. Exemplo: Network, Filesystem, Secrets, Runtime, Configuration.

## 25. Health

Estados válidos: Healthy, Degraded, Failed, Unknown.

## 26. Lifecycle

Todos os componentes implementam: Created, Initialized, Ready, Running, Stopping, Stopped, Failed.

## 27. Invariantes

INV-001: Existe apenas um Runtime ativo.
INV-002: Existe apenas um Gateway ativo.
INV-003: Existe apenas um modo ativo.
INV-004: Todo Runtime possui exatamente um Runtime Provider.
INV-005: Toda Configuration válida gera exatamente um Runtime.
INV-006: Todo Plugin possui Manifest.
INV-007: Toda Route pertence a um Gateway.
INV-008: Toda SecurityPolicy pertence a uma Configuration.

## 28. Identificadores

Todos os IDs deverão utilizar UUID v7. Exemplo: 0198c8b7-a42d-7db5-9c77-...

## 29. Versionamento

Toda entidade possui: CreatedAt, UpdatedAt, Version.

## 30. Eventos de Domínio

Exemplos: RuntimeStarted, RuntimeStopped, ModeChanged, VPNConnected, VPNDisconnected, TorBootstrapped, ProxySelected, PluginInstalled, ConfigurationReloaded, HealthChanged. Esses eventos são eventos de domínio, independentes da tecnologia de transporte.

## 31. Estados Persistentes

Persistidos: Configuration, Profiles, Plugins, Histórico, Runtime State.
Não persistidos: sockets, processos, conexões, threads, interfaces temporárias.

## 32. Regras de Negócio

- Apenas um modo pode estar ativo.
- O Runtime Provider deve estar disponível antes do startup.
- Um Profile não pode referenciar outro inexistente.
- Um Plugin não pode solicitar permissões não declaradas.
- Uma Configuration inválida nunca poderá gerar um Runtime.

## 33. Extensibilidade

Novos provedores de VPN, DNS, Proxy, Gateway ou Runtime devem ser adicionados por meio de plugins ou implementações das interfaces oficiais, sem alterar o modelo de domínio.

## Definition of Done

Este documento estará concluído quando:
- Todas as entidades do domínio estiverem identificadas.
- Todos os agregados estiverem definidos.
- Todos os relacionamentos estiverem documentados.
- Todas as invariantes estiverem especificadas.
- Os contratos (REST, gRPC, SDK e banco) puderem ser derivados exclusivamente deste documento.
