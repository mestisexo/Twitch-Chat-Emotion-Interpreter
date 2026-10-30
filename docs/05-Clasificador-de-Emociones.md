---
title: Clasificador de emociones
tags: [ia, clasificador, python]
status: borrador
updated: 2026-10-09
---

# Clasificador de emociones

Contexto: [[Arquitectura]]. Contrato gRPC: [[06-Contratos-e-Integracion]]. Decisión de servicio separado: [[0004-clasificador-servicio-separado]].

## Rol
Recibe lotes de mensajes cortos y devuelve, para cada uno, una probabilidad por emoción. No guarda estado entre llamadas.

## Particularidades del chat de Twitch
- Mensajes muy cortos, con **emotes** que cargan la emoción (`KEKW`, `Sadge`).
- Argot, errores ortográficos, mayúsculas para énfasis, repeticiones ("jajajaja", "VAMOOOOS").
- **Multilenguaje** y mezcla de idiomas en el mismo chat.
- Muchos mensajes repetidos (copypasta, spam).

Un modelo genérico de sentimiento suele fallar aquí; por eso el diseño combina modelo + reglas.

## Escalera de optimización (de más barato a más costoso)
Cada peldaño reduce cuántos mensajes llegan al siguiente. Se activa solo si las mediciones de [[07-Latencia-y-Rendimiento]] lo justifican.

1. **Fast path por emotes y léxico:** diccionario de emotes y palabras con su emoción. Resuelve una fracción grande sin tocar el modelo. (Este paso vive en el backend.)
2. **Caché de textos exactos** (LRU): el chat repite mucho.
3. **Micro-batching:** un lote aprovecha mucho mejor la CPU/GPU que mensajes sueltos.
4. **Modelo pequeño y cuantizado:** exportar a ONNX y cuantizar (por ejemplo INT8) para bajar latencia en CPU.
5. **Distillation o modelo más liviano** si el modelo base es el cuello de botella.
6. **GPU opcional** (imagen aparte) para despliegues de alto volumen.

## Candidatos de modelo (a evaluar, no decididos)
| Familia | Idea | Ojo con |
|---|---|---|
| Modelos de emociones en español (por ejemplo, de la familia *pysentimiento*) | Buen punto de partida para español | Taxonomía propia (incluye clases como `others`) y dominio distinto al chat de Twitch |
| Modelos multilingües tipo XLM-R, afinados a emociones | Chats multilenguaje | Más pesados |
| Modelo propio afinado con datos de chat | Mejor calidad en argot y emotes | Requiere datos etiquetados |
| LLM pequeño | Comprende contexto | Latencia y costo generalmente altos para este uso |

> [!warning] Licencias
> Revisar la licencia de cada modelo y dataset antes de incluirlos o redistribuirlos en un proyecto open source.

## Evaluación de calidad
- **Conjunto de oro:** mensajes de chat etiquetados por humanos (idealmente 2 o más anotadores, midiendo acuerdo). Se guarda **solo con consentimiento y anonimizado**; ver [[13-Seguridad-y-Privacidad]].
- **Métricas:** F1 macro, matriz de confusión, y calidad **a nivel de ventana** (el termómetro agrega, así que importa que la emoción dominante de una ventana coincida con la percibida por humanos, más que acertar cada mensaje).
- **Regresión:** el conjunto de oro corre en CI cuando cambie el modelo o el preprocesado; no se acepta un cambio que baje la métrica acordada.
- Reporte de versión del modelo: la respuesta gRPC incluye `model_version`.

## Interfaz con el backend
- gRPC, método `ClassifyBatch`, lote de hasta N mensajes. Se detalla en [[06-Contratos-e-Integracion]].
- Servicio de salud estándar de gRPC para readiness.
- Límite de tamaño de lote y de longitud de texto para proteger el servicio.

## Opciones de ejecución
| Opción | Ventaja | Desventaja |
|---|---|---|
| **Servicio Python (v1)** | Ecosistema ML, rapidez de desarrollo | Segundo proceso, imagen más pesada |
| ONNX Runtime dentro del backend Go | Un solo binario, sin salto de red | Preprocesado/tokenización más laborioso en Go |
| Servidor de inferencia dedicado | Escala y optimizaciones | Más complejidad operativa |

v1: servicio Python. El contrato gRPC permite migrar luego sin cambiar el resto del sistema.

## Estructura propuesta
```
classifier/
├── app/
│   ├── server.py        # gRPC
│   ├── model.py         # carga e inferencia
│   └── preprocess.py
├── eval/                # conjunto de oro y scripts de evaluación
├── tests/
├── pyproject.toml
└── Dockerfile
```

## Imagen y reproducibilidad
Los pesos del modelo se descargan en build con **revisión fijada** (hash/tag) y quedan dentro de la imagen, para que "imagen fija" signifique resultados reproducibles. Ver [[09-Despliegue-Docker]].




**Cálculo (incremental):** en vez de recorrer todos los mensajes de la ventana, se mantiene por canal un acumulador ponderado por emoción (`N_e`) y la suma de pesos (`D`). Al llegar un mensaje, ambos se multiplican por el factor de decaimiento `0.5^(Δt / HALF_LIFE)` y luego se suma el mensaje nuevo (`N_e += w·p_e`, `D += w`). El score es `N_e / D`. Es equivalente a la fórmula, con costo constante por mensaje y sin guardar historial. Intensidad y tendencia siguen el mismo principio con medias móviles exponenciales (rápida y lenta).




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
