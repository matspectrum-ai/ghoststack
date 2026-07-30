# GhostStack — Threat Model & Security Architecture

Versão: 0.1.0
Status: Draft
Dependência: PRD + System Architecture

## 1. Objetivo

Este documento define: modelo de ameaças, objetivos de segurança, limites do sistema, garantias, não garantias e mitigações obrigatórias.

Este documento é normativo. Toda implementação deverá obedecer suas regras.

## 2. Filosofia

GhostStack não promete anonimato absoluto. GhostStack reduz a superfície de ataque, minimiza vazamentos e fornece uma plataforma auditável para roteamento seguro.

## 3. Objetivos de Segurança

SG-001: Todo tráfego deve sair apenas pelo Gateway.
SG-002: Nenhum processo pode acessar Internet diretamente.
SG-003: Falhas devem resultar em estado seguro (fail closed).
SG-004: Detectar vazamentos automaticamente.
SG-005: Configuração reproduzível.
SG-006: Logs auditáveis.
SG-007: Zero telemetria obrigatória.

## 4. Ativos Protegidos

Rede: IP público, IP local, rotas, DNS (consultas, cache, resolver).
Identidade: perfil de navegação, sessões, cookies.
Configuração: ghost.yaml, certificados, credenciais VPN, credenciais Proxy.
Estado: plugins, banco SQLite, histórico, eventos.

Observação: O GhostStack protege o tráfego de rede. O gerenciamento de perfis e identidades de navegador é responsabilidade do navegador utilizado.

## 5. Trust Boundaries

Usuário -> CLI -> Ghost Daemon -> Gateway -> Internet

Cada fronteira deve validar entradas. Nunca confiar em dados externos.

## 6. Adversários

Nível 1 — Observador local (Wi-Fi público, ISP): observar tráfego, observar DNS.
Nível 2 — Servidor remoto (site acessado): observar IP de origem, fingerprint HTTP/TLS.
Nível 3 — Proxy malicioso: modificar tráfego, registrar conexões.
Nível 4 — VPN comprometida: observar metadados.
Nível 5 — Exit Node Tor malicioso: observar tráfego sem criptografia fim a fim.

Fora do escopo: comprometimento do kernel, malware com privilégios de root, ataques físicos, firmware malicioso, vulnerabilidades do navegador, engenharia social, comprometimento do provedor VPN ou da rede Tor em nível global.

## 7. STRIDE

Spoofing: TLS, certificados, validação.
Tampering: HTTPS, WireGuard, OpenVPN.
Repudiation: logs, IDs, timestamps.
Information Disclosure: DNS, Kill Switch, Gateway.
Denial of Service: Fail Closed.
Elevation of Privilege: seccomp, AppArmor, capacidades mínimas.

## 8. LINDDUN

Linkability: mitigar correlação de sessões quando possível.
Identifiability: ocultar IP através do Gateway configurado.
Non-Repudiation: não registrar mais dados do que o necessário.
Detectability: reduzir exposição de tráfego observável.
Disclosure: criptografar canais suportados.
Unawareness: toda configuração deve ser transparente ao usuário.
Non-Compliance: compatível com GDPR e LGPD quando aplicável.

## 9. Ameaças

DNS Leak: consulta DNS fora do Gateway. Severidade: CRÍTICA. Mitigação: Unbound obrigatório + teste automático.
IPv6 Leak: mitigação: desabilitar ou rotear IPv6 conforme configuração. Teste obrigatório.
WebRTC Leak: específico de navegadores. GhostStack fornecerá diagnósticos, mas mitigação depende da configuração do navegador.
Direct Connection: aplicação ignora Gateway. Mitigação: Kill Switch.
VPN Down: mitigação: Fail Closed.
Tor Down: mitigação: bloquear saída até nova decisão do usuário.
Proxy Down: mitigação: remover rota. Nunca sair diretamente.

## 10. Objetivos de Privacidade

GhostStack busca minimizar: exposição do IP de origem, vazamentos de DNS, vazamentos de rotas, conexões diretas não autorizadas.

GhostStack não garante: invisibilidade, impossibilidade de rastreamento ou anonimato absoluto.

## 11. Kill Switch

Obrigatório. Estados: VPN OFF -> BLOCK. Nunca: VPN OFF -> Internet direta.

## 12. DNS

Obrigatório. Nunca usar ISP DNS sem autorização explícita.

## 13. Logs

Permitidos: Container Started, VPN Connected, Health, Errors.
Proibidos: senhas, tokens, queries DNS completas, payloads HTTP, conteúdo do tráfego.

## 14. Segredos

Nunca armazenar em: Git, YAML. Utilizar: Docker Secrets, arquivos protegidos por permissões, variáveis de ambiente quando apropriado.

## 15. Containers

Todos deverão: Read Only RootFS, No New Privileges, Seccomp, AppArmor, Drop Capabilities. Sempre que tecnicamente possível.

## 16. Dependências

Toda dependência deverá possuir: versão fixa, checksum, SBOM (Software Bill of Materials).

## 17. Atualizações

Atualizações nunca serão automáticas por padrão. Usuário decide.

## 18. Supply Chain

Obrigatório: imagens assinadas quando disponíveis, hashes, SBOM, verificação de integridade.

## 19. Configuração Segura

Padrão:
mode: direct
kill_switch: true
dns:
  secure: true
ipv6:
  mode: disabled
telemetry: false

## 20. Health Checks

Cada módulo expõe: Health, Version, Status, Latency.

## 21. Auditoria

Comando: ghost doctor. Executará verificações: Gateway ativo, DNS, VPN, Tor, Firewall, Containers, Rotas, Estado do Kill Switch, Configuração consistente.

## 22. Matriz de Riscos

Ameaça | Severidade | Mitigação
DNS Leak | Crítica | Unbound + testes
IPv6 Leak | Alta | Desabilitar ou rotear IPv6
VPN Down | Crítica | Kill Switch
Proxy Down | Alta | Fail Closed
Tor Down | Alta | Bloqueio ou troca controlada
Configuração inválida | Média | Validação antes do startup
Container comprometido | Alta | Isolamento + privilégios mínimos
Atualização maliciosa | Alta | Assinaturas + hashes + SBOM

## 23. Casos de Segurança

SC-001: VPN cai -> Kill Switch ativa, nenhum pacote sai, evento registrado.
SC-002: DNS Leak -> Detectado, alerta emitido, ghost doctor retorna falha.
SC-003: Tor indisponível -> sessões não migram automaticamente para Internet direta, usuário escolhe a ação.
SC-004: Container do Gateway falha -> Core tenta recuperação, se não recuperar, estado final é FAIL CLOSED.

## 24. Limitações Declaradas

GhostStack não garante: anonimato absoluto, resistência a adversários com comprometimento do sistema operacional, proteção contra malware instalado pelo usuário, proteção contra fingerprinting de navegador, proteção contra login voluntário em contas pessoais, impossibilidade de correlação por comportamento do usuário.

## 25. Requisitos Normativos

Todos os módulos devem: implementar fail closed, implementar health checks, validar configuração antes da execução, emitir eventos estruturados, não registrar informações sensíveis, executar com o menor conjunto possível de privilégios.

## Definition of Done

Este documento será considerado completo quando:
- O modelo de ameaças estiver documentado.
- Todos os ativos protegidos estiverem identificados.
- As limitações estiverem explicitadas.
- Todas as ameaças relevantes possuírem mitigação ou justificativa.
- Os requisitos de segurança puderem ser convertidos em testes automatizados.
