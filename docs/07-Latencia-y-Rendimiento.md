---
title: Latencia y rendimiento
tags: [rendimiento, latencia, benchmarks]
status: borrador
updated: 2026-10-09
---

# Latencia y rendimiento

Contexto: [[Arquitectura]]. Pruebas: [[11-Pruebas]]. Métricas: [[12-Observabilidad]].

> [!important] Todas las cifras de esta nota son **hipótesis de partida**
> Se validan o corrigen con los benchmarks. Cuando haya datos, se actualizan aquí y se cita la medición.

## Dos conceptos que no hay que mezclar
- **Latencia de pipeline:** tiempo desde que el mensaje llega al backend hasta que la aguja se mueve en pantalla. Se optimiza con ingeniería.
- **Suavizado (ventana):** la ventana de agregación hace que la aguja sea estable pero inevitablemente "recuerde" unos segundos. Es una decisión de **producto**, ajustable (`WINDOW_SECONDS`, `HALF_LIFE_SECONDS`).

Un buen termómetro reacciona rápido a un cambio real (por ejemplo un pico de hype) y no tiembla con ruido.

## Objetivos iniciales (hipótesis)
| Métrica | Objetivo |
|---|---|
| Mensaje → aguja, p50 | < 500 ms |
| Mensaje → aguja, p95 | < 1 s |
| Reacción visible ante un cambio brusco de ánimo | pocos segundos (depende del suavizado elegido) |
| Capacidad sostenida por canal | cientos de mensajes por segundo, con muestreo si hace falta |
| Clientes SSE simultáneos por instancia | decenas (el overlay y el panel; no es un servicio masivo) |

## Presupuesto de latencia por etapa (p95, hipótesis)
| Etapa | Presupuesto |
|---|---|
| Ingesta y normalización | < 5 ms |
| Fast path y caché | < 1 ms |
| Espera del micro-batch | ≤ 150 ms |
| Transporte gRPC (red local) | < 5 ms |
| Inferencia del lote (CPU, modelo pequeño) | < 100 ms |
| Agregación | < 5 ms |
| Espera al siguiente tick | ≤ 500 ms (media 250 ms) |
| Push SSE y render | < 50 ms |

El mayor consumo está en el batch y el tick, que son **decisiones de diseño**, no costos de cómputo. Reducirlos mejora la latencia a cambio de más llamadas y más actualizaciones.

> La latencia entre Twitch y nuestro backend no se controla; se mide aparte si es posible.

## Palancas de optimización (en orden de impacto esperado)
1. **No clasificar todo:** fast path por emotes/léxico y caché de textos repetidos ([[05-Clasificador-de-Emociones]]).
2. **Muestreo bajo saturación:** el termómetro es estadístico; perder un porcentaje de mensajes en picos extremos apenas cambia el resultado.
3. **Micro-batching** con límites de tamaño y de espera.
4. **Modelo ligero y cuantizado** (ONNX/INT8).
5. **Evitar copias y asignaciones** en el camino caliente del backend (perfilar con `pprof` antes de tocar nada).
6. **Un solo tick** que emite a todos los clientes, en vez de trabajo por cliente.
7. **Conexiones persistentes** gRPC reutilizadas.

Regla: **medir primero, optimizar después**. Ninguna optimización entra sin un benchmark que la justifique.

## Cómo medir
### Marcas de tiempo
Cada mensaje lleva `t_ingest` desde la llegada. Al emitir el estado se calcula la edad del mensaje más reciente que contribuyó. Se publica como histograma (`e2e_latency_seconds`). En el cliente se puede medir además la diferencia con `ts` del evento.

### Escenarios de carga (en `loadtest/`)
| Escenario | Descripción |
|---|---|
| Base | ~20 msg/s constantes, mezcla realista |
| Ocupado | ~200 msg/s |
| Pico | ~1000 msg/s durante 30 s (prueba de muestreo y backpressure) |
| Cambio de ánimo | Chat neutral que pasa abruptamente a hype: mide tiempo de reacción |
| Muchos clientes | 50+ conexiones SSE simultáneas |
| Clasificador caído | Verifica modo degradado y recuperación |
| Clasificador lento | Inyecta latencia para validar timeouts y colas |

### Herramientas
- **Simulador/replay de chat** propio (genera mensajes con distribución realista de emotes, repeticiones y ráfagas; también reproduce grabaciones anonimizadas).
- **Go `testing.B` + `benchstat`** para micro-benchmarks de `pipeline` y `aggregate`.
- **k6** (o similar) para la capa SSE con muchos clientes.
- **`pprof`** para CPU y memoria.
- Benchmark de inferencia del clasificador por separado: latencia por lote según tamaño de lote.

### Entorno de medición
Declarar siempre: CPU, RAM, versión de imagen y parámetros. Una cifra sin entorno no sirve. Guardar resultados versionados en `loadtest/results/`.

## Criterios de aceptación y regresión
- Los objetivos de la tabla se cumplen en los escenarios Base, Ocupado y Pico (con muestreo declarado).
- CI ejecuta micro-benchmarks y **falla si hay una regresión mayor al umbral acordado** frente a la rama principal.
- Las pruebas de carga completas corren de forma programada (nocturna) y antes de cada release, no en cada PR.

## Plan de trabajo
Se implementa por etapas como cambios de OpenSpec: ver [[Roadmap]] (`add-load-testing`).
