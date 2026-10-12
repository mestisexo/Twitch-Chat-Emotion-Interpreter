---
title: "ADR-0004 Clasificador como servicio separado"
tags: [adr]
status: propuesta
updated: 2026-10-09
---

# ADR-0004: Clasificador como servicio independiente (Python en v1)

- **Estado:** propuesta
- **Fecha:** 2026-10-09

## Contexto
La primera versión del clasificador está en Python y el backend es Go. Falta definir si integrarlos en un solo proceso o mantenerlos separados.

## Opciones
1. **Servicio Python separado:** aprovecha el ecosistema ML, iteración rápida, escalado independiente.
2. **ONNX Runtime dentro del backend Go:** un solo binario, sin salto de red; la tokenización y el preprocesado en Go son más laboriosos.
3. **Servidor de inferencia dedicado:** potente, pero complejidad operativa alta para v1.

## Decisión
Servicio separado en Python para v1, detrás del contrato de [[0003-grpc-para-clasificador]]. Se reevalúa la opción 2 con datos de rendimiento.

## Consecuencias
- (+) Rapidez de desarrollo, experimentación de modelos sencilla.
- (+) Cambiar de modelo o lenguaje no afecta al backend.
- (−) Dos imágenes y un salto de red local.
- (−) Imagen del clasificador más pesada ([[09-Despliegue-Docker]]).

## Revisión
Cuando existan mediciones: si el salto de red o el tamaño de la imagen son un problema real, evaluar ONNX en Go.
