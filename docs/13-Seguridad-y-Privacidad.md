---
title: Seguridad y privacidad
tags: [seguridad, privacidad]
status: borrador
updated: 2026-10-09
---

# Seguridad y privacidad

Contexto: [[08-Datos-y-Almacenamiento]], [[0005-no-persistir-mensajes-crudos]].

> [!warning] Esto no es asesoría legal
> Antes de publicar, revisar los términos de Twitch para desarrolladores y la normativa de protección de datos aplicable (por ejemplo GDPR si hay usuarios en la UE).

## Privacidad por diseño
- Los mensajes del chat son datos de personas. **Se procesan en memoria y se descartan.**
- No se muestran mensajes ni nombres de usuario en el overlay: solo agregados.
- Los logs no incluyen el texto de los mensajes en niveles normales.
- Datasets de entrenamiento: solo con consentimiento y anonimizados ([[05-Clasificador-de-Emociones]]).

## Uso responsable
- El termómetro describe **ánimo colectivo**, no juzga a individuos.
- No es una herramienta de moderación ni de perfilado. Documentarlo en el README para evitar usos indebidos.
- La emoción detectada es una estimación y puede equivocarse (ironía, argot, sarcasmo). Mostrarla como indicador, no como verdad.

## Seguridad del servicio
| Tema | Medida |
|---|---|
| Secretos (tokens de Twitch si se usan en v2) | Variables de entorno o secretos del orquestador; nunca en el repo ni en imágenes |
| Acceso al overlay | Opción de token en la URL para que no cualquiera lea el estado de un canal |
| CORS | Lista explícita de orígenes permitidos |
| Entrada al clasificador | Límites de tamaño de lote y de longitud de texto; timeouts |
| Contenedores | Usuario no root, imagen mínima, sistema de archivos de solo lectura |
| Dependencias | Dependabot, escaneo de imágenes, versiones fijadas ([[GitHub]]) |
| Cadena de suministro | Acciones de CI fijadas; revisar licencias de modelos |
| gRPC interno | Red interna de Compose; TLS si el clasificador corre en otra máquina |

## Reporte de vulnerabilidades
`SECURITY.md` con un canal privado (por ejemplo "Security advisories" de GitHub) y plazos de respuesta.

## Licencias
- Código: pendiente ([[Decisiones Abiertas]]).
- Modelos y datasets: verificar compatibilidad con la licencia del proyecto y documentarla.
