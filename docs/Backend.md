---
title: Backend
tags:
  - backend
  - go
status: borrador
updated: 2026-10-09
---

# Backend

Contexto: [[Arquitectura]].


## 💎 Estructura


```
backend/
├── cmd/server/          # main, wiring, arranque y apagado ordenado
├── internal/
│   ├── auth/            # flujo OAuth, state, callback, cookie de sesión   
│   ├── tokens/          # store cifrado, refresh, validate, revocación     
│   ├── session/         # Manager: una goroutine por usuario conectado     
│   ├── twitch/          # cliente EventSub WS, reconexión, suscripciones   
│   ├── pipeline/        # normalización, fast path
│   ├── classify/        # cliente gRPC + fallback
│   ├── aggregate/       # ventana deslizante; una instancia por canal     
│   ├── api/             # SSE por usuario, health, readiness
│   ├── config/          # client_id, secret, redirect_uri, claves               
├── gen/                 # código generado desde proto/
└── Dockerfile
```


## ✴️ Autenticación

- Autenticar al streamer con OAuth y custodiar sus credenciales. 
- Mantener una sesión EventSub con su chat.
- Renovar, validar y revocar tokens sin intervención del usuario. 
- Usar SQLite para almacenar los datos de sesión y refresco de cada usuario.

| Endpoint                    | Para qué sirve                                       | Protocolo |
| --------------------------- | ---------------------------------------------------- | --------- |
| `GET /auth/twitch/login`    | Inicia el flujo OAuth                                | http      |
| `GET /auth/twitch/callback` | Recibe el `code`, intercambia tokens y crea sesión   | http      |
| `POST /auth/logout`         | Cierra sesión, revoca el token y cierra el WebSocket | http      |
| `GET /api/me`               | Devuelve el usuario y el estado de su conexión       | http      |


## ✉️ Mensajes

-  Recibir
-  Normalizar 


## 🔷 Clasificación de sentimientos

-  Resolver por fast path lo que se pueda, y enviar el resto al clasificador.
-  Exponer la clasificación de sentimiento por SSE. 

| Endpoint          | Para qué sirve                             | Protoclo |
| ----------------- | ------------------------------------------ | -------- |
| `GET /api/stream` | Notifica cada cambio de estado al frontend | SSE      |


## 🩻 Operación

-  Exponer end points operativos del estado del sistema.

| Endpoint         | Para qué sirve                     | Protocolo |
| ---------------- | ---------------------------------- | --------- |
| `GET /health`    | Comprueba si el servicio está vivo | http      |
| `GET /api/state` | Devuelve las emociones calculadas  | http      |


## 🪤 Seguridad

* Tokens cifrados en reposo (AES-GCM) y nunca en logs; enmascarar campos sensibles en `slog`.
* Cookie de sesión `HttpOnly`, `Secure`, `SameSite=Lax`.
* Pedir solo el scope de lectura (`user:read:chat`) y mostrárselo al usuario antes de autorizar.
* El archivo SQLite vive en un volumen con permisos restringidos (`0600`) y, al tener los tokens cifrados, no expone credenciales utilizables por sí solo. La clave de cifrado vive fuera del volumen, en una variable de entorno.


## 🪖 Calidad de código

- `gofmt`, `go vet`, `golangci-lint`.
- Pruebas con `go test -race`; benchmarks con `testing.B` en `pipeline/` y `aggregate/`.
- Registro estructurado con `log/slog`.
- `pprof` habilitable por configuración para perfilar en pruebas de rendimiento.
- Pruebas de `tokens` y `session` con un servidor Twitch simulado (refresh, 401, revocación, `session_reconnect`). 


## 🎸 Ciclo de vida de las credenciales

 - Refresco proactivo antes de `expires_at`; ante un 401, renovar y reintentar una vez.
 - Persistir siempre el refresh token que devuelva Twitch.
 - Si el refresco falla, la sesión pasa a `reauth_required`: se cierra su WebSocket y se avisa al frontend por SSE.
 - Atender `session_reconnect` y el keepalive timeout. Las suscripciones se recrean solo al abrir una sesión WS nueva.
- Validación horaria con `/validate`.
- Suscripción a `user.authorization.revoke` para borrar tokens al instante.