# GhostStack — System Architecture

Versão: 0.1.0
Status: Draft
Dependência: PRD v0.1

## 1. Objetivo

Definir a arquitetura técnica completa do GhostStack.

Este documento será a fonte única de verdade para Engenharia, Codex, Claude Code, OpenCode, CI/CD, QA e documentação.

Nenhuma implementação poderá divergir desta arquitetura sem uma ADR aprovada.

## 2. Princípios Arquiteturais

PA-001 — API First
Toda funcionalidade deverá ser acessível por API. A CLI é apenas um cliente da API.

PA-002 — Core Stateless
O Core não deverá manter estado em memória além do necessário. Todo estado persistente ficará em armazenamento próprio.

PA-003 — Modularidade
Cada módulo pode ser substituído. Exemplo: Tor -> I2P -> Nym sem alterar outros módulos.

PA-004 — Least Privilege
Cada container recebe apenas as permissões mínimas necessárias. Nunca executar como root sem necessidade técnica comprovada.

PA-005 — Security by Default
Após instalar o GhostStack: todo tráfego deverá seguir as regras do Gateway. Nenhum bypass será permitido por padrão.

## 3. Arquitetura Geral

                             User
                               │
                   ghost CLI / Dashboard
                               │
──────────────────────────────────────────────────────────
                    Ghost Daemon (Core)
──────────────────────────────────────────────────────────
      Config      Runtime      Scheduler     Health
          │            │             │            │
──────────────────────────────────────────────────────────
                Gateway (sing-box)
──────────────────────────────────────────────────────────
     VPN       Tor       Proxy       DNS       Firewall
──────────────────────────────────────────────────────────
                 Docker Runtime
──────────────────────────────────────────────────────────
                    Linux Kernel

## 4. Arquitetura C4

Context Diagram:
                 User
                   │
             Ghost CLI
                   │
             Ghost Daemon
                   │
──────────────────────────────────
      VPN Provider
      Tor Network
      Proxy Provider
      DNS Resolver
──────────────────────────────────
                Internet

## 5. Containers

Core
- ghost-core: API, Scheduler, Runtime, Estado, Eventos, Gateway

Gateway
- ghost-gateway: sing-box, TUN, Routing, Policy

Tor
- ghost-tor: Tor daemon, SOCKS, Control Port, DNS

DNS
- ghost-dns: Unbound, DNSSEC, DoH, DoT

Firewall
- ghost-firewall: nftables, Kill Switch, Leak Protection

VPN
- ghost-vpn: WireGuard, OpenVPN

Dashboard
- ghost-dashboard: Tauri + React

## 6. Componentes

Ghost Core: iniciar stack, parar stack, health, logs, eventos, configuração, plugins. Nunca roteia pacotes.

Gateway: responsável exclusivamente pelo roteamento. Toda decisão de rede passa por ele. Nunca faz persistência.

DNS: resolve nomes. Nunca inicia VPN. Nunca gerencia Tor.

Firewall: único responsável pelas regras do sistema. Nenhum outro módulo pode alterar nftables.

VPN: responsável apenas por conectar, desconectar, trocar servidor. Nunca modifica DNS.

Tor: responsável apenas pela rede Tor.

## 7. Comunicação

CLI -> REST -> Core -> gRPC -> Gateway
Eventos: Gateway -> Core -> Dashboard

## 8. Event Bus

Todos os módulos publicam eventos. Exemplos: VPN_CONNECTED, TOR_READY, DNS_LEAK, HEALTH_FAIL, MODE_CHANGED, CONTAINER_STARTED, CONTAINER_STOPPED.

Formato:
{
  "id": "...",
  "timestamp": "...",
  "module": "...",
  "event": "...",
  "payload": {}
}

## 9. Redes Docker

ghost-control: comunicação interna, nunca exposta.
ghost-runtime: comunicação entre módulos.
ghost-public: somente Gateway. Apenas ele conversa com Internet.

## 10. Volumes

config/, logs/, runtime/, cache/, state/, plugins/, certificates/

## 11. Fluxo de Inicialização

