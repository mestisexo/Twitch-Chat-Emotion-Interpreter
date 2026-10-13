---
title: Definición clara de la meta
tags:
  - definicion
  - objetivo
  - problema
status: borrador
updated: 2026-10-09
---

# Definición clara de la meta

Volver a [[Home]].

## El problema
En los chat de Twitch con mucha actividad, el streamer difícilmente alcanza a leer los mensajes. Pierde la señal más valiosa: **cómo se siente la audiencia ahora mismo**.

## La solución
Un termómetro en tiempo real que resume la emoción dominante del chat, su distribución, su intensidad y su tendencia, visible como overlay para la audiencia y como panel para el streamer.

## Objetivo general
Dar al streamer una percepción **fiel y rápida** del estado emocional de su chat.

## Objetivos específicos
1. Que un cambio real en el ánimo del chat se refleje en el termómetro casi de inmediato. Ver en [[07-Latencia-y-Rendimiento]].
2. Que cada emoción mostrada coincida razonablemente con lo que un humano percibiría (como se logrará en [[05-Clasificador-de-Emociones]]).
3. Que la aguja no tiemble con cada mensaje; su momentum debe ser suavizado y controlado.
4. Que sea fácil de instalar y configurar un canal. Ver [[09-Despliegue-Docker]].
5. Cualquier persona puede entender y contribuir. Ver [[14-Onboarding-y-Contribuir]].
6. No guardar mensajes de personas. Ver [[13-Seguridad-y-Privacidad]].

## Usuarios
- **Streamer** (principal).
- **Moderadores** (secundario). opcionalmente podrían usar el panel para detectar cambios de ánimo del chat.
- **Contribuidores y desarrolladores:** podrían extender emociones, modelos o visualizaciones.

## Fuera de alcance (v1)
- Moderación automática o detección de toxicidad.
- Análisis individual de usuarios o perfiles.
- Histórico de largo plazo y analítica post-stream (se evalúa en [[08-Datos-y-Almacenamiento]]).
- Plataformas distintas de Twitch.

## Principios de diseño
- **Desacoplar:** lo que se clasifica, lo que se agrega y lo que se pinta son etapas independientes.
- **Degradar con gracia:** si falla una pieza, el sistema sigue dando una lectura aproximada.
- **Contratos explícitos** entre componentes, para poder cambiar la tecnología de uno sin tocar los otros.
- **Todo medible:** si no se mide, no se optimiza.

## Taxonomía preliminar de emociones
Punto de partida; se valida con datos reales del chat (ver [[Decisiones Abiertas]]).

| Etiqueta                 | Ejemplos de señal en el chat                |
| ------------------------ | ------------------------------------------- |
| `joy` (diversión)        | `KEKW`, `LUL`, "jajaja", xd, xdd            |
| `hype` (emoción/euforia) | `PogChamp`, mayúsculas, "VAMOOOS", SII      |
| `surprise`               | "QUÉ", "no way", `OMEGALUL`, what           |
| `anger`                  | insultos, "qué asco", críticas fuertes      |
| `sadness`                | `F`, "qué pena", `Sadge`, que triste        |
| `fear`                   | "cuidado", "no no no" en situaciones tensas |

Los emotes y el argot propio de Twitch pesan mucho y los modelos genéricos los entienden mal: es uno de los retos centrales. (pendiente por revisar)

## Definición de "v1 lista"
- [ ] Se conecta a un canal de Twitch y recibe mensajes.
- [ ] Clasifica emoción con una calidad mínima acordada en [[05-Clasificador-de-Emociones]].
- [ ] Overlay funcionando en OBS con los datos en vivo.
- [ ] Cumple el presupuesto de latencia [[07-Latencia-y-Rendimiento]] en los escenarios definidos.
- [ ] Imágenes Docker publicadas y guía de instalación probada por alguien ajeno al proyecto.
- [ ] CI verde y documentación de contribución completa.
