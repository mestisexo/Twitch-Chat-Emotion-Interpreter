---
title: Visión y objetivos
tags: [vision, producto]
status: borrador
updated: 2026-10-09
---

# Visión y objetivos

Volver a [[00-Home]].

## El problema
En chats de Twitch con mucha actividad, el streamer no puede leer los mensajes. Pierde la señal más valiosa: **cómo se siente el chat ahora mismo** (hype, risa, enojo, tristeza, sorpresa). Hoy tiene que adivinar por fragmentos que alcanza a ver.

## La propuesta
Un termómetro en tiempo real que resume la emoción dominante del chat, su distribución, su intensidad y su tendencia, visible como overlay en OBS y como panel para el streamer.

## Objetivo general
Dar al streamer una percepción **fiel y rápida** del estado emocional de su chat, sin que tenga que leerlo.

## Objetivos específicos
1. **Rápido:** que un cambio real en el ánimo del chat se refleje en el termómetro en pocos segundos. Detalle en [[07-Latencia-y-Rendimiento]].
2. **Fiel:** que la emoción mostrada coincida razonablemente con lo que un humano percibiría (se mide en [[05-Clasificador-de-Emociones]]).
3. **Estable:** que la aguja no tiemble con cada mensaje; suavizado controlado.
4. **Fácil de instalar:** `docker compose up` y configurar un canal. Ver [[09-Despliegue-Docker]].
5. **Abierto:** cualquier persona puede entender, ejecutar y contribuir. Ver [[14-Onboarding-y-Contribuir]].
6. **Respetuoso con los datos:** no guardar mensajes de personas por defecto. Ver [[13-Seguridad-y-Privacidad]].

## Usuarios
- **Streamer** (principal): quiere un vistazo, no un análisis.
- **Moderadores** (secundario): pueden usar el panel para detectar cambios de ánimo.
- **Contribuidores y desarrolladores:** quieren extender emociones, modelos o visualizaciones.

## Fuera de alcance (v1)
- Moderación automática o detección de toxicidad como producto. La emoción **no** es lo mismo que toxicidad.
- Análisis individual de usuarios o perfiles de personas.
- Histórico de largo plazo y analítica post-stream (se evalúa en [[08-Datos-y-Almacenamiento]]).
- Plataformas distintas de Twitch.

## Principios de diseño
- **Latencia primero, pero con sentido estadístico:** el termómetro es una estimación, no un registro. Muestrear bajo carga es aceptable.
- **Desacoplar:** lo que se clasifica, lo que se agrega y lo que se pinta son etapas independientes.
- **Degradar con gracia:** si falla una pieza, el sistema sigue dando una lectura aproximada.
- **Contratos explícitos** entre componentes, para poder cambiar la tecnología de uno sin tocar los otros.
- **Todo medible:** si no se mide, no se optimiza.

## Taxonomía preliminar de emociones
Punto de partida; se valida con datos reales del chat (ver [[Decisiones-Abiertas]]).

| Etiqueta | Ejemplos de señal en el chat |
|---|---|
| `joy` (diversión) | `KEKW`, `LUL`, "jajaja" |
| `hype` (emoción/euforia) | `PogChamp`, mayúsculas, "VAMOOOS" |
| `surprise` | "QUÉ", "no way", `OMEGALUL` según contexto |
| `anger` | insultos, "qué asco", críticas fuertes |
| `sadness` | `F`, "qué pena", `Sadge` |
| `fear` | "cuidado", "no no no" en situaciones tensas |
| `neutral` | conversación sin carga emocional |

Los emotes y el argot propio de Twitch pesan mucho y los modelos genéricos los entienden mal: es uno de los retos centrales.

## Definición de "v1 lista"
- [ ] Se conecta a un canal de Twitch y recibe mensajes.
- [ ] Clasifica emoción con una calidad mínima acordada en [[05-Clasificador-de-Emociones]].
- [ ] Overlay funcionando en OBS con los datos en vivo.
- [ ] Cumple el presupuesto de latencia de [[07-Latencia-y-Rendimiento]] en los escenarios definidos.
- [ ] Imágenes Docker publicadas y guía de instalación probada por alguien ajeno al proyecto.
- [ ] CI verde y documentación de contribución completa.
