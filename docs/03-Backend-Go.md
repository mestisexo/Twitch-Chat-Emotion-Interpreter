---
title: Backend (Go)
tags: [backend, go]
status: borrador
updated: 2026-10-09
---

# Backend (Go)

Contexto: [[Arquitectura]]. Contratos: [[06-Contratos-e-Integracion]].

## Responsabilidades
1. Conectarse al chat de Twitch y mantener la conexión (reconexión con backoff).
2. Preprocesar y filtrar mensajes.
3. Resolver por fast path lo que se pueda, y enviar el resto al clasificador por lotes.
4. Agregar resultados en ventanas y calcular el estado del termómetro.
5. Publicar el estado por SSE y exponer endpoints operativos.
6. Servir el frontend compilado (archivos estáticos embebidos) en producción.

## Fuente de mensajes de Twitch
| Opción | Cuándo | Nota |
|---|---|---|
| IRC de solo lectura (anónimo) | v1, uso por el propio streamer, sin registrar app | Más simple. |
| EventSub (WebSocket) | v2, cuando se requiera OAuth o varios canales | Es el camino recomendado para integraciones nuevas. |

> [!warning] Verificar
> Confirmar en la documentación oficial vigente de Twitch qué métodos están soportados y sus límites antes de implementar. Se abordará como un cambio de OpenSpec (`add-twitch-chat-ingest`).

La ingesta se esconde detrás de una interfaz (`ChatSource`) para poder cambiar de IRC a EventSub o a un **simulador de chat** (clave para pruebas, ver [[11-Pruebas]]).

## Estructura de paquetes propuesta
```
backend/
├── cmd/server/          # main, wiring, arranque y apagado ordenado
├── internal/
│   ├── twitch/          # ChatSource: IRC / EventSub / replay
│   ├── pipeline/        # normalización, filtros, fast path, caché, batcher
│   ├── classifier/      # cliente gRPC + fallback
│   ├── aggregate/       # ventana deslizante, intensidad, tendencia
│   ├── api/             # SSE, health, readiness
│   ├── config/          # variables de entorno
│   └── metrics/         # Prometheus
├── gen/                 # código generado desde proto/
└── Dockerfile
```

## Modelo de concurrencia
- Una goroutine por etapa, conectadas con **canales con buffer acotado**.
- Política ante saturación: descartar el mensaje más antiguo y contar el descarte en una métrica (`messages_dropped_total`).
- Un `ticker` único emite el estado cada `TICK_INTERVAL` (por defecto 500 ms).
- Cierre ordenado con `context.Context`: drenar colas, cerrar conexiones.

## Configuración (variables de entorno)
| Variable | Por defecto | Descripción |
|---|---|---|
| `TWITCH_CHANNEL` | (obligatoria) | Canal a escuchar |
| `CLASSIFIER_ADDR` | `classifier:50051` | Dirección gRPC |
| `BATCH_MAX_SIZE` | `32` | Máximo de mensajes por lote |
| `BATCH_MAX_WAIT_MS` | `150` | Espera máxima antes de enviar un lote |
| `WINDOW_SECONDS` | `10` | Ventana de agregación |
| `HALF_LIFE_SECONDS` | `4` | Decaimiento del peso por antigüedad |
| `TICK_INTERVAL_MS` | `500` | Frecuencia de emisión |
| `QUEUE_SIZE` | `1024` | Tamaño de colas entre etapas |
| `LOG_LEVEL` | `info` | Nivel de logs |

> Valores iniciales: son hipótesis a ajustar con los benchmarks de [[07-Latencia-y-Rendimiento]].

## Agregación (propuesta)
Para cada emoción *e* en el instante *t*:

`score_e(t) = Σ ( w_i · p_i,e ) / Σ w_i`, con `w_i = 0.5^((t − t_i) / HALF_LIFE)`

- `w_i` también se multiplica por un factor de **deduplicación** (repeticiones del mismo texto pesan menos) y por un **tope por usuario**.
- **Intensidad:** ritmo de mensajes frente a una línea base móvil del propio canal (media móvil de varios minutos), normalizada a 0..1.
- **Tendencia:** comparación del `score` dominante actual contra el de hace unos segundos (`rising`, `steady`, `falling`).

## Calidad de código
- `gofmt`, `go vet`, `golangci-lint`.
- Pruebas con `go test -race`; benchmarks con `testing.B` en `pipeline/` y `aggregate/`.
- Registro estructurado con `log/slog`.
- `pprof` habilitable por configuración para perfilar en pruebas de rendimiento.

## Preguntas abiertas
- Biblioteca para IRC: ¿propia (el protocolo es simple) o dependencia? Se decide en el cambio de OpenSpec correspondiente.
- ¿Un solo canal por instancia en v1? Recomendado: sí, para simplificar. Ver [[Roadmap]].
