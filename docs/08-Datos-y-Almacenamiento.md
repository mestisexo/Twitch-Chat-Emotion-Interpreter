---
title: Datos y almacenamiento
tags: [datos, almacenamiento, privacidad]
status: borrador
updated: 2026-10-09
---

# Datos y almacenamiento

Contexto: [[Arquitectura]]. Decisión: [[0005-no-persistir-mensajes-crudos]]. Privacidad: [[13-Seguridad-y-Privacidad]].

## v1: sin base de datos
Todo el estado vive **en memoria**:
- Buffer circular con las predicciones de la ventana activa.
- Caché LRU de textos ya clasificados.
- Línea base móvil del ritmo de mensajes.
- Un historial corto de estados para el panel (por ejemplo la última hora), también en memoria.

Si el proceso se reinicia, el termómetro se reconstruye en pocos segundos. No hace falta persistencia para cumplir el objetivo de [[Meta]].

## Qué NO se guarda
- **El texto de los mensajes.** Se procesa y se descarta.
- Datos que identifiquen a personas del chat (nombres de usuario) más allá del tiempo necesario para el tope por usuario en la ventana en memoria.

## Qué podría guardarse en el futuro (opcional)
Solo **agregados**, nunca texto crudo:

| Dato | Para qué |
|---|---|
| Serie temporal del estado emitido (cada pocos segundos) | Reportes post-stream, "momentos de hype" |
| Métricas operativas | Ya cubiertas por Prometheus ([[12-Observabilidad]]) |

Opciones si se necesita: SQLite (autoalojado simple), o Postgres/TimescaleDB (más volumen o varios usuarios). Se decide cuando exista el requisito real; se propone entonces como un cambio de OpenSpec.

## Datos para entrenar y evaluar el modelo
Son una cosa **separada** del sistema en ejecución.
- Se recopilan solo con consentimiento y anonimización (sin nombres de usuario ni identificadores).
- Se almacenan fuera de la imagen Docker y fuera del repositorio principal si son sensibles o grandes; en el repo solo van muestras sintéticas o publicables.
- Se documenta su origen, licencia y proceso de etiquetado en `classifier/eval/`. Ver [[05-Clasificador-de-Emociones]].

## Retención
Para los agregados opcionales: política explícita y configurable (por defecto, corta). Documentarla en el README del despliegue.
