# GhostStack — Configuration Specification (ghost.yaml)

Versão: 1.0
Status: Normativo
Dependências: PRD, System Architecture, Threat Model, Runtime & Gateway Specification, ADR-001 Runtime Abstraction, Domain Model

## 1. Objetivo

Este documento define a especificação oficial do arquivo de configuração do GhostStack. Toda configuração persistente do sistema deverá ser representável através de ghost.yaml.

## 2. Princípios

O arquivo deve ser: humano-legível, determinístico, amigável para IA, facilmente validável, independente da implementação, compatível com Git. Nunca utilizar formatos binários.

## 3. Configuração x Estado

O arquivo ghost.yaml representa configuração permanente. Ele nunca armazenará: estado atual, uptime, PID, conexões, métricas, modo atualmente ativo. Essas informações pertencem ao Runtime State.

## 4. Localização

Sistema: /etc/ghoststack/ghost.yaml
Usuário: ~/.config/ghoststack/ghost.yaml
Projeto: ./ghost.yaml

Prioridade: Projeto -> Usuário -> Sistema

## 5. Estrutura

version: 1
compatibility:
metadata:
runtime:
gateway:
firewall:
dns:
vpn:
tor:
proxy:
plugins:
dashboard:
logging:
storage:
security:
profiles:
include:

## 6. Compatibilidade

version: 1
compatibility:
  schema: "1.0"
  daemon: ">=1.0"
  api: "v1"
  gateway: ">=1.0"

## 7. Metadata

metadata:
  name:
  description:
  author:
  created:
  updated:

Campos apenas informativos.

## 8. Runtime

runtime:
  provider: native
  default_mode: direct
  auto_start: true
  auto_update: false
  reload_strategy: hot
  health_interval: 5s

Providers oficiais: native, docker, podman

## 9. Gateway

gateway:
  backend: sing-box
  tun:
    enabled: true
    interface: ghost0
  inbound:
    socks: 1080
    http: 8080

## 10. Firewall

firewall:
  enabled: true
  backend: nftables
  kill_switch: true
  fail_closed: true
  ipv6:
    mode: disable

Valores permitidos: disable, allow, route

## 11. DNS

dns:
  provider: unbound
  secure: true
  cache: true
  dnssec: true
  doh: false
  dot: true

## 12. VPN

vpn:
  enabled: false
  provider: wireguard
  profile: default
  reconnect: true

## 13. Tor

tor:
  enabled: false
  socks_port: 9050
  control_port: 9051
  auto_rotate: false
  rotate_interval: 10m

## 14. Proxy

proxy:
  enabled: false
  strategy: single
  timeout: 10s

Estratégias: single, priority, random, round_robin, latency

## 15. Logging

logging:
  level: info
  format: json
  stdout: true
  file:
    enabled: true
    path: logs/

## 16. Storage

storage:
  database: sqlite
  backup: true
  retention: 30d

## 17. Dashboard

dashboard:
  enabled: true
  port: 9090
  auth: local

## 18. Plugins

plugins:
  enabled: true
  repositories:
    - official
    - local
  allow_unsigned: false

## 19. Security

security:
  telemetry: false
  metrics: true
  verify_checksums: true
  verify_signatures: true

## 20. Secrets

secrets:
  backend: env
  providers:
    - env
    - file
    - keyring

Nunca armazenar senhas em texto puro.

## 21. Includes

Permite dividir configurações.

include:
  - vpn.yaml
  - dns.yaml
  - firewall.yaml

## 22. Profiles

profiles:
  base:
    dns:
      secure: true
  office:
    extends:
      - base
    runtime:
      default_mode: vpn

Herança múltipla é suportada.

## 23. Merge

Ordem: System -> User -> Project -> Include -> Profile -> Environment -> CLI

Política:
- Objetos: Deep Merge
- Arrays: Substituição
- Escalares: Último vence

## 24. Variáveis de Ambiente

Suportadas.

vpn:
  username: ${VPN_USERNAME}
  password: ${VPN_PASSWORD}

## 25. Secrets por Arquivo

vpn:
  password_file: /run/secrets/vpn_password

## 26. Valores Padrão

Na ausência de configuração:
runtime:
  provider: native
  default_mode: direct
firewall:
  enabled: true
logging:
  level: info

## 27. Validação

Toda chave deverá possuir: tipo, obrigatoriedade, enum, valor padrão, documentação. Nenhuma configuração inválida poderá iniciar o Runtime.

## 28. JSON Schema

Todo release deverá publicar: schemas/ghost.schema.json. Esse schema é a referência normativa.

## 29. Migração

Comando: ghost migrate
Fluxo: v1 -> v2 -> v3
Backup obrigatório antes da migração.

## 30. Runtime State

O estado do sistema será persistido separadamente. Local padrão: /var/lib/ghoststack/state.json. Esse arquivo nunca substitui ghost.yaml.

## 31. API de Configuração

Endpoints oficiais:
- GET /config
- PUT /config
- POST /config/validate
- POST /config/reload
- POST /config/migrate

## 32. Eventos

Eventos oficiais: CONFIG_CHANGED, CONFIG_VALIDATED, CONFIG_RELOADED, PROFILE_CHANGED, CONFIG_MIGRATED

## 33. Requisitos Normativos

O sistema deverá: validar toda configuração antes da execução, gerar erros determinísticos, suportar YAML 1.2, gerar JSON Schema automaticamente, permitir reload parcial, impedir configuração inconsistente.

## Definition of Done

Este documento será considerado concluído quando:
- Toda configuração estiver formalmente especificada.
- O JSON Schema puder ser gerado automaticamente.
- Todos os campos possuírem tipo e validação.
- O parser puder ser implementado exclusivamente com base nesta especificação.
- O Runtime puder carregar uma configuração sem depender de documentação adicional.
