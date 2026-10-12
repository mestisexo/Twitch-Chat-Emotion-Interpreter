---
title: CI/CD y GitHub
tags: [ci, github, devops]
status: borrador
updated: 2026-10-09
---

# CI/CD y GitHub

Contexto: [[09-Despliegue-Docker]], [[11-Pruebas]].

## Flujo de ramas
- `main` siempre desplegable y protegida.
- Ramas cortas por cambio; Pull Request obligatorio con revisión y CI verde.
- Commits con **Conventional Commits** (`feat:`, `fix:`, `docs:`...) para generar changelog y versiones.
- Cada PR significativo referencia su cambio de OpenSpec ([[15-Flujo-con-OpenSpec]]).

## Workflows (GitHub Actions)
| Workflow | Cuándo | Qué hace |
|---|---|---|
| `backend.yml` | PR / push a `main` con cambios en `backend/` o `proto/` | gofmt, vet, golangci-lint, `go test -race`, benchmarks |
| `frontend.yml` | Cambios en `frontend/` | lint, typecheck, Vitest, build |
| `classifier.yml` | Cambios en `classifier/` | lint (ruff), pytest, evaluación sobre el conjunto de oro |
| `proto.yml` | Cambios en `proto/` | `buf lint` y `buf breaking` |
| `integration.yml` | PR / nocturno | Levanta compose con simulador de chat y corre pruebas de integración y E2E |
| `perf-nightly.yml` | Programado y pre-release | Escenarios de carga ([[07-Latencia-y-Rendimiento]]); publica resultados |
| `docker.yml` | Push a `main` y tags | Construye imágenes multi-arch, escanea (por ejemplo con Trivy) y publica a GHCR en tags |
| `docs.yml` | Cambios en `docs/` | Revisión de enlaces rotos y formato |
| `release.yml` | Tag `vX.Y.Z` | Crea release con changelog y artefactos |

Filtros por ruta (`paths:`) para no correr todo en cada cambio.

## Calidad y seguridad automatizadas
- **Dependabot** para Go, npm, pip y GitHub Actions.
- Escaneo de dependencias e imágenes.
- Análisis estático (por ejemplo CodeQL) y detección de secretos.
- Acciones fijadas por versión/commit.
- Caché de dependencias y de capas Docker para acelerar CI.

## Branch protection (sugerido)
- Requerir PR, al menos 1 aprobación, checks obligatorios, historial lineal.
- `CODEOWNERS` por componente para dirigir revisiones.

## Archivos de comunidad
`README.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, `LICENSE`, `.github/ISSUE_TEMPLATE/`, `.github/PULL_REQUEST_TEMPLATE.md`. Detalle en [[14-Onboarding-y-Contribuir]].

## Etiquetas de issues recomendadas
`good first issue`, `help wanted`, `bug`, `enhancement`, `docs`, `performance`, `classifier`, `backend`, `frontend`, `needs-spec` (requiere propuesta de OpenSpec antes de implementar).

## Licencia
Pendiente de decidir ([[Decisiones-Abiertas]]). Considerar también la licencia de modelos y datasets que se distribuyan.
