# GhostStack — CLI Specification

Versão: 1.0.0
Status: Normativo
Tipo: User Interface & Administration Specification
Dependências: PRD, System Architecture, Cross Platform Architecture Specification, API Governance & Versioning Specification, Configuration Specification, Runtime Specification, Security Operations Specification, Observability & Monitoring Specification

## 1. Objetivo

Definir a interface oficial de linha de comando do GhostStack.

## 2. Filosofia

A CLI deve ser: simples, previsível, automatizável, segura, compatível entre sistemas.

## 3. Binário Oficial

Nome: ghost
Exemplos:
- Linux: ghost status
- Windows: ghost status
- macOS: ghost status

## 4. Arquitetura

A CLI não contém lógica de negócio. Fluxo: CLI -> API Client -> Ghost Daemon -> Core Runtime.

## 5. Modos de Saída

Suporte obrigatório:
- Human Mode: padrão
- JSON Mode: ghost status --json
- Quiet Mode: ghost command --quiet

## 6. Estrutura de Comandos

Formato: ghost <resource> <action>. Exemplo: ghost vpn connect.

## 7. Comandos Principais

7.1 Status: ghost status
7.2 Runtime: ghost runtime start/stop/restart/reload
7.3 Configuração: ghost config show/validate/reload/edit
7.4 Profiles: ghost profile list/use/create
7.5 VPN: ghost vpn status/connect/disconnect/reconnect
7.6 Tor: ghost tor status/start/stop/new-circuit
7.7 Proxy: ghost proxy status/test/reload
7.8 Firewall: ghost firewall status/enable/disable/rules
7.9 DNS: ghost dns status/test/flush
7.10 Plugins: ghost plugin list/install/remove/update
7.11 Storage: ghost storage verify/backup/restore
7.12 Security: ghost security check/scan/events
7.13 Diagnostics: ghost diagnose/export

## 8. Flags Globais

--version, --help, --json, --verbose, --config path.yaml

## 9. Autenticação

Comandos administrativos podem exigir: ghost auth login/status/logout.

## 10. Permissões

Operações críticas exigem privilégio: Linux sudo, Windows Run as Administrator, macOS sudo.

## 11. Confirmações

Operações destrutivas exigem confirmação. Modo automático: --yes.

## 12. Interactive Mode

ghost setup com fluxo: Detect System -> Choose Profile -> Configure Network -> Validate -> Start.

## 13. Autocomplete

Suporte: Bash, Zsh, Fish, PowerShell. Exemplo: ghost completion bash.

## 14. Erros

Formato humano: Error: VPN provider unavailable, Code: VPN-002.
Formato JSON: { "error":"VPN-002", "message":"Provider unavailable" }.

## 15. Exit Codes

0: Sucesso, 1: Erro genérico, 2: Configuração inválida, 3: Permissão negada, 4: Serviço indisponível, 5: Falha de segurança.

## 16. Scripts

A CLI deve permitir automação: ghost status --json | jq.

## 17. Compatibilidade

A mesma CLI deve funcionar em: Linux, Windows, macOS.

## 18. Atualização

ghost update com fluxo: Check Version -> Verify Signature -> Download -> Install -> Restart.

## 19. Logs da CLI

Modo debug: ghost --debug status. Não deve exibir secrets, tokens, credenciais.

## 20. Segurança

A CLI nunca deve: salvar senha em histórico, exibir secrets, ignorar políticas do Core.

## 21. Testes

Obrigatórios: todos comandos, flags, JSON output, permissões, erros, autocomplete.

## 22. Definition of Done

Este documento será considerado concluído quando:
- todos recursos possuírem comandos;
- automação for possível;
- saída for previsível;
- funcionar nos três sistemas operacionais;
- usuários puderem administrar o GhostStack sem interface gráfica.
