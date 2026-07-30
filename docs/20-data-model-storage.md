# GhostStack — Data Model & Storage Specification

Versão: 1.0.0
Status: Normativo
Tipo: Data Architecture Specification
Dependências: PRD, System Architecture, Domain Model, Configuration Specification, API Governance & Versioning Specification, Security Operations Specification, Observability & Monitoring Specification, Deployment & Distribution Specification

## 1. Objetivo

Definir o modelo oficial de dados e armazenamento do GhostStack.

## 2. Princípios

O armazenamento deve ser: local por padrão, seguro, versionável, migrável, resiliente, independente do sistema operacional.

## 3. Estratégia de Persistência

O GhostStack utiliza arquitetura híbrida: Configuration Files + Local Database + Runtime State + Cache.

## 4. Armazenamento Principal

Banco inicial: SQLite. Motivos: zero dependência externa, portátil, funciona em Linux, Windows e macOS, adequado para agente local.

## 5. Localização do Banco

Linux: /var/lib/ghoststack/ghost.db
Windows: C:\ProgramData\GhostStack\Data\ghost.db
macOS: /Library/Application Support/GhostStack/Data/ghost.db

## 6. Camada de Abstração

O Core não acessa SQLite diretamente. Interface StorageProvider com Open, Close, Query, Insert, Update, Delete, Migrate. Implementações futuras: SQLiteProvider, PostgreSQLProvider, RemoteStorageProvider.

## 7. Estrutura do Banco

Schema: ghoststack com tabelas system, runtime, network, vpn, tor, plugins, events, audit, migrations.

## 8. Tabela System

Responsável pela identidade da instalação. Campos: id, installation_id, version, created_at, updated_at.

## 9. Tabela Runtime State

Armazena estado atual. Campos: id, status, mode, started_at, last_health_check, updated_at. Estados: starting, running, stopped, failed, recovering.

## 10. Tabela Network State

Armazena informações operacionais. Campos: interface, gateway, dns_provider, ipv6_mode, updated_at. Não armazenar: histórico de navegação, destinos acessados, conteúdo de tráfego.

## 11. Tabela VPN

Tabela vpn_profiles. Campos: id, name, provider, status, created_at, updated_at. Segredos não ficam nesta tabela.

## 12. Tabela Tor

Tabela tor_state. Campos: status, bootstrap_percent, last_circuit_change, updated_at. Não armazenar: histórico de circuitos, destinos, conexões.

## 13. Tabela Plugins

Tabela plugins. Campos: id, name, version, enabled, permissions, installed_at, updated_at.

## 14. Tabela Events

Armazena eventos internos. Tabela events. Campos: id, type, source, severity, timestamp, payload.

## 15. Tabela Audit

Eventos de segurança. Tabela audit_log. Campos: id, action, actor, timestamp, result.

## 16. Tabela Migrations

Controle de schema. Tabela schema_migrations. Campos: version, applied_at.

## 17. Migrações

Fluxo: Current Schema -> Migration Check -> Apply Migration -> Validate -> Commit. Falha: Rollback automático.

## 18. Versionamento de Schema

Formato: major.minor. Exemplo: 1.0. Mudanças incompatíveis: 2.0.

## 19. Backup

Comando: ghost backup. Inclui: banco, configuração, plugins, metadados. Não inclui: secrets, tokens, chaves privadas.

## 20. Backup Seguro

Fluxo: Backup -> Checksum -> Compressão -> Armazenamento. Formato: ghost-backup-YYYY-MM-DD.tar.

## 21. Restore

Comando: ghost restore backup.tar. Fluxo: Validate Backup -> Stop Runtime -> Restore Data -> Run Migration -> Start Runtime.

## 22. Cache

Cache permitido: estado temporário, métricas, resultados de health check. Cache proibido: credenciais, histórico de rede, dados privados.

## 23. Retenção

Eventos: retention.events: 30d, retention.logs: 30d. Configuração customizável.

## 24. Limpeza Automática

Processo: Scheduler -> Check Expiration -> Delete Expired Data -> Record Cleanup Event.

## 25. Integridade

Banco deve possuir: checksum, validação de schema, proteção contra corrupção. Comando: ghost storage verify.

## 26. Criptografia

Dados sensíveis devem usar armazenamento seguro do SO ou criptografia quando necessário. O banco não deve armazenar: senhas, tokens, chaves privadas.

## 27. Concorrência

O Storage Layer deve suportar: múltiplos módulos lendo, escrita controlada, transações. SQLite em modo WAL.

## 28. Performance

Requisitos: inicialização < 2 segundos, consultas comuns indexadas, limpeza assíncrona.

## 29. Testes

Obrigatórios: criação do banco, migrações, rollback, backup, restore, integridade, corrupção simulada.

## 30. Definition of Done

Este documento será considerado concluído quando:
- o modelo de dados estiver definido;
- o armazenamento for multiplataforma;
- migrações forem seguras;
- backups forem reproduzíveis;
- nenhuma informação sensível for persistida incorretamente.
