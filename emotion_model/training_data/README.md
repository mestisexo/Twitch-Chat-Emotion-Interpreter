# Datos para ajustar el clasificador

Completa `messages.csv` con una frase por fila. El archivo debe guardarse como CSV UTF-8 y conservar exactamente estos encabezados:

```csv
text,label
```

Ejemplo de formato (no incluido en `messages.csv` para que no se use accidentalmente como dato de entrenamiento):

```csv
text,label
"¡Qué buena partida, estoy feliz!",joy
"No puedo creer que hayamos perdido otra vez",sadness
"Qué rabia, eso fue injusto",anger
"¿En serio pasó eso?",surprise
"Qué asco, fue desagradable",disgust
"Me da miedo que vuelva a ocurrir",fear
"Hola, ¿cómo están?",others
```

## Etiquetas permitidas

Usa solo una de estas etiquetas, escrita exactamente así:

- `joy`: alegría, entusiasmo o diversión.
- `sadness`: tristeza, decepción o pena.
- `anger`: enojo, frustración o indignación.
- `surprise`: sorpresa o asombro.
- `disgust`: rechazo o asco.
- `fear`: miedo o preocupación.
- `others`: mensaje neutral, saludo, información sin emoción clara o texto imposible de clasificar.

Asigna la emoción predominante del texto. Si no hay una emoción clara, usa `others`. No agregues nombres de usuario, columnas, encabezados adicionales ni etiquetas nuevas. Si una frase contiene comas, déjala entre comillas dobles como en el ejemplo. Evita saltos de línea dentro de una frase.

Las frases se escriben en español; las etiquetas van en inglés porque así están definidas en el modelo y en el servicio. Por ejemplo, `"Me vale verga",others` puede etiquetarse como neutral/indiferente si ese es el sentido en el contexto del chat. Si se usa para expresar enojo u otra emoción clara, etiquétala con esa emoción en vez de `others`.

Incluye ejemplos variados y representativos del chat, no solo casos obvios. No uses automáticamente como etiqueta la predicción del modelo actual; revisa y etiqueta el sentido del mensaje. Antes de guardar mensajes reales, elimina nombres, enlaces u otros datos personales que no sean necesarios y confirma que tienes permiso para usar el texto.

Este archivo es solo la plantilla de datos. Todavía no ejecuta entrenamiento; primero reuniremos y revisaremos los ejemplos y después se preparará el script de ajuste y evaluación.
