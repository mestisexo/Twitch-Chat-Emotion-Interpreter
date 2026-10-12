---
title: "ADR-0003 gRPC para el clasificador"
tags: [adr]
status: propuesta
updated: 2026-10-09
---

# ADR-0003: gRPC por lotes entre backend y clasificador

- **Estado:** propuesta
- **Fecha:** 2026-10-09

## Contexto
El backend (Go) debe enviar lotes de mensajes al clasificador (Python en v1) con baja sobrecarga y un contrato claro que permita reemplazar el clasificador. Ver [[05-Clasificador-de-Emociones]] y [[07-Latencia-y-Rendimiento]].

## Opciones
1. **gRPC + Protobuf:** contrato tipado, generación de código en ambos lenguajes, conexión persistente, validación de compatibilidad con `buf`.
2. **HTTP/JSON:** más simple de depurar, sin esquema estricto, más sobrecarga de serialización.
3. **Cola de mensajes (Redis, NATS):** desacopla, pero añade infraestructura innecesaria para v1.

## Decisión
gRPC con un método unario por lote (`ClassifyBatch`). La inferencia domina el tiempo, pero el contrato fuerte es el motivo principal.

## Consecuencias
- (+) Contrato versionado y verificado en CI.
- (+) Cambiar el clasificador de Python a otra tecnología no toca el backend.
- (−) Depuración menos directa que con JSON (mitigable con herramientas como `grpcurl`).
- (−) Dependencia de generación de código en el flujo de build.

## Revisión
Si las mediciones muestran que el costo del transporte es relevante frente a la inferencia.
