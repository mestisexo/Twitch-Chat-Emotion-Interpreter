---
title: Frontend (Vite)
tags: [frontend, vite]
status: borrador
updated: 2026-10-09
---

# Frontend (Vite)

Contexto: [[02-Arquitectura]]. Eventos que consume: [[06-Contratos-e-Integracion]].

## Vistas
1. **Overlay (OBS):** fuente de navegador con fondo transparente. Muestra el termómetro: emoción dominante, aguja o barra de intensidad y, opcionalmente, la distribución. Debe ser muy ligero.
2. **Panel del streamer / mods:** vista más completa con historial corto de la sesión (en memoria), indicador de estado (conectado, degradado) y controles.

Configuración por parámetros de URL para que el streamer solo pegue un enlace en OBS, por ejemplo:
`/overlay?theme=dark&show=distribution&scale=1`

## Stack
- Vite + TypeScript.
- Framework: **pendiente** (ver [[Decisiones-Abiertas]]). Criterio: el overlay debe pesar poco y arrancar rápido; para una vista tan pequeña, una librería ligera o TypeScript sin framework es razonable. El panel puede justificar más estructura.
- Tests con Vitest; E2E con Playwright (ver [[11-Pruebas]]).

## Conexión en tiempo real
- `EventSource` contra `/api/v1/stream` (SSE), con reconexión automática. Razón: [[0002-sse-para-tiempo-real]].
- Si el estado deja de llegar por más de N segundos, la UI muestra "sin señal" en vez de congelar el último valor como si fuera actual.
- Si `degraded: true`, mostrar un indicador discreto.

## Animación y legibilidad
- Interpolar entre estados recibidos (la aguja no salta; se desplaza suavemente) usando `requestAnimationFrame`.
- Colores por emoción **con alternativa no basada solo en color** (icono o etiqueta) para accesibilidad (daltonismo).
- Respetar `prefers-reduced-motion`.
- Texto con buen contraste sobre fondos de juego variados.

## Estructura propuesta
```
frontend/
├── src/
│   ├── overlay/
│   ├── dashboard/
│   ├── shared/        # cliente SSE, tipos del contrato, utilidades
│   └── main.ts
├── tests/
├── vite.config.ts
└── package.json
```

## Tipos del contrato
Los tipos TypeScript del evento `thermometer` se mantienen alineados con el esquema descrito en [[06-Contratos-e-Integracion]]. Idealmente se generan o se validan automáticamente para detectar divergencias en CI.

## Integración con el backend
- **Desarrollo:** `vite dev` con proxy a `localhost:8080` para `/api`.
- **Producción:** `vite build` y los archivos estáticos se embeben en el binario del backend. Detalle en [[09-Despliegue-Docker]].

## Criterios de rendimiento
- Peso inicial del overlay: objetivo de unos pocos cientos de KB (medir y fijar un límite en CI).
- Sin trabajo continuo innecesario cuando la pestaña no es visible.
