# GhostStack — Deployment & Distribution Specification

Versão: 1.0.0
Status: Normativo
Dependências: PRD, System Architecture, Threat Model & Security Architecture, Runtime Specification, ADR-001 Runtime Abstraction, Configuration Specification, OpenAPI Specification, Plugin SDK Specification, Testing & Verification Specification, Implementation Roadmap, Development Standards Specification

## 1. Objetivo

Definir o modelo oficial de distribuição e implantação do GhostStack.

## 2. Sistemas Suportados

Versão inicial:
- Debian: Oficial
- Ubuntu: Oficial
- antiX Linux: Compatível

Arquitetura: x86_64
Suporte futuro: ARM64

## 3. Modelo de Instalação

Três métodos oficiais:

Método 1 — Pacote nativo (.deb): Recomendado para Debian, Ubuntu, antiX.
Método 2 — Binário standalone (ghoststack-linux-amd64.tar.gz): Para servidores e ambientes personalizados.
Método 3 — Build Source: Para desenvolvedores (git clone, make build).

## 4. Estrutura do Sistema

Após instalação:
/
├── usr/
│   ├── bin/
│   │   └── ghost
│   └── lib/
│       └── ghoststack/
├── etc/
│   └── ghoststack/
│       ├── ghost.yaml
│       └── schemas/
├── var/
│   ├── lib/
│   │   └── ghoststack/
│   └── log/
│       └── ghoststack/
└── run/
    └── ghoststack/

## 5. Componentes Instalados

Pacote principal: ghoststack-core (daemon, CLI, schemas, serviços, arquivos padrão).
Pacotes opcionais: ghoststack-dashboard, ghoststack-plugins, ghoststack-sdk.

## 6. Serviço Systemd

Serviço: ghoststack.service em /etc/systemd/system/
Exemplo:
[Unit]
Description=GhostStack Privacy Runtime
After=network.target

[Service]
ExecStart=/usr/bin/ghost daemon
Restart=always

[Install]
WantedBy=multi-user.target

## 7. Usuário do Sistema

Instalação cria usuário ghoststack: sem login, sem shell, permissões mínimas.

## 8. Permissões

Config (/etc/ghoststack): root:ghoststack 750
Logs (/var/log/ghoststack): ghoststack:ghoststack 750

## 9. Primeiro Start

Fluxo: Install -> Generate Default Config -> Validate Config -> Initialize Runtime -> Start Service -> Health Check.

## 10. Inicialização

Comando: ghost init. Responsável por criar configuração, gerar identidade local, validar dependências e preparar diretórios.

## 11. Configuração Inicial

Arquivo criado: /etc/ghoststack/ghost.yaml. Perfil padrão: profile: direct.

## 12. Atualização

Comando: ghost update. Fluxo: Check Version -> Download Package -> Verify Signature -> Backup -> Install -> Migration -> Restart.

## 13. Verificação de Integridade

Todo pacote deve possuir checksum, assinatura e versão. Exemplo: ghoststack.deb, ghoststack.deb.sha256, ghoststack.deb.sig.

## 14. Rollback

Suportado. Comando: ghost rollback. Fluxo: Current Version -> Backup -> Previous Version -> Restore -> Restart.

## 15. Migrações

Antes de atualizar: Config Backup -> Database Migration -> Runtime Migration -> Validation. Falha durante migração: Abort -> Rollback.

## 16. Remoção

Comando: ghost uninstall. Opções: --keep-data (manter dados) ou --purge (remover tudo).

## 17. Diagnóstico

Comando: ghost doctor. Verifica: configuração, permissões, serviços, rede, DNS, firewall, dependências.

## 18. Backup

Comando: ghost backup. Inclui: configuração, profiles, plugins, estado necessário. Nunca inclui senhas ou secrets em texto puro.

## 19. Recuperação

Comando: ghost restore backup.tar. Fluxo: Validate Backup -> Restore Files -> Validate Config -> Restart.

## 20. Compatibilidade Kernel

Requisitos mínimos: Linux Kernel >= 5.10. Recursos utilizados: namespaces, nftables, WireGuard, systemd.

## 21. Dependências Externas

Dependências opcionais: nftables, wireguard-tools, tor, unbound, sing-box. O GhostStack deve detectar ausência e informar.

## 22. Modo Portable

Suportado. Execução: ./ghost daemon. Limitações: sem systemd, sem integração completa, sem atualização automática.

## 23. Containerização

Não é requisito obrigatório. O GhostStack deve funcionar: nativamente no Linux, em VM, em container, em ambiente imutável.

## 24. Observabilidade

Instalação deve disponibilizar: ghost status, ghost logs, ghost metrics, ghost health.

## 25. Pacotes de Release

Release oficial:
- ghoststack-v1.0.0-amd64.deb
- ghoststack-v1.0.0-linux-amd64.tar.gz

## 26. Canais de Distribuição

Canais: stable, beta, nightly.

## 27. Política de Segurança

Atualizações críticas podem exigir: assinatura obrigatória, atualização forçada, bloqueio de versões vulneráveis.

## 28. Definition of Done

Este documento será considerado concluído quando:
- o GhostStack puder ser instalado em Debian/Ubuntu/antiX;
- o serviço iniciar automaticamente;
- upgrades forem seguros;
- rollback funcionar;
- remoção puder ser controlada;
- releases forem reproduzíveis.
