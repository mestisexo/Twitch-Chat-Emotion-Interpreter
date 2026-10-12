---
title: Pruebas
tags: [pruebas, calidad]
status: borrador
updated: 2026-10-09
---

# Pruebas

Contexto: [[10-CI-CD-GitHub]], [[07-Latencia-y-Rendimiento]].

## Estrategia por niveles
| Nivel | Qué verifica | Herramientas |
|---|---|---|
| Unitarias | Lógica pura: normalización, filtros, agregación, caché | `go test -race`, `pytest`, Vitest |
| Contrato | Que backend, clasificador y frontend respeten los contratos | `buf breaking`, validación de esquemas JSON |
| Integración | Backend + clasificador reales, Twitch simulado | Docker Compose + simulador de chat |
| E2E | Del mensaje al píxel: overlay mostrando el estado esperado | Playwright |
| Calidad del modelo | Acierto sobre el conjunto de oro y por ventana | Scripts en `classifier/eval/` |
| Rendimiento | Latencia, throughput, regresiones | Benchmarks Go, k6, simulador |
| Resiliencia | Fallos de clasificador, caída de conexión a Twitch, reconexión | Pruebas de integración con inyección de fallos |

## El simulador de chat
Pieza central para probar sin depender de Twitch real:
- Implementa la interfaz `ChatSource` del backend ([[03-Backend-Go]]).
- Modos: **sintético** (genera mensajes con ritmo y emociones configurables), **replay** (reproduce grabaciones anonimizadas) y **guion** (secuencias deterministas para pruebas, por ejemplo "neutral 20 s, luego hype").
- Permite verificar de forma determinista que el termómetro reacciona como se espera.

## Qué probar sí o sí (primeros casos)
- Un mensaje con solo emotes se resuelve por fast path y **no** llama al clasificador.
- Un texto repetido 100 veces pesa menos que 100 textos distintos.
- Un usuario no puede dominar la lectura aunque envíe muchos mensajes.
- Bajo saturación se descartan mensajes antiguos y se cuenta el descarte.
- Si el clasificador no responde, se emite `degraded: true` y el sistema sigue.
- El tick se emite aunque no lleguen mensajes.
- El frontend muestra "sin señal" si deja de recibir eventos.
- Reconexión automática de SSE y de la fuente de chat.

## Datos de prueba
- Fixtures sintéticos y pequeños en el repo.
- **Nunca** mensajes reales identificables en el repositorio ([[13-Seguridad-y-Privacidad]]).

## Criterios en CI
- Pruebas unitarias, de contrato e integración bloquean el merge.
- Rendimiento: micro-benchmarks con umbral de regresión; cargas completas en ejecución programada.
- Calidad del modelo: no se aceptan cambios que reduzcan la métrica acordada.

## Cobertura
Se usa como señal, no como objetivo único. Prioridad: cubrir bien el camino crítico (pipeline, agregación) antes que perseguir un porcentaje global.
