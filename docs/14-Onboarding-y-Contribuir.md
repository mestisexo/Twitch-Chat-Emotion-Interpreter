---
title: Onboarding y cómo contribuir
tags: [comunidad, onboarding, opensource]
status: borrador
updated: 2026-10-09
---

# Onboarding y cómo contribuir

Contexto: [[Home]]. Método de trabajo: [[15-Flujo-con-OpenSpec]].

Objetivo: que una persona nueva **ejecute el proyecto en menos de 15 minutos** y sepa dónde empezar a aportar.

## Ruta de bienvenida
1. Leer [[Meta]] (5 min) y [[Arquitectura]] (10 min).
2. Levantar el sistema con el simulador de chat (sin necesidad de Twitch):
   ```
   git clone <repo>
   cp .env.example .env
   docker compose -f compose.dev.yaml up
   ```
3. Abrir el overlay y ver reaccionar el termómetro.
4. Elegir un issue `good first issue`.

## Requisitos de desarrollo
| Componente | Herramientas |
|---|---|
| Backend | Go (versión indicada en `go.mod`), golangci-lint |
| Frontend | Node LTS, npm o pnpm |
| Clasificador | Python (versión indicada en `pyproject.toml`), uv o pip |
| Contratos | `buf` |
| General | Docker y Docker Compose, Node (para la CLI de OpenSpec) |

Un `Makefile` (o `justfile`) con comandos comunes: `make dev`, `make test`, `make lint`, `make bench`, `make up`.

## Cómo proponer un cambio
1. **Cambio pequeño** (typo, bug evidente): abre un PR directo.
2. **Cambio que altera comportamiento o contratos:** primero una propuesta con OpenSpec (`/opsx:propose`), discutida en un issue o PR. Ver [[15-Flujo-con-OpenSpec]].
3. **Decisión de arquitectura:** un ADR usando [[0000-plantilla-adr]].

## Convenciones
- Commits: Conventional Commits.
- PRs pequeños, con descripción del *por qué*, pruebas y enlace al cambio de OpenSpec.
- Estilo: lo que imponen los linters de cada componente (CI lo verifica).
- Documentación: si cambias comportamiento, actualiza `docs/` y las specs.
- Idioma: **pendiente de decidir** (español, inglés o ambos); ver [[Decisiones Abiertas]].

## Roles y gobernanza (inicial)
- **Maintainers:** revisan y fusionan; definen la dirección.
- **Contribuidores:** cualquiera que aporte código, docs, datos o ideas.
- Código de conducta obligatorio (`CODE_OF_CONDUCT.md`).

## Cosas fáciles para empezar
- Ampliar el diccionario de emotes y léxico.
- Mejorar mensajes de error y documentación.
- Añadir escenarios al simulador de chat.
- Probar el despliegue en un sistema nuevo y reportar fricciones.
- Etiquetar ejemplos para el conjunto de oro (con las reglas de consentimiento de [[13-Seguridad-y-Privacidad]]).

## Checklist de archivos de comunidad (para el repo)
- [ ] `README.md` con demo y arranque rápido
- [ ] `CONTRIBUTING.md`
- [ ] `CODE_OF_CONDUCT.md`
- [ ] `SECURITY.md`
- [ ] `LICENSE`
- [ ] Plantillas de issue y PR
- [ ] `AGENTS.md` con instrucciones para asistentes de código (Codex, etc.)
