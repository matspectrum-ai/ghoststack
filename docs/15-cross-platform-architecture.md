# GhostStack — Cross Platform Architecture Specification

Versão: 1.0.0
Status: Normativo
Tipo: Architecture Specification
Dependências: PRD, System Architecture, Threat Model & Security Architecture, Runtime, Networking & Gateway Specification, ADR-001 Runtime Abstraction, Domain Model, Configuration Specification, Contracts Specification, Event Specification, Plugin SDK Specification, Testing & Verification Specification, Implementation Roadmap, Development Standards Specification

## 1. Objetivo

Definir a arquitetura oficial multiplataforma do GhostStack. Diferenças específicas de sistema operacional deverão existir somente nas camadas de infraestrutura.

## 2. Princípio Arquitetural

"Platform specific code must stay at the edge." Código específico de sistema operacional nunca deve contaminar o domínio.

## 3. Sistemas Suportados

Versão inicial:
- Debian/Ubuntu: Oficial
- antiX Linux: Oficial
- Windows 10/11: Oficial
- macOS 13+: Oficial

Arquiteturas:
- Inicial: x86_64
- Futuro: ARM64

## 4. Arquitetura Geral

Ghost CLI -> Ghost Daemon -> Core Runtime -> Network/Security/Runtime -> Platform Layer -> Linux/Windows/macOS

## 5. Core Multiplataforma

O Core deve ser independente do sistema operacional. Responsabilidades: gerenciamento de estado, eventos, configuração, plugins, API, regras de negócio. O Core não pode acessar firewall diretamente, interfaces de rede diretamente, serviços do sistema diretamente.

## 6. Platform Layer

Nova camada obrigatória. Estrutura:
```
platform/
├── linux/
├── windows/
└── darwin/
```
Responsável por implementar contratos do sistema.

## 7. Contratos de Plataforma

Todos os recursos dependentes do SO devem possuir interfaces. Exemplo: FirewallProvider com Enable, Disable, AddRule, RemoveRule. Implementações: LinuxFirewall, WindowsFirewall, MacOSFirewall.

## 8. Gerenciamento de Serviços

- Linux: systemd / LinuxServiceManager
- Windows: Windows Service Manager / WindowsServiceManager
- macOS: launchd / MacOSServiceManager

## 9. Firewall Abstraction

Interface FirewallProvider com Enable, Disable, AddRule, RemoveRule.
- Linux: nftables
- Windows: Windows Filtering Platform
- macOS: pf firewall

## 10. Routing Layer

Interface RouteProvider com AddRoute, RemoveRoute, ListRoutes.
- Linux: ip route, netlink
- Windows: Windows Routing API
- macOS: route, NetworkExtension

## 11. DNS Layer

Interface DNSProvider com SetResolver, RestoreResolver, Status.
- Linux: systemd-resolved, resolv.conf, Unbound
- Windows: Windows DNS Client API
- macOS: scutil, NetworkExtension

## 12. VPN Layer

Interface VPNProvider com Connect, Disconnect, Reconnect, Status. Provider inicial: WireGuard.
- Linux: wireguard-tools
- Windows: WireGuard Tunnel Service
- macOS: WireGuardKit

## 13. Tor Layer

Interface TorProvider com Start, Stop, NewCircuit, Status. Implementação inicial: Tor daemon. Compatível com Linux, Windows, macOS.

## 14. Gateway Layer

Gateway permanece independente. Provider oficial: sing-box. Responsabilidade do GhostStack: gerar configuração, iniciar processo, monitorar estado, aplicar regras.

## 15. Process Management

Interface ProcessManager com Start, Stop, Restart, Status.
- Linux: os/exec, systemd
- Windows: CreateProcess, Service Manager
- macOS: launchd, os/exec

## 16. Filesystem Abstraction

Interface FileSystemProvider com Read, Write, Exists. Objetivo: evitar caminhos fixos.

## 17. Diretórios por Sistema

Linux:
- /etc/ghoststack
- /var/lib/ghoststack
- /var/log/ghoststack

Windows:
- C:\ProgramData\GhostStack
- C:\Program Files\GhostStack

macOS:
- /Library/Application Support/GhostStack
- /Library/Logs/GhostStack

## 18. CLI Universal

Comando ghost em todos os sistemas. Comandos: ghost status, ghost mode tor, ghost vpn connect, ghost config validate.

## 19. Configuração Universal

O mesmo arquivo ghost.yaml deve funcionar nos três sistemas. Diferenças ficam internas.

## 20. Build Matrix

CI obrigatório:
- Linux amd64
- Windows amd64
- macOS amd64
- macOS arm64

## 21. Testes Multiplataforma

Cada release deve executar: Unit Tests em todos os sistemas, Integration Tests por plataforma, Security Tests por plataforma.

## 22. Compatibilidade

O comportamento público deve ser idêntico. Exemplo: ghost status deve retornar o mesmo modelo independente do sistema.

## 23. Plugins

Plugins devem declarar compatibilidade:
```
platforms:
  - linux
  - windows
  - darwin
```
Plugins específicos permitidos: linux.firewall.nftables. Não permitido: modificar Core.

## 24. Segurança

Cada plataforma deve implementar: isolamento, permissões mínimas, proteção de secrets, controle de processos.

## 25. Atualizações

O sistema de atualização deve reconhecer:
- ghoststack-linux-amd64
- ghoststack-windows-amd64.exe
- ghoststack-macos-arm64

## 26. Limitações Conhecidas

Algumas funcionalidades podem variar:
- nftables: Linux apenas
- WFP: Windows apenas
- pf: macOS apenas
- WireGuard: Todos
- Tor: Todos

## 27. Definition of Done

Este documento será considerado concluído quando:
- o Core puder executar sem conhecimento do SO;
- cada sistema possuir seus adapters;
- contratos de plataforma estiverem definidos;
- CI validar todas as plataformas;
- o GhostStack possuir comportamento consistente em Linux, Windows e macOS.
