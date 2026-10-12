---
title: "ADR-0005 No persistir mensajes crudos"
tags: [adr]
status: propuesta
updated: 2026-10-09
---

# ADR-0005: No guardar el texto de los mensajes

- **Estado:** propuesta
- **Fecha:** 2026-10-09

## Contexto
Los mensajes del chat son datos de personas. El objetivo del termómetro solo requiere agregados. Ver [[08-Datos-y-Almacenamiento]] y [[13-Seguridad-y-Privacidad]].

## Opciones
1. **No persistir texto:** procesar en memoria y descartar.
2. **Persistir todo:** habilita análisis futuro y entrenamiento, pero implica obligaciones de privacidad y almacenamiento.

## Decisión
No persistir texto de mensajes ni identificadores de usuario. Si se guarda algo, serán agregados. Los datasets de entrenamiento se gestionan aparte, con consentimiento y anonimización.

## Consecuencias
- (+) Menor riesgo legal y de privacidad; menos infraestructura.
- (+) Más fácil de explicar y de confiar para un proyecto open source.
- (−) No se pueden reanalizar mensajes pasados; los datos de evaluación requieren un proceso deliberado.

## Revisión
Si surge un caso de uso que justifique guardar más, requiere nueva propuesta con evaluación de privacidad.
