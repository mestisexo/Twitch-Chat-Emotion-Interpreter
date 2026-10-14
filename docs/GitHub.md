---
title: GitHub
tags:
  - ci
  - github
  - devops
status: borrador
updated: 2026-10-09
---

# GitHub

Contexto: [[09-Despliegue-Docker]], [[11-Pruebas]].


## 🍀 Flujo de ramas

- La rama `main` siempre protegida.
- Ramas cortas por cambio. 
- Pull Request obligatorio para mis cambios, los externos sí obligatorio con revisión.
- Se fusiona con squash merge. 
- El título del Pull Request sigue Conventional Commits (ver a continuación)

| Tipo       | Cuando                                   | Version               |
| ---------- | ---------------------------------------- | --------------------- |
| `feat`     | Nueva funcionalidad                      | minor (0.1.0 → 0.2.0) |
| `fix`      | Corrección de bug                        | patch (0.1.0 → 0.1.1) |
| `docs`     | Solo documentación                       | no sube               |
| `refactor` | Reestructurar sin cambiar comportamiento | no sube               |
| `perf`     | Mejora de rendimiento                    | patch                 |
| `test`     | Agregar o corregir tests                 | no sube               |
| `build`    | Sistema de build o dependencias          | no sube               |
| `ci`       | Workflows de CI                          | no sube               |
| `chore`    | Mantenimiento que no encaja en otro      | no sube               |
- Los Pull Request referencian su cambio de OpenSpec.


## 🌋 GitHub Actions

| Workflow       | Cuándo             | Qué hace                                |
| -------------- | ------------------ | --------------------------------------- |
| `pr-title.yml` | Existe un nuevo PR | Revisión semantica Conventional Commits |


## 🌩️ Archivos de comunidad

- Ver detalles en [[14-Onboarding-y-Contribuir]].

| Archivo                            | Función |
| ---------------------------------- | ------- |
| `README.md`                        |         |
| `CONTRIBUTING.md`                  |         |
| `CODE_OF_CONDUCT.md`               |         |
| `SECURITY.md`                      |         |
| `LICENSE`                          |         |
| `.github/ISSUE_TEMPLATE/`          |         |
| `.github/PULL_REQUEST_TEMPLATE.md` |         |


## 🪂 Licencia

Pendiente de decidir ([[Decisiones Abiertas]]). Considerar también la licencia de modelos y datasets que se distribuyan.
