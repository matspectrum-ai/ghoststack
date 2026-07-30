# GhostStack — Dashboard UX & Management Interface Specification

Versão: 1.0.0
Status: Normativo
Tipo: User Interface Architecture Specification
Dependências: PRD, System Architecture, API Governance & Versioning Specification, CLI Specification, Configuration Specification, Observability & Monitoring Specification, Security Operations Specification, Plugin System & Extension Architecture Specification

## 1. Objetivo

Definir a arquitetura e experiência da interface gráfica oficial do GhostStack.

## 2. Filosofia

O Dashboard deve ser: simples para usuários comuns, poderoso para usuários avançados, transparente, seguro, consistente com CLI.

## 3. Arquitetura

O Dashboard é uma aplicação separada. Browser/Desktop UI -> Frontend Application -> API Gateway -> Ghost Core.

## 4. Tecnologia

Implementação inicial: Frontend React + TypeScript, Build Vite, Comunicação REST API + WebSocket Events.

## 5. Modos de Execução

Suporte: Local Web Dashboard (padrão, http://localhost:9090) e Desktop Wrapper (futuro: GhostStack Desktop baseado em Tauri/Electron).

## 6. Segurança de Acesso

Padrão: dashboard.enabled: true, bind: localhost, auth: local. Nunca: expor publicamente, abrir portas externas automaticamente.

## 7. Autenticação

Métodos: Inicial Local Authentication. Futuro: passkeys, certificados, MFA.

## 8. Estrutura Visual

Layout: Sidebar + Dashboard Content + Status Bar.

## 9. Navegação Principal

Menu: Dashboard, Runtime, Network, VPN, Tor, Proxy, DNS, Firewall, Plugins, Security, Logs, Settings.

## 10. Tela Dashboard

Visão geral do sistema: estado geral, perfil ativo, conexão atual, alertas, recursos.

## 11. Tela Runtime

Controla: daemon, módulos, estado. Ações: iniciar, parar, reiniciar, reload.

## 12. Tela Network

Exibe: gateway, interfaces, rotas, DNS, IPv6. Não exibe: histórico de navegação, destinos acessados.

## 13. Tela VPN

Informações: provider, status, latência, tempo conectado. Ações: conectar, desconectar, trocar perfil.

## 14. Tela Tor

Exibe: status, bootstrap, circuito atual. Ações: iniciar, parar, novo circuito.

## 15. Tela Proxy

Gerenciamento: proxies configurados, latência, disponibilidade.

## 16. Tela DNS

Exibe: servidor ativo, DNSSEC, DoT, estado de segurança.

## 17. Tela Firewall

Exibe: estado, regras, bloqueios. Ações críticas exigem confirmação.

## 18. Tela Plugins

Funções: listar, instalar, atualizar, remover, configurar. Cada plugin mostra: Name, Version, Permissions, Status.

## 19. Tela Security

Painel: ameaças, eventos, integridade, recomendações.

## 20. Tela Logs

Filtros: módulo, nível, período. Níveis: TRACE, DEBUG, INFO, WARN, ERROR, FATAL.

## 21. Tela Configuração

Permite: visualizar, editar, validar, aplicar. Fluxo: Edit -> Validate -> Preview Changes -> Apply.

## 22. Sistema de Eventos

O Dashboard recebe eventos: CONFIG_CHANGED, VPN_CONNECTED, TOR_READY, SECURITY_ALERT, PLUGIN_FAILED. Comunicação: WebSocket.

## 23. Componentes Reutilizáveis

Biblioteca: ghost-ui. Componentes: StatusCard, MetricChart, AlertBox, ConfigEditor, EventTimeline.

## 24. Tema

Suporte: Dark Mode, Light Mode. Configuração: dashboard.theme: dark.

## 25. Responsividade

Suporte: Desktop, Notebook, Tablet. Mobile: Modo leitura.

## 26. Acessibilidade

Requisitos: teclado, contraste, labels, navegação sem mouse.

## 27. Internacionalização

Sistema preparado para: português, inglês, espanhol. Arquivos: locales/pt-BR.json, locales/en-US.json.

## 28. Performance

Metas: carregamento inicial < 2 segundos, baixo uso de memória, atualização em tempo real eficiente.

## 29. Extensão por Plugins

Plugins podem adicionar: páginas, widgets, indicadores. Não podem: substituir telas críticas, alterar segurança.

## 30. Segurança Frontend

Obrigatório: CSP, sanitização, proteção XSS, validação de respostas.

## 31. Testes

Obrigatórios: componentes, navegação, autenticação, permissões, integração API, eventos em tempo real.

## 32. Definition of Done

Este documento será considerado concluído quando:
- Dashboard consumir apenas APIs oficiais;
- possuir controle completo do GhostStack;
- manter paridade com CLI;
- respeitar permissões;
- funcionar em todas plataformas suportadas.
