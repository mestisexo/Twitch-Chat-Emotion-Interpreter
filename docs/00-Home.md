---
title: Home
tags: [moc, proyecto]
status: vivo
updated: 2026-10-09
---

# 🌡️ Termómetro de Emociones del Chat

> [!note] Nombre provisional
> El nombre definitivo está pendiente (ver [[Decisiones-Abiertas]]). En estos documentos usamos "el termómetro".

**En una frase:** un sistema open source que lee el chat de Twitch en tiempo real, clasifica la emoción de los mensajes y se la muestra al streamer como un termómetro, para que perciba el ánimo del chat aunque no pueda leer todos los mensajes.

## Mapa del vault

### Por qué y qué
- [[01-Vision-y-Objetivos]]: problema, objetivo, alcance, métricas de éxito
- [[16-Roadmap]]: hitos y su relación con los cambios de OpenSpec
- [[17-Glosario]]

### Cómo está construido
- [[02-Arquitectura]]: visión general y flujo de datos
- [[03-Backend-Go]]
- [[04-Frontend-Vite]]
- [[05-Clasificador-de-Emociones]]
- [[06-Contratos-e-Integracion]]: gRPC, eventos SSE, versionado
- [[08-Datos-y-Almacenamiento]]

### Calidad y rendimiento
- [[07-Latencia-y-Rendimiento]]: presupuesto de latencia y benchmarks
- [[11-Pruebas]]
- [[12-Observabilidad]]

### Producción y comunidad
- [[09-Despliegue-Docker]]
- [[10-CI-CD-GitHub]]
- [[13-Seguridad-y-Privacidad]]
- [[14-Onboarding-y-Contribuir]]

### Método de trabajo
- [[15-Flujo-con-OpenSpec]]: cómo se reparten la responsabilidad estos docs y `openspec/`
- Decisiones: [[Decisiones-Abiertas]] y carpeta de ADRs ([[0000-plantilla-adr]])

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
