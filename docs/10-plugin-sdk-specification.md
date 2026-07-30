# GhostStack
## Documento 10 — Plugin SDK Specification
### Versão: 1.0.0
### Status: Normativo

## Dependências

- PRD
- System Architecture
- Threat Model
- Runtime Specification
- ADR-001 Runtime Abstraction
- Domain Model
- Configuration Specification
- Contracts Specification
- OpenAPI Specification
- Event Specification

## 1. Objetivo

Definir o SDK oficial para desenvolvimento de plugins do GhostStack.

O SDK estabelece:

- arquitetura de plugins;
- ciclo de vida;
- APIs disponíveis;
- permissões;
- isolamento;
- distribuição;
- versionamento;
- compatibilidade.

Todo plugin deverá utilizar este SDK.

## 2. Filosofia

Plugins são componentes desacoplados do Core.

O Core nunca conhecerá implementações concretas.

Toda comunicação ocorre através de contratos oficiais.

## 3. Arquitetura

```
Ghost Core
    ↓
Plugin Manager
    ↓
Plugin Runtime
    ↓
Plugin
```

O Plugin Runtime é responsável por carregar, isolar e descarregar plugins.

## 4. Estrutura

Todo plugin deverá possuir:

```
plugin/
├── manifest.yaml
├── plugin.so
├── README.md
├── LICENSE
└── assets/
```

Caso a linguagem não produza bibliotecas compartilhadas (.so), o Plugin Runtime poderá suportar outros formatos definidos pela ABI oficial.

## 5. Manifest

Arquivo obrigatório.

```yaml
id: wireguard
name: WireGuard Provider
version: 1.0.0
author: GhostStack
sdk: 1.0
entry: plugin.so
description: Official WireGuard provider
license: Apache-2.0
```

## 6. Campos Obrigatórios

| Campo | Obrigatório |
|-------|-------------|
| id    | Sim         |
| version | Sim       |
| sdk   | Sim         |
| entry | Sim         |
| author | Sim        |
| description | Sim   |

## 7. Capabilities

O plugin declara exatamente o que implementa.

Exemplo:

```yaml
capabilities:
  - vpn.provider
  - metrics.provider
```

Capabilities desconhecidas invalidam o plugin.

## 8. Permissões

Todo plugin declara permissões explicitamente.

```yaml
permissions:
  - network
  - filesystem
  - runtime
  - events
```

Permissões não declaradas nunca serão concedidas.

## 9. Categorias Oficiais

- vpn.provider
- proxy.provider
- dns.provider
- firewall.provider
- gateway.backend
- runtime.provider
- dashboard.widget
- auth.provider
- storage.provider
- metrics.provider
- notification.provider

## 10. Interface Principal

Todo plugin implementa:

```rust
interface IPlugin {
  initialize(): Promise<void>
  enable(): Promise<void>
  disable(): Promise<void>
  unload(): Promise<void>
  manifest(): PluginManifest
}
```

## 11. Ciclo de Vida

```
Discover
  ↓
Validate
  ↓
Load
  ↓
Initialize
  ↓
Enable
  ↓
Running
  ↓
Disable
  ↓
Unload
```

## 12. Descoberta

Diretórios oficiais:

- `/usr/lib/ghoststack/plugins`
- `~/.local/share/ghoststack/plugins`
- `./plugins`

## 13. Validação

Antes do carregamento:

- manifest válido;
- assinatura (quando exigida);
- compatibilidade de SDK;
- dependências;
- permissões.

Falhas impedem o carregamento.

## 14. Dependências

Manifesto:

```yaml
dependencies:
  - unbound
  - runtime>=1.0
```

Dependências não satisfeitas impedem a ativação.

## 15. Versionamento

SemVer obrigatório.

```
MAJOR.MINOR.PATCH
```

Compatibilidade é determinada pela versão do SDK.

## 16. ABI

O SDK define uma ABI estável.