ghost up -> Docker -> Core -> Config -> Gateway -> Firewall -> DNS -> VPN -> Tor -> Health -> READY

Se qualquer etapa falhar: ROLLBACK

## 12. Fluxo de Desligamento

ghost down -> Stop Dashboard -> Stop VPN -> Stop Tor -> Flush Gateway -> Persist State -> Remove Runtime -> DONE

## 13. Trust Boundaries

Usuário -> CLI -> Ghost Daemon -> Gateway -> Internet

Toda comunicação entre níveis deve ser autenticada e validada. Nunca confiar em dados externos.

## 14. Estado

Persistido em SQLite. Estruturas: settings, profiles, proxies, vpn, plugins, history, events.

## 15. Configuração

Arquivo único: ghost.yaml

Exemplo:
version: 1
mode: direct
gateway:
  provider: sing-box
dns:
  provider: unbound
vpn:
  enabled: false
tor:
  enabled: false
firewall:
  enabled: true

## 16. Runtime

Estados possíveis (máquina de estados obrigatória):
CREATED -> INITIALIZING -> STARTING -> READY -> RUNNING -> DEGRADED -> RECOVERING -> FAILED -> STOPPING -> STOPPED

Não poderão existir estados adicionais.

## 17. Modo de Operação

DIRECT: Aplicação -> Gateway -> Internet
VPN: Aplicação -> Gateway -> WireGuard/OpenVPN -> Internet
TOR: Aplicação -> Gateway -> Tor -> Internet
VPN_TOR: Aplicação -> Gateway -> VPN -> Tor -> Internet
PROXY: Aplicação -> Gateway -> SOCKS5 -> Internet
VPN_PROXY: Aplicação -> Gateway -> VPN -> Proxy -> Internet

## 18. Falhas

Cada módulo implementa: Retry -> Backoff -> Recovery -> Rollback. Nunca falhar silenciosamente.

## 19. Observabilidade

Todos os módulos devem expor: Health, Metrics, Logs, Events, Version. Endpoints internos: /health, /metrics, /version.

## 20. Arquitetura do Monorepo

ghoststack/
apps/    -> daemon, cli, dashboard
services/ -> gateway, firewall, dns, tor, vpn
packages/ -> config, events, logger, contracts, sdk, testing
docker/   -> compose, images
docs/
spec/
tests/
scripts/
.github/

## 21. Decisões Arquiteturais (pré-ADRs)

ID | Decisão | Status
AD-001 | O gateway único será o sing-box. | Accepted
AD-002 | O firewall utilizará nftables como backend principal. | Accepted
AD-003 | O DNS padrão será Unbound. | Accepted
AD-004 | O estado persistente será armazenado em SQLite. | Accepted
AD-005 | Todo o tráfego de saída deverá passar exclusivamente pelo Gateway. | Accepted
AD-006 | A CLI nunca manipulará Docker diretamente; todas as operações serão mediadas pelo Daemon. | Accepted
AD-007 | O Dashboard consumirá a mesma API pública utilizada pela CLI. | Accepted
AD-008 | Cada módulo será executado como um serviço/container independente, facilitando substituição e testes. | Accepted
ADR-001 | Runtime Abstraction — Core depende de IRuntimeProvider, não de Docker diretamente. | Accepted

## 22. Riscos Arquiteturais

Risco | Mitigação
Vazamento de DNS | Gateway + Unbound + testes automáticos
Bypass do firewall | nftables gerenciado exclusivamente pelo módulo Firewall
Falha na VPN | Kill Switch + rollback para estado seguro
Configuração inconsistente | Validação de ghost.yaml antes da inicialização
Falha de um container | Supervisão pelo Core com reinicialização controlada

## Definition of Done

A arquitetura será considerada aprovada quando:
- Todos os componentes tiverem responsabilidades claramente definidas.
- Todas as dependências entre módulos estiverem documentadas.
- Os limites de confiança estiverem especificados.
- A máquina de estados do runtime estiver definida.
- Os fluxos de inicialização, desligamento e recuperação estiverem completos.
- Não houver dependências circulares entre módulos.
