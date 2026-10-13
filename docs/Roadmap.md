---
title: Roadmap
tags: [roadmap, planificacion]
status: borrador
updated: 2026-10-09
---

# Roadmap

Cada paso se ejecuta como un cambio de OpenSpec [[15-Flujo-con-OpenSpec]] (`id` entre comillas invertidas).

## M0: Fundaciones
- [ ] `bootstrap-monorepo`: estructura, Makefile, `.env.example`, README
- [ ] `add-ci-pipeline`: workflows básicos por componente ([[CI CD GitHub]])
- [ ] `add-proto-contracts`: `proto/` + `buf` ([[06-Contratos-e-Integracion]])
- [ ] Archivos de comunidad y licencia ([[14-Onboarding-y-Contribuir]])

## M1: Tubería mínima de extremo a extremo (sin IA real)
- [ ] `add-chat-source-simulator`: simulador de chat ([[11-Pruebas]])
- [ ] `add-pipeline-preprocessing`: normalización, filtros, dedupe
- [ ] `add-emotion-aggregation`: ventana, intensidad, tendencia (con clasificador falso)
- [ ] `add-realtime-stream-api`: SSE con tick
- [ ] `add-thermometer-overlay`: overlay básico
- **Resultado:** ver el termómetro moverse con datos simulados.

## M2: Clasificador real
- [ ] `add-emotion-classifier-service`: servicio Python con gRPC ([[05-Clasificador-de-Emociones]])
- [ ] `add-micro-batching`: batcher y cliente gRPC con timeouts y modo degradado
- [ ] `add-fastpath-emotes`: emotes, léxico y caché
- [ ] `add-golden-set-eval`: conjunto de oro y métricas de calidad

## M3: Twitch real
- [ ] `add-twitch-chat-ingest`: fuente real de mensajes ([[03-Backend-Go]])
- [ ] Reconexión y pruebas de resiliencia

## M4: Rendimiento y operación
- [ ] `add-load-testing`: escenarios y `loadtest/` ([[07-Latencia-y-Rendimiento]])
- [ ] `add-observability`: métricas y tablero ([[12-Observabilidad]])
- [ ] Primera ronda de optimización guiada por mediciones
- [ ] Documentar recursos mínimos medidos

## M5: Producción y comunidad
- [ ] `add-docker-deployment`: imágenes y Compose ([[09-Despliegue-Docker]])
- [ ] `add-release-pipeline`: publicación a GHCR, releases
- [ ] Panel del streamer
- [ ] Revisión de seguridad y privacidad ([[13-Seguridad-y-Privacidad]])
- [ ] **Lanzamiento v1.0**

## Más adelante (no comprometido)
- EventSub y varios canales
- Persistencia opcional de agregados ([[08-Datos-y-Almacenamiento]])
- Modelo propio afinado con datos de chat
- Resumen del chat con un LLM (sería un agente real)
- Otras plataformas
