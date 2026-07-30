# GhostStack Architecture

## Overview

GhostStack is a privacy orchestration platform built with Go and TypeScript.

## Components

- Core: Go runtime, CLI, providers.
- Dashboard: React + Vite.
- Plugins: extensible via Go/TS SDK.

## Data Flow

1. CLI sends commands to Core.
2. Core manages providers and networking.
3. Dashboard polls API for state.
4. Events flow through Event Bus.
