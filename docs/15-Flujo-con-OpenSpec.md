---
title: Flujo con OpenSpec
tags: [openspec, metodo]
status: borrador
updated: 2026-10-09
---

# Flujo con OpenSpec

Contexto: [[Home]], [[Roadmap]].

## Reparto de responsabilidades
Dos lugares, dos propósitos. No se duplican: se enlazan.

| | `docs/` (este vault, Obsidian) | `openspec/` |
|---|---|---|
| Responde a | **Por qué** y **cómo está armado** | **Qué debe hacer** el sistema, verificable |
| Contenido | Visión, arquitectura, decisiones, despliegue, operación, onboarding | Requisitos con escenarios, propuestas de cambio, tareas |
| Cambia | Poco; con decisiones | Con cada funcionalidad |
| Lo lee | Personas (y asistentes como contexto) | Personas y el agente que implementa |

Regla práctica: si describe **comportamiento observable y comprobable**, va en una spec de OpenSpec. Si explica **motivos, estructura u operación**, va en `docs/`.

## Capacidades (specs) iniciales sugeridas
Una carpeta por capacidad en `openspec/specs/`:

| Capacidad | Resumen | Nota relacionada |
|---|---|---|
| `chat-ingest` | Conexión a Twitch y entrega de mensajes | [[03-Backend-Go]] |
| `message-preprocessing` | Normalización, filtros, deduplicación, tope por usuario | [[03-Backend-Go]] |
| `emotion-classification` | Contrato y calidad mínima del clasificador | [[05-Clasificador-de-Emociones]] |
| `emotion-aggregation` | Ventana, intensidad, tendencia | [[03-Backend-Go]] |
| `realtime-stream-api` | SSE, ticks, modo degradado | [[06-Contratos-e-Integracion]] |
| `thermometer-ui` | Overlay y panel | [[04-Frontend-Vite]] |

## Ejemplo de requisito (formato típico de OpenSpec)
```markdown
## ADDED Requirements

### Requirement: Emisión periódica del estado
El sistema SHALL emitir el estado del termómetro a todos los clientes conectados
en cada tick configurado, aunque no hayan llegado mensajes nuevos.

#### Scenario: Chat sin actividad
- **WHEN** no llega ningún mensaje durante un tick
- **THEN** el sistema emite igualmente un evento `thermometer` con `msgs_per_s` bajo

#### Scenario: Clasificador no disponible
- **WHEN** el clasificador no responde dentro del tiempo límite
- **THEN** el evento incluye `degraded: true`
- **AND** el sistema sigue emitiendo con la lectura de emotes y léxico
```
> El formato exacto puede variar según tu versión de OpenSpec; usa el que genere `/opsx:propose`.

## Contexto del proyecto para OpenSpec
Los archivos de configuración de OpenSpec (según la versión: `openspec/config.yaml` o `openspec/project.md`) deben contener un resumen breve y **apuntar a este vault**. Texto sugerido:

```
Proyecto: termómetro de emociones del chat de Twitch (open source).
Componentes: backend Go (backend/), frontend Vite+TS (frontend/),
clasificador Python (classifier/), contratos gRPC (proto/).
Prioridades: baja latencia mensaje->pantalla, privacidad (no guardar mensajes),
despliegue con Docker, contribución abierta.
Antes de proponer cambios, leer docs/00-Home.md, docs/02-Arquitectura.md
y docs/06-Contratos-e-Integracion.md.
Los contratos (proto/ y evento SSE) no se cambian sin propuesta explícita.
Toda optimización de rendimiento requiere un benchmark que la justifique (docs/07).
```

## Flujo para cada cambio
```
Explore → Propose → Apply → Archive
```
1. **Explore:** se usa al comienzo de una funcionalidad incierta (por ejemplo, elegir modelo, o IRC vs EventSub). Las conclusiones relevantes se vuelcan a un ADR o a [[Decisiones Abiertas]].
2. **Propose:** genera `openspec/changes/<id>/` con propuesta, deltas de spec y tareas. **Se revisa** que la propuesta respete [[Arquitectura]] y los contratos.
3. **Apply:** el agente implementa las tareas; las pruebas de [[11-Pruebas]] son parte de las tareas, no un extra.
4. **Archive:** al cerrar, las specs se actualizan. Si el cambio alteró la arquitectura o una decisión, se actualiza la nota o el ADR correspondiente en `docs/`.

## Reglas de coherencia
- Un PR = un cambio de OpenSpec (salvo trivialidades).
- Si un cambio contradice una nota de `docs/`, se corrige la nota **en el mismo PR**.
- Las notas de `docs/` enlazan a las capacidades; las specs pueden enlazar a las notas como contexto.

## Primeros cambios y prompts de arranque
Ver el orden completo en [[Roadmap]]. Para empezar:

```
/opsx:propose bootstrap-monorepo: crear la estructura de carpetas de docs/02-Arquitectura.md
(backend/, frontend/, classifier/, proto/, deploy/, loadtest/), con Makefile, .gitignore,
.env.example y README mínimo. Sin lógica de negocio todavía.
```

```
/opsx:explore elegir fuente de mensajes de Twitch para v1 (IRC anónimo vs EventSub)
considerando docs/03-Backend-Go.md. Verificar documentación oficial vigente.
```

(Escríbelos con la sintaxis que te mostró `openspec init` para Codex.)
