---
title: Frontend
tags:
  - frontend
  - vite
status: borrador
updated: 2026-10-09
---

# Frontend

Contexto: [[Arquitectura]].


## 💎 Estructura

```
frontend/
├── src/
│   ├── overlay/
│   ├── dashboard/
│   ├── shared/        # cliente SSE y tipos del contrato
│   └── main.ts
├── tests/
├── vite.config.ts
└── package.json
```


## 🎷 Overlay

- Muestra un polígono generado por pixijs que cambia de color de acuerdo al veredicto de emoción proveniente del backend.
- Configuración por parámetros de URL para que el streamer solo pegue un enlace en OBS, por ejemplo: `/overlay?theme=dark&show=distribution&scale=1`

## 📚 Dashboard

- Estado de la conexión (`connecting`, `active`, `reauth_required`), que viene de `/api/me` del backend


## 🧭 Conexión en tiempo real

* `EventSource` contra `/api/v1/stream` (SSE), con reconexión automática del navegador y reintento manual con espera creciente (1 s hasta 30 s, con variación aleatoria) cuando el servidor cierra la conexión con un error HTTP.
* El backend emite un evento `heartbeat` cada 5 s aunque el estado no cambie. Si pasan más de N segundos sin recibir ningún evento (por defecto 10 s en el dashboard y 6 s en el overlay), la UI muestra "sin señal" en vez de congelar el último valor como si fuera actual. Al llegar el siguiente evento vuelve al estado normal.
* Si `degraded: true` (el clasificador no responde y solo actúa el fast path), mostrar un indicador discreto en el dashboard. En el overlay no se muestra, porque es un dato interno que no le interesa a la audiencia.
* Si llega `session: reauth_required`, la UI deja de reconectar y pide al streamer iniciar sesión de nuevo, con un mensaje distinto al de "sin señal".
* El overlay se autentica con una clave de solo lectura en la URL (`?key=`), y el dashboard con la cookie de sesión.

## 🚤 Animación

- Colores por emoción (etiqueta) para accesibilidad de daltonismo.
- Respetar `prefers-reduced-motion`.
- Texto con buen contraste sobre fondos de juego variados.


## Criterios de rendimiento

- Peso inicial del overlay: objetivo de unos pocos cientos de KB (medir y fijar un límite en CI).
- Sin trabajo continuo innecesario cuando la pestaña no es visible.
