---
title: Observabilidad
tags: [observabilidad, metricas]
status: borrador
updated: 2026-10-09
---

# Observabilidad

Contexto: [[07-Latencia-y-Rendimiento]], [[09-Despliegue-Docker]].

## Métricas (Prometheus, en `/metrics`)
| Métrica | Tipo | Para qué |
|---|---|---|
| `messages_received_total` | contador | Volumen de ingesta |
| `messages_dropped_total{reason}` | contador | Descartes por saturación o filtros |
| `fastpath_resolved_total` | contador | Cuánto evita el modelo |
| `cache_hits_total` / `cache_misses_total` | contador | Eficacia de la caché |
| `batch_size` | histograma | Tamaño real de los lotes |
| `classifier_latency_seconds` | histograma | Latencia gRPC + inferencia |
| `classifier_errors_total` | contador | Fallos del clasificador |
| `e2e_latency_seconds` | histograma | **Mensaje → estado emitido** (métrica principal) |
| `queue_depth{stage}` | gauge | Detectar cuellos de botella |
| `sse_clients` | gauge | Clientes conectados |
| `degraded_mode` | gauge | 1 si funciona sin clasificador |

El clasificador expone además latencia por lote, tamaño de lote y uso de memoria.

## Logs
- Estructurados (JSON), con nivel configurable.
- **No registrar el texto de los mensajes** en niveles normales ([[13-Seguridad-y-Privacidad]]).
- Identificador de correlación por lote para seguir un problema entre servicios.

## Trazas (opcional)
OpenTelemetry si se necesita ver el recorrido de un lote entre backend y clasificador. No es requisito de v1.

## Tableros
Un tablero de Grafana de ejemplo en `deploy/observability/`, activable con el perfil `observability` de Compose.

## Alertas sugeridas (para despliegues propios)
- `e2e_latency_seconds` p95 por encima del objetivo durante varios minutos.
- `degraded_mode == 1` sostenido.
- `queue_depth` creciendo sin parar.
- Reconexiones frecuentes a Twitch.
