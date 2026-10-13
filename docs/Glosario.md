---
title: Glosario
tags: [glosario]
status: vivo
updated: 2026-10-09
---

# Glosario

| Término | Significado |
|---|---|
| **Termómetro** | Resumen en tiempo real de la emoción del chat (dominante, distribución, intensidad, tendencia) |
| **Emote** | Imagen/código del chat de Twitch (`KEKW`, `Sadge`) que suele cargar emoción |
| **Clasificador** | Servicio que convierte texto en probabilidades por emoción ([[05-Clasificador-de-Emociones]]) |
| **Fast path** | Camino barato que resuelve mensajes sin llamar al modelo (emotes, léxico, caché) |
| **Micro-batch** | Agrupar varios mensajes en un solo llamado al clasificador |
| **Tick** | Intervalo fijo en que el backend emite el estado a los clientes |
| **Ventana** | Periodo reciente de mensajes que se agrega para calcular el estado |
| **Vida media (half-life)** | Tiempo en que el peso de un mensaje cae a la mitad dentro de la ventana |
| **Intensidad** | Qué tan activo está el chat frente a su ritmo habitual |
| **Backpressure** | Mecanismo para que una etapa lenta no acumule trabajo sin límite |
| **Modo degradado** | Funcionamiento sin clasificador, solo con fast path |
| **SSE** | Server-Sent Events: flujo unidireccional del servidor al navegador sobre HTTP ([[0002-sse-para-tiempo-real]]) |
| **gRPC** | Llamadas remotas tipadas sobre HTTP/2 ([[0003-grpc-para-clasificador]]) |
| **ADR** | Registro de decisión de arquitectura ([[0000-plantilla-adr]]) |
| **OpenSpec** | Método de desarrollo guiado por especificaciones ([[15-Flujo-con-OpenSpec]]) |
| **Conjunto de oro** | Mensajes etiquetados por humanos para evaluar el modelo |
| **Overlay** | Vista del termómetro incrustada en el stream vía OBS |
