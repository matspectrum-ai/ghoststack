# GhostStack — Deployment & Distribution Specification (Cross Platform)

Versão: 1.0.0
Status: Normativo
Tipo: Deployment Specification
Dependências: PRD, System Architecture, Cross Platform Architecture Specification, Runtime Specification, Configuration Specification, Plugin SDK Specification, Testing & Verification Specification, Development Standards Specification

## 1. Objetivo

Definir o processo oficial de distribuição, instalação, atualização e manutenção do GhostStack em múltiplos sistemas operacionais.

## 2. Princípios

O sistema de distribuição deve ser: seguro, reproduzível, verificável, automatizável, compatível com atualização segura.

## 3. Artefatos Oficiais

Cada release gera artefatos independentes:
- GhostStack-v1.0.0-linux-amd64
- GhostStack-v1.0.0-windows-amd64.exe
- GhostStack-v1.0.0-macos-amd64.pkg
- GhostStack-v1.0.0-macos-arm64.pkg

## 4. Canais de Release

- Stable: usuários finais, produção
- Beta: testes, novos recursos
- Nightly: desenvolvimento, validação automática

## 5. Estrutura de Instalação

Core + CLI + Daemon + Providers + Plugins

## 6. Linux Deployment

Formatos suportados: .deb (principal), .tar.gz (secundário)

## 7. Linux Paths

Binários: /usr/bin/ghost
Core: /usr/lib/ghoststack
Configuração: /etc/ghoststack
Dados: /var/lib/ghoststack
Logs: /var/log/ghoststack

## 8. Linux Service

Gerenciador: systemd
Serviço: ghoststack.service
Estados: active, inactive, failed, restarting

## 9. Windows Deployment

Formatos suportados: .exe installer (principal), .msi (alternativo)

## 10. Windows Paths

Instalação: C:\Program Files\GhostStack
Dados: C:\ProgramData\GhostStack
Logs: C:\ProgramData\GhostStack\Logs
Configuração: C:\ProgramData\GhostStack\Config

## 11. Windows Service

O GhostStack executa como GhostStack Service, gerenciado pelo Windows Service Control Manager.

## 12. macOS Deployment

Formatos suportados: .pkg (principal), .dmg (alternativo)

## 13. macOS Paths

Aplicação: /Applications/GhostStack
Dados: /Library/Application Support/GhostStack
Logs: /Library/Logs/GhostStack
Configuração: /Library/Application Support/GhostStack/config

## 14. macOS Service

Gerenciador: launchd
Serviço: com.ghoststack.daemon

## 15. Primeiro Setup

Fluxo: Install -> Initialize -> Validate Environment -> Create Identity -> Generate Config -> Start Service

## 16. Comando de Inicialização

Universal: ghost init. Executa: criação de diretórios, validação, geração de configuração, verificação de dependências.

## 17. Verificação de Ambiente

Comando: ghost doctor. Verifica sistema, rede e dependências.

## 18. Atualizações

Comando: ghost update. Fluxo: Check Version -> Download -> Verify Signature -> Backup -> Migration -> Install -> Restart -> Health Check.

## 19. Assinatura de Releases

Todo release oficial deve possuir artifact, checksum e signature.

## 20. Rollback

Comando: ghost rollback. Fluxo: Current Version -> Backup -> Previous Version -> Restore -> Restart.

## 21. Migração

Atualizações podem alterar configuração, banco, plugins, schemas. Fluxo: Backup -> Migration -> Validation -> Commit. Falha: Rollback automático.

## 22. Remoção

Comando: ghost uninstall. Modos: --keep-data (manter dados) ou --purge (remover tudo).

## 23. Portable Mode

Disponível em todos os sistemas. Executável: ghost portable. Características: sem instalação, sem serviço, configuração local.

## 24. Docker/Container

Não é requisito. O GhostStack deve funcionar: nativamente, em VM, em container.

## 25. Build Pipeline

Toda release passa por: Source -> Build -> Unit Tests -> Integration Tests -> Security Scan -> Package -> Sign -> Publish.

## 26. Compatibilidade de Arquitetura

Suporte inicial: AMD64. Planejado: ARM64.

## 27. Distribuição Oficial

- Linux: repositório APT
- Windows: site oficial, package manager futuro
- macOS: package oficial, Homebrew futuro

## 28. Telemetria

Por padrão: security.telemetry: false. Nenhum dado é coletado sem consentimento explícito.

## 29. Recuperação

Modo seguro: ghost recovery. Permite: iniciar sem plugins, restaurar configuração, desativar providers defeituosos.

## 30. Definition of Done

Este documento será considerado concluído quando:
- existir processo de instalação para Linux, Windows e macOS;
- releases forem reproduzíveis;
- atualizações forem seguras;
- rollback funcionar;
- o sistema puder ser operado sem conhecimento interno da arquitetura.
