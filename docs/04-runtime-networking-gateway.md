# GhostStack — Runtime, Networking & Gateway Specification

Versão: 0.1.0
Status: Normativo
Dependências: PRD, Architecture, Threat Model

## 1. Objetivo

Especificar completamente o Runtime do GhostStack. Nenhuma implementação poderá modificar estes comportamentos.

## 2. Runtime

O Runtime é responsável por: iniciar serviços, parar serviços, reiniciar serviços, monitorar serviços, recuperar serviços. O Runtime não toma decisões de roteamento. Essas decisões pertencem ao Gateway.

## 3. Componentes

ghost-cli -> REST API -> ghostd -> Runtime Manager, Configuration Manager, Health Manager, Plugin Manager, Container Manager, Event Bus -> Gateway Controller -> Docker Engine

## 4. Estados

Todo módulo deverá implementar a seguinte máquina de estados:
CREATED -> INITIALIZING -> STARTING -> READY -> RUNNING -> DEGRADED -> RECOVERING -> FAILED -> STOPPING -> STOPPED

Não poderão existir estados adicionais.

## 5. Estado Global

O GhostDaemon mantém um estado global:
runtime:
  state: READY
  mode: vpn
  uptime: ...
  components:
    gateway: READY
    dns: READY
    vpn: READY
    tor: STOPPED
    firewall: READY

## 6. Modo de Operação

DIRECT: Aplicação -> Gateway -> Internet
VPN: Aplicação -> Gateway -> WireGuard/OpenVPN -> Internet
TOR: Aplicação -> Gateway -> Tor -> Internet
VPN_TOR: Aplicação -> Gateway -> VPN -> Tor -> Internet
PROXY: Aplicação -> Gateway -> SOCKS5 -> Internet
VPN_PROXY: Aplicação -> Gateway -> VPN -> Proxy -> Internet

## 7. Modos futuros

Reservados: I2P, NYM, TAILSCALE, SSH, VMESS, VLESS, SHADOWSOCKS. Nenhuma implementação deverá depender deles.

## 8. Gateway

O Gateway será implementado usando sing-box. Responsável por: TUN, Routing, Outbounds, DNS, Rules. Nunca: Logs, Configuração, Banco.

## 9. Interface TUN

O Gateway criará ghost0. Interface virtual. Todo tráfego será encaminhado para ela.

## 10. Outbounds

Cada backend gera um outbound: DIRECT, VPN, TOR, PROXY.

## 11. Seleção de Outbound

O Gateway possui apenas um outbound ativo. Trocar o outbound: ghost mode tor -> Gateway atualiza configuração -> Reload -> Novo outbound. Sem reiniciar containers não relacionados.

## 12. Inbounds

Sempre existirão: TUN, SOCKS, HTTP.

## 13. DNS

Sempre interceptado. Fluxo: Aplicação -> Gateway -> Unbound -> Resolver. Nunca: Aplicação -> ISP.

## 14. Routing Rules

Prioridade: BLOCK -> LOCAL -> DNS -> VPN -> TOR -> PROXY -> DIRECT. Primeira regra válida vence.

## 15. Redes Docker

ghost-control: comunicação interna, sem Internet.
ghost-runtime: comunicação entre módulos.
ghost-public: apenas Gateway. Nenhum outro container poderá entrar.

## 16. Firewall

Toda saída deverá obedecer: Container -> Gateway -> Firewall -> Internet. Nunca: Container -> Internet.

## 17. Container Manager

Responsável por: start, stop, restart, inspect, logs. Nunca: Routing, Configuração.

Implementação via IRuntimeProvider (ADR-001):
- NativeRuntimeProvider (padrão)
- DockerRuntimeProvider (opcional)
- PodmanRuntimeProvider (futuro)

## 18. Runtime Recovery

Falha: FAILED -> Restart -> Health -> READY. Caso falhe novamente: FAILED -> Rollback.

## 19. Health

Cada módulo expõe: Status, Version, Latency, Errors, StartedAt.

## 20. Heartbeat

Intervalo: 5 segundos (configurável).

## 21. Scheduler

Executa: health, rotação Tor, limpeza cache, atualização estado. Nunca: Routing.

## 22. Eventos

Formato:
{
  "event": "VPN_CONNECTED",
  "module": "vpn",
  "timestamp": "..."
}

## 23. Startup

Load Config -> Validate -> Create Networks -> Start Core -> Start Gateway -> Apply Firewall -> Start DNS -> Start Backend -> Health -> READY

## 24. Shutdown

Stop Scheduler -> Persist State -> Stop Backend -> Stop DNS -> Flush Gateway -> Stop Containers -> STOPPED

## 25. Reload

Permitido: ghost reload. Não deverá reiniciar a stack inteira. Somente componentes afetados.

## 26. Rollback

Toda alteração: Backup -> Apply -> Validate -> Commit. Falhou? Rollback. Obrigatório.

## 27. Configuração

Estrutura:
version: 1
runtime:
  mode: vpn
gateway:
  backend: sing-box
dns:
  provider: unbound
firewall:
  backend: nftables
vpn:
  provider: wireguard
tor:
  enabled: false

## 28. Performance

Startup: <= 5 segundos (sem considerar tempo de conexão externa).
Reload: <= 2 segundos.
Troca de modo: <= 5 segundos.
Health: < 100 ms.

## 29. Limites

Até: 100 perfis de configuração, 1 Gateway ativo, N backends instalados, 1 modo ativo por vez.

## 30. Tratamento de Erros

Todo erro deve possuir: código, mensagem, causa, ação recomendada.

Exemplo:
{
  "code": "GW-001",
  "message": "Gateway failed to initialize",
  "cause": "TUN interface creation failed",
  "recommendation": "Verify CAP_NET_ADMIN capability."
}

## 31. Contratos Internos

Todos os módulos devem implementar: Initialize(), Start(), Stop(), Reload(), Health(), Status(), Version(). Sem exceções.

## 32. Interfaces Lógicas

IGateway, IVPN, ITor, IDNS, IFirewall, IPlugin, IHealth.

## 33. Observabilidade

Todos os módulos deverão emitir: logs estruturados (JSON), eventos, métricas, health checks, informações de versão. Nenhum módulo poderá escrever logs em formato arbitrário.

## 34. Requisitos de Implementação

O Runtime deve ser determinístico. As operações devem ser idempotentes sempre que possível. Nenhum módulo pode acessar diretamente outro módulo sem passar pelos contratos definidos.

## Definition of Done

Este documento será considerado completo quando:
- O ciclo de vida de todos os componentes estiver definido.
- A topologia de rede estiver especificada.
- O comportamento do Gateway estiver completamente documentado.
- Todos os estados do Runtime forem implementáveis.
- Os contratos internos puderem ser derivados diretamente desta especificação.
