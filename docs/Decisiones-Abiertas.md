---
title: Decisiones abiertas
tags: [decisiones, pendiente]
status: vivo
updated: 2026-10-09
---

# Decisiones abiertas

Volver a [[00-Home]]. Cuando una decisión se cierre, se convierte en un ADR ([[0000-plantilla-adr]]) y se quita de aquí.

| #   | Tema                                | Opciones                      | Mi recomendación inicial                                                                                | Dónde se discute                        |
| --- | ----------------------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------------- | --------------------------------------- |
| 1   | **Nombre** del proyecto             | Pendiente                     | Elegir antes de crear el repo y las imágenes                                                            |                                         |
| 2   | **Licencia**                        | MIT / Apache-2.0 / AGPL       | Apache-2.0 o MIT para máxima adopción; AGPL si quieres frenar servicios cerrados basados en el proyecto | [[10-CI-CD-GitHub]]                     |
| 3   | **Idioma** de docs y código         | Español / inglés / ambos      | Código y README en inglés para alcanzar más contribuidores; docs internas pueden empezar en español     | [[14-Onboarding-y-Contribuir]]          |
| 4   | **Taxonomía** de emociones          | 7 clases propuestas / otra    | Empezar con la de [[01-Vision-y-Objetivos]] y validar con datos reales                                  | [[05-Clasificador-de-Emociones]]        |
| 5   | **Modelo** del clasificador         | Ver candidatos                | Evaluar 2 o 3 candidatos contra un conjunto de oro pequeño                                              | [[05-Clasificador-de-Emociones]]        |
| 6   | **Framework** del frontend          | Ligero / React / Vue / Svelte | Algo liviano para el overlay                                                                            | [[04-Frontend-Vite]]                    |
| 7   | **Fuente de Twitch** en v1          | IRC anónimo / EventSub        | IRC anónimo para v1, verificando documentación vigente                                                  | [[03-Backend-Go]]                       |
| 8   | **Un canal o varios** por instancia | 1 / N                         | 1 en v1                                                                                                 | [[16-Roadmap]]                          |
| 9   | **Ubicación del clasificador**      | Servicio Python / ONNX en Go  | Servicio Python en v1                                                                                   | [[0004-clasificador-servicio-separado]] |
| 10  | **Parámetros de suavizado**         | Ventana y vida media          | Empezar con 10 s y 4 s y ajustar con usuarios                                                           | [[07-Latencia-y-Rendimiento]]           |
| 11  | **Herramienta de tareas**           | Make / just / task            | Make por simplicidad                                                                                    | [[14-Onboarding-y-Contribuir]]          |
