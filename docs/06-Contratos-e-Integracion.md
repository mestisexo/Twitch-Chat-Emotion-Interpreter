---
title: Contratos e integración
tags: [contratos, api, grpc, sse]
status: borrador
updated: 2026-10-09
---

# Contratos e integración

Contexto: [[Arquitectura]]. Decisiones: [[0002-sse-para-tiempo-real]], [[0003-grpc-para-clasificador]].

Los contratos son **la fuente de verdad** entre componentes. Cualquier cambio de contrato pasa por una propuesta de OpenSpec ([[15-Flujo-con-OpenSpec]]).

## 1. Backend ↔ Clasificador (gRPC)
Archivo: `proto/emotion/v1/classifier.proto`

```proto
syntax = "proto3";
package emotion.v1;

service EmotionClassifier {
  rpc ClassifyBatch(ClassifyBatchRequest) returns (ClassifyBatchResponse);
}

message ClassifyBatchRequest {
  repeated Message messages = 1;
}

message Message {
  string id = 1;          // identificador del mensaje dentro del lote
  string text = 2;        // texto ya normalizado
  string lang_hint = 3;   // opcional
}

message ClassifyBatchResponse {
  repeated Prediction predictions = 1;
  string model_version = 2;
  string taxonomy_version = 3;
}

message Prediction {
  string id = 1;
  repeated EmotionScore scores = 2;  // suman ~1.0
}

message EmotionScore {
  string label = 1;   // "joy", "hype", ...
  float score = 2;
}
```

Notas:
- Las etiquetas son `string` + `taxonomy_version` para poder evolucionar la taxonomía sin romper compatibilidad.
- Salud: servicio estándar `grpc.health.v1.Health`.
- **Compatibilidad:** se valida en CI con `buf lint` y `buf breaking` ([[CI CD GitHub]]).
- Código cliente/servidor se genera; no se edita a mano.

## 2. Backend → Frontend (SSE)
Endpoint: `GET /api/v1/stream`

Evento `thermometer` (JSON), emitido cada tick:

```json
{
  "channel": "nombre_del_canal",
  "ts": "2026-10-09T18:30:00.500Z",
  "window_s": 10,
  "msgs_per_s": 42.3,
  "intensity": 0.73,
  "dominant": "hype",
  "distribution": { "hype": 0.41, "joy": 0.30, "surprise": 0.12, "neutral": 0.10, "anger": 0.04, "sadness": 0.02, "fear": 0.01 },
  "trend": "rising",
  "degraded": false,
  "model_version": "x.y.z"
}
```

Reglas:
- Se emite un evento **por tick**, aunque no haya mensajes nuevos (con `msgs_per_s` bajo), así el cliente distingue "chat tranquilo" de "sin conexión".
- Un evento `ping` periódico mantiene viva la conexión a través de proxies.
- Campos nuevos se añaden como opcionales; no se renombran ni eliminan dentro de `v1`.
- Un cambio incompatible crea `/api/v2/`.

## 3. Endpoints operativos
| Ruta | Uso |
|---|---|
| `GET /healthz` | El proceso está vivo |
| `GET /readyz` | Listo (conectado a Twitch y clasificador disponible o en modo degradado declarado) |
| `GET /metrics` | Métricas Prometheus ([[12-Observabilidad]]) |

## 4. Configuración del overlay
Parámetros por URL en el frontend (ver [[04-Frontend-Vite]]). El backend no necesita saber de ellos.

## Versionado
- Contratos: versión en la ruta/paquete (`v1`).
- Releases del proyecto: SemVer. Un cambio incompatible de contrato implica versión mayor.