Mudanças incompatíveis exigem nova versão major.

## 17. Isolamento

Plugins não compartilham memória diretamente.

Toda comunicação ocorre por:

- contratos;
- Event Bus;
- APIs oficiais.

## 18. Sandboxing

Plugins poderão ser executados em sandbox.

Backends oficiais:

- processo nativo;
- seccomp (Linux);
- namespaces (Linux).

Backends futuros:

- WASI;
- WebAssembly;
- containers.

## 19. APIs Disponíveis

Plugins podem acessar:

- Configuration API
- Event API
- Logging API
- Metrics API
- Runtime API
- Secret API
- Health API

## 20. APIs Proibidas

Plugins nunca poderão:

- modificar memória do Core;
- acessar arquivos internos sem permissão;
- substituir contratos oficiais;
- interceptar chamadas internas.

## 21. Eventos

Plugins podem publicar:

- PluginEvent

E consumir qualquer evento público autorizado.

## 22. Logging

Todo plugin recebe um logger próprio.

```
logger.info()
logger.error()
logger.warn()
```

Todos os logs devem incluir o identificador do plugin.

## 23. Configuração

Plugins possuem configuração própria.

```
plugins/
  wireguard/
    config.yaml
```

Nunca modificar ghost.yaml.

## 24. Secrets

Secrets são obtidos apenas pela Secret API.

Nunca por leitura direta de arquivos internos.

## 25. Atualizações

Plugin Manager suporta:

- instalar;
- atualizar;
- remover;
- rollback.

## 26. Assinaturas

Plugins podem ser assinados digitalmente.

Arquivos:

- plugin.so
- plugin.sig

O Core pode ser configurado para aceitar apenas plugins assinados.

## 27. Repositórios

Fontes oficiais:

- Official Repository
- Local Repository

Futuras:

- Enterprise Repository
- Custom Repository

## 28. Compatibilidade

O Core rejeita plugins incompatíveis com:

- SDK;
- ABI;
- versão mínima do Runtime.

## 29. Hot Reload

Plugins podem suportar:

```
Disable
  ↓
Reload
  ↓
Enable
```

Sem reiniciar o Runtime.

## 30. Recursos

Todo plugin deve declarar recursos esperados.

Exemplo:

```yaml
resources:
  memory: 64MB
  cpu: low
  network: true
```

Esses valores auxiliam planejamento e observabilidade, mas não substituem mecanismos de limitação.

## 31. Métricas

Plugins podem expor:

- tempo de inicialização;
- uso de memória;
- uso de CPU;
- erros;
- chamadas de API.

## 32. Segurança

Plugins nunca poderão:

- elevar privilégios;
- acessar secrets não autorizados;
- executar código fora da sandbox configurada;
- alterar configurações globais sem autorização.

## 33. Testes

Todo plugin oficial deverá possuir:

- testes unitários;
- testes de contrato;
- testes de integração;
- testes de compatibilidade de SDK.

## 34. Distribuição

Pacotes oficiais:

```
ghost-plugin-wireguard-1.0.0.gsp
```

Formato `.gsp` (GhostStack Plugin Package) é reservado para distribuição oficial.

## 35. Ferramentas

CLI oficial:

```
ghost plugin create
ghost plugin build
ghost plugin test
ghost plugin publish
ghost plugin install
```

## 36. Requisitos Normativos

O SDK deverá:

- manter ABI estável por versão major;
- impedir carregamento de plugins inválidos;
- validar manifestos;
- suportar atualização segura;
- isolar plugins do Core.

## 37. Definition of Done

Este documento será considerado concluído quando:

- o ciclo de vida dos plugins estiver formalmente especificado;
- o formato do manifesto estiver definido;
- as permissões e capacidades forem verificáveis;
- o Plugin Runtime puder ser implementado sem ambiguidades;
- terceiros puderem desenvolver plugins compatíveis utilizando apenas esta especificação.
