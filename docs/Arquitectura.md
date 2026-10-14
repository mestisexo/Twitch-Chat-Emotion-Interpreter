---
title: Arquitectura
tags: [arquitectura]
status: borrador
updated: 2026-10-09
---

# Arquitectura

Volver a [[Home]].


## 💎 Estructura


```
chat-emotion-meter/
├── backend/          # Go
├── frontend/         # Vite
├── classifier/       # Python
├── proto/            # Contratos gRPC (fuente de verdad)
├── deploy/           # docker-compose y configuración
├── loadtest/         # Escenarios de carga y simulador de chat
├── docs/             # Este vault de Obsidian
├── openspec/         # Especificaciones y cambios
├── .github/          # Workflows, plantillas de issue y PR
├── AGENTS.md
├── CONTRIBUTING.md
└── README.md
```


## 🫳🏽 Componentes

| Componente | Tecnología | Responsabilidad                                                               | Nota                             |
| ---------- | ---------- | ----------------------------------------------------------------------------- | -------------------------------- |
| Backend    | Go         | Autenticación, Ingesta de mensajes, normalización, clasificación y respuesta. | [[Backend]]                      |
| Classifier | Python     | Conversión de mensajes en emoción categorizada                                | [[05-Clasificador-de-Emociones]] |
| Frontend   | Vite       | Panel de Autenticación, Overlay para Obs                                      | [[Frontend]]             |
