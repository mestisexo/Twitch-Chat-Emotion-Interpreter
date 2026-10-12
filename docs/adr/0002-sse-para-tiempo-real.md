---
title: "ADR-0002 SSE para tiempo real"
tags: [adr]
status: propuesta
updated: 2026-10-09
---

# ADR-0002: SSE para enviar el estado al frontend

- **Estado:** propuesta
- **Fecha:** 2026-10-09

## Contexto
El frontend solo necesita **recibir** un estado cada ~500 ms. OBS usa un navegador embebido para el overlay. Contratos en [[06-Contratos-e-Integracion]].

## Opciones
1. **SSE (Server-Sent Events):** unidireccional sobre HTTP, reconexión automática en el navegador, simple de operar detrás de proxies.
2. **WebSocket:** bidireccional, más flexible, más piezas (heartbeats, reconexión manual, manejo de estado).
3. **Polling:** simple pero con mayor latencia o más tráfico.

## Decisión
SSE. La comunicación cliente→servidor, si hiciera falta (configuración), se hace con peticiones HTTP normales.

## Consecuencias
- (+) Menos código y menos modos de fallo.
- (+) Reconexión nativa en `EventSource`.
- (−) Solo texto y un sentido; si en el futuro se requiere interacción en tiempo real bidireccional, se reevalúa WebSocket.
- Límite de conexiones por origen en HTTP/1.1; con HTTP/2 deja de ser un problema práctico.

## Revisión
Si aparece un caso de uso bidireccional (por ejemplo controles remotos en vivo).
