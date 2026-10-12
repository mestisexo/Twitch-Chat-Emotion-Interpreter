---
title: "ADR-0001 Monorepo"
tags: [adr]
status: propuesta
updated: 2026-10-09
---

# ADR-0001: Un solo repositorio (monorepo)

- **Estado:** propuesta
- **Fecha:** 2026-10-09

## Contexto
El sistema tiene tres componentes (Go, Vite, Python) unidos por contratos. El proyecto será open source y queremos que quien llegue lo ejecute fácil. Ver [[02-Arquitectura]].

## Opciones
1. **Monorepo:** todo junto. Cambios de contrato atómicos, un solo lugar para issues, docs y OpenSpec, una sola clonación.
2. **Un repo por componente:** despliegues independientes, pero cambios coordinados en varios PRs, contratos desincronizados y onboarding más difícil.

## Decisión
Monorepo, con carpetas por componente y filtros de rutas en CI ([[10-CI-CD-GitHub]]).

## Consecuencias
- (+) Cambios que tocan contrato y consumidores van en un solo PR.
- (+) Un único `openspec/` y un único vault `docs/`.
- (−) CI debe filtrar por rutas para no ser lento.
- (−) Permisos granulares por componente se resuelven con `CODEOWNERS`, no con repos separados.

## Revisión
Si algún componente necesita ciclo de release totalmente independiente.
