# GhostStack — VPN/Tor/Proxy Operational Specification

Versão: 1.0.0
Status: Normativo
Tipo: Privacy Runtime Specification
Dependências: Runtime, Networking & Gateway Specification, Provider Integration Specification, Threat Model & Security Architecture, Firewall Specification, DNS Specification, Configuration Specification, Security Response Specification

## 1. Objetivo

Definir o comportamento operacional dos módulos: VPN, Tor, Proxy, DNS seguro, Firewall, Gateway de tráfego.

## 2. Filosofia

O GhostStack deve permitir diferentes níveis de privacidade sem forçar um único modelo. Princípios: controle explícito, comportamento previsível, nenhuma rota fora da política, falha segura.

## 3. Arquitetura de Rede

Applications -> Ghost Gateway -> Privacy Layer -> Firewall -> Network Interface -> Internet

## 4. Componentes

Camadas: Application Layer, Gateway Layer, Proxy Layer, Tor Layer, VPN Layer, Firewall Layer, Physical Network.

## 5. Modos Operacionais

Perfis oficiais: direct, vpn, tor, vpn_tor, proxy, vpn_proxy, ghost.

## 6. Perfil Direct

Objetivo: uso normal. Fluxo: Application -> Network -> Internet. Proteções: firewall, DNS seguro.

## 7. Perfil VPN

Objetivo: todo tráfego passa pelo túnel VPN. Fluxo: Application -> Gateway -> VPN Tunnel -> Internet. Requisitos: túnel ativo, rota validada, kill switch habilitado.

## 8. Perfil Tor

Objetivo: encaminhar tráfego pelo Tor. Fluxo: Application -> SOCKS Proxy -> Tor Network -> Internet. Requisitos: daemon Tor ativo, SOCKS disponível, circuito válido.

## 9. Perfil VPN + Tor

Modo Ghost. Fluxo padrão: Application -> Tor -> VPN -> Internet. Alternativa: Application -> VPN -> Tor -> Internet. A ordem deve ser definida explicitamente em privacy.chain.

## 10. Justificativa da Ordem

Tor sobre VPN: ISP não vê uso direto do Tor, entrada Tor fica escondida atrás da VPN. VPN sobre Tor: VPN não vê IP real, configuração mais complexa, depende do suporte do provedor.

## 11. Perfil Proxy

Fluxo: Application -> Proxy Gateway -> Internet. Suporta: HTTP, HTTPS, SOCKS5.

## 12. Perfil VPN + Proxy

Fluxo: Application -> Proxy -> VPN -> Internet. Uso: ambientes específicos, testes, controle de saída.

## 13. Perfil Ghost

Perfil recomendado: runtime: vpn_tor, firewall.kill_switch: true, dns.secure: true. Características: VPN, Tor, DNS protegido, bloqueio de fuga.

## 14. Firewall Kill Switch

Objetivo: bloquear tráfego caso a proteção caia. Estados: ACTIVE, BLOCKING, RECOVERING, DISABLED.

## 15. DNS Leak Prevention

O DNS nunca deve sair fora da rota configurada. Bloquear: DNS do ISP, resolução direta, fallback inseguro.

## 16. IPv6 Handling

Opções: disable, allow, route. Padrão: disable quando não houver suporte seguro.

## 17. Traffic Routing Rules

Permite regras por aplicação. Exemplo: routes.browser.mode: tor, routes.gaming.mode: direct.

## 18. Application Isolation

Aplicações podem possuir políticas próprias. Exemplo: apps.firefox.network: tor.

## 19. Proxy Chain

Suporte: Proxy A -> Proxy B -> VPN -> Internet. Limites: máximo configurável, timeout obrigatório.

## 20. Health Monitoring

Monitorar: VPN (handshake, latência, perda), Tor (bootstrap, circuito, disponibilidade), Proxy (resposta, tempo), DNS (resolução, segurança).

## 21. Reconexão Automática

Fluxo: Failure Detected -> Stop Traffic -> Reconnect -> Validate -> Resume Traffic.

## 22. Circuit Rotation

Tor: ghost tor new-circuit. Configuração: tor.rotate.enabled: true, interval: 10m.

## 23. Privacy Levels

Níveis: Standard (Firewall, Secure DNS), Enhanced (VPN, Secure DNS, Kill Switch), Maximum (VPN, Tor, Proxy, Kill Switch, Isolation).

## 24. Performance Trade-off

Maior privacidade pode gerar: maior latência, menor velocidade, maior consumo. O usuário deve escolher o equilíbrio.

## 25. Logs

Não registrar: sites acessados, conteúdo, dados pessoais. Registrar apenas: estado, erros, métricas técnicas.

## 26. Segurança Contra Falhas

Falhas críticas: VPN_DOWN, TOR_FAILED, DNS_LEAK, ROUTE_BYPASS. Resposta: BLOCK, ALERT, RECOVER.

## 27. Testes

Obrigatórios: teste DNS leak, teste VPN drop, teste Tor failure, teste proxy failure, teste recuperação, teste IPv6.

## 28. Definition of Done

Este documento será considerado concluído quando:
- todas as combinações de privacidade estiverem definidas;
- roteamento for determinístico;
- vazamentos forem bloqueados;
- falhas forem recuperáveis;
- usuários puderem escolher níveis de privacidade.
