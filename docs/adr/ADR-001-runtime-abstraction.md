# ADR-001 — Runtime Abstraction

Status: Accepted

## Contexto

Inicialmente o GhostStack foi projetado utilizando Docker como runtime principal para execução dos serviços.

Durante a revisão arquitetural identificou-se que essa decisão cria um acoplamento desnecessário entre o Core e uma tecnologia específica.

O objetivo do GhostStack é orquestrar a rede do sistema operacional, não containers.

## Decisão

O Ghost Core não conhecerá Docker, Podman ou qualquer outro runtime.

Ele dependerá exclusivamente da interface:

IRuntimeProvider

Implementações oficiais:

- NativeRuntimeProvider
- DockerRuntimeProvider
- PodmanRuntimeProvider

Implementações futuras:

- ContainerdRuntimeProvider
- KubernetesRuntimeProvider
- LXCRuntimeProvider

## Consequências

- Docker deixa de ser requisito.
- Native Runtime passa a ser a implementação padrão.
- Docker torna-se opcional.
- A CLI permanece inalterada.
- O Daemon passa a resolver dinamicamente o Runtime Provider.

## Justificativa

- menor acoplamento;
- maior portabilidade;
- melhor desempenho;
- melhor integração com Linux;
- maior facilidade de testes;
- arquitetura preparada para expansão.
