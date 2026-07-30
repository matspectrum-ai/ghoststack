# GhostStack — Provider Integration Specification

Versão: 1.0.0
Status: Normativo
Tipo: Integration Architecture Specification
Dependências: PRD, System Architecture, Runtime Specification, Networking & Gateway Specification, API Governance & Versioning Specification, Plugin System & Extension Architecture Specification, Security Operations Specification, Data Model & Storage Specification

## 1. Objetivo

Definir a arquitetura oficial de integração entre o GhostStack e provedores externos.

## 2. Filosofia

Ghost Core -> Provider Interface -> Adapter -> External Provider.

## 3. Objetivos

Permitir: troca de provedores, múltiplas implementações, testes isolados, evolução independente.

## 4. Tipos de Providers

VPN, Proxy, Tor, DNS, Gateway, Firewall, Storage, Identity.

## 5. Provider Interface

Contrato base com Name, Version, Initialize, Start, Stop, Status, Health.

## 6. Provider Lifecycle

Discovered -> Validated -> Initialized -> Running -> Stopped -> Removed.

## 7. Provider Manifest

Todo provider deve possuir provider.yaml com name, version, type, ghoststack.min_version, permissions.

## 8. VPN Providers

Interface VPNProvider com Connect, Disconnect, Reconnect, GetStatus. Suporte inicial: WireGuard, OpenVPN.

## 9. VPN Provider Requirements

Deve fornecer: estado da conexão, erro detalhado, métricas, health check.

## 10. Proxy Providers

Interface ProxyProvider com AddProxy, RemoveProxy, Test, Select. Tipos suportados: HTTP, HTTPS, SOCKS5.

## 11. Proxy Pool

Suporte: Single, Round Robin, Random, Latency Based, Priority.

## 12. Tor Provider

Interface TorProvider com Start, Stop, NewCircuit, Status. Responsabilidades: iniciar daemon, validar bootstrap, controlar SOCKS, monitorar estado.

## 13. DNS Providers

Interface DNSProvider com Resolve, Configure, Test. Suporte: Unbound, DoT, DoH, DNSSEC.

## 14. Gateway Providers

Interface GatewayProvider com Route, Tunnel, Configure. Implementações: inicial sing-box, futuro custom gateway.

## 15. Firewall Providers

Interface FirewallProvider com Enable, Disable, AddRule, RemoveRule, Status. Backends: Linux nftables/iptables, Windows Windows Filtering Platform, macOS pf.

## 16. Descoberta de Providers

Processo: Scan -> Validate Manifest -> Check Compatibility -> Register. Comando: ghost provider list.

## 17. Registro de Providers

Database: tabela providers com id, name, type, version, status, installed_at.

## 18. Configuração

Exemplo: providers.vpn.active: wireguard, providers.dns.active: unbound.

## 19. Autenticação

Providers podem usar: arquivos, variáveis de ambiente, secret manager. Nunca senha em texto puro.

## 20. Secrets

Fluxo: Provider -> Secret Provider -> Secure Storage.

## 21. Health Check

Todo provider deve responder Health(). Estados: Healthy, Degraded, Failed.

## 22. Fallback

Quando um provider falhar: configurável, registrado, auditável.

## 23. Timeout

Todos providers devem possuir timeout connect: 30s, operation: 10s.

## 24. Retry Policy

Tentativas: retry.attempts: 3, retry.backoff: exponential.

## 25. Atualização de Providers

Processo: Check Version -> Download -> Verify Signature -> Install -> Restart.

## 26. Compatibilidade

Provider declara compatibility.ghoststack.min: 1.0, max: 2.x.

## 27. Segurança

Providers não podem: ignorar firewall, acessar secrets sem permissão, alterar configuração global.

## 28. Logs

Cada provider possui namespace provider.vpn.wireguard.

## 29. Métricas

Providers devem expor: provider_status, provider_latency, provider_errors.

## 30. Testes

Obrigatórios: inicialização, falha, reconexão, timeout, troca de provider, permissões.

## 31. Definition of Done

Este documento será considerado concluído quando:
- providers puderem ser substituídos;
- integrações forem desacopladas;
- falhas forem controladas;
- novos adapters puderem ser adicionados sem alterar o Core.
