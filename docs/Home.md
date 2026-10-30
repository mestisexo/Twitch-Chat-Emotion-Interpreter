---
title: Home
tags: [moc, proyecto]
status: vivo
updated: 2026-10-09
---

# Home


## 🎇 Definición

- [[Meta]] 


## Cómo está construido

- [[Arquitectura]]: visión general y flujo de datos
- [[Backend]]
- [[Frontend]]
- [[05-Clasificador-de-Emociones]]
- [[06-Contratos-e-Integracion]]: gRPC, eventos SSE, versionado
- [[08-Datos-y-Almacenamiento]]

### Calidad y rendimiento
- [[07-Latencia-y-Rendimiento]]: presupuesto de latencia y benchmarks
- [[11-Pruebas]]
- [[12-Observabilidad]]

### Producción y comunidad
- [[09-Despliegue-Docker]]
- [[GitHub]]
- [[13-Seguridad-y-Privacidad]]
- [[14-Onboarding-y-Contribuir]]

### Método de trabajo
- [[15-Flujo-con-OpenSpec]]: cómo se reparten la responsabilidad estos docs y `openspec/`
- Decisiones: [[Decisiones Abiertas]] y carpeta de ADRs ([[0000-plantilla-adr]])

## Decisiones tomadas (propuestas)
| ADR | Decisión | Estado |
|---|---|---|
| [[0001-monorepo]] | Un solo repositorio con tres componentes | Propuesta |
| [[0002-sse-para-tiempo-real]] | SSE hacia el navegador/OBS | Propuesta |
| [[0003-grpc-para-clasificador]] | gRPC por lotes entre backend y clasificador | Propuesta |
| [[0004-clasificador-servicio-separado]] | Clasificador como servicio aparte (Python en v1) | Propuesta |
| [[0005-no-persistir-mensajes-crudos]] | No guardar el texto de los mensajes | Propuesta |

## Cómo abrir este vault
1. En Obsidian: *Open folder as vault* y elige la carpeta `docs/` del repositorio.
2. Los diagramas usan Mermaid (Obsidian lo renderiza de forma nativa).
3. Añade a `.gitignore`: `docs/.obsidian/workspace*.json` para no versionar el estado personal de cada persona.
4. Convención: nombres de archivo sin tildes ni espacios; enlaces con wikilinks (doble corchete); cada nota lleva frontmatter con `tags`, `status` y `updated`.
