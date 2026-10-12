# Overlay de Twitch con PixiJS

## Arranque

Primero inicia el clasificador local de emociones (la primera ejecución descarga el modelo):

```powershell
cd emotion_model
pip install -r requirements.txt
python server.py
```

En otra terminal, inicia el backend con el nombre del canal:

```powershell
go run . -channel nombredelcanal
```

En una tercera terminal, ejecuta Vite:

```powershell
cd frontend
npm install
npm run dev
```

Abre la URL que muestra Vite. El overlay conecta por defecto a `ws://localhost:8080/ws`.
Para usar otra dirección, agrega el parámetro `?ws=ws://servidor:8080/ws` a la URL.

El chat comparte una figura hexagonal cuyo color representa la emoción grupal estimada: alegría (amarillo), sorpresa (cian), tristeza (azul), enojo (naranja), disgusto (verde), miedo (violeta) y neutral/sin clasificar (gris). Los colores se configuran en `json/scores.json`. El halo pulsa según la confianza de la emoción dominante.

El ticker de PixiJS está limitado a 30 FPS para reducir el trabajo de renderizado en OBS; la figura sigue animándose con una fluidez menor que a 60 FPS.

El backend analiza lotes de hasta 16 mensajes. Al completar 16, envía el lote enseguida; si queda incompleto, lo envía al cumplirse 1 segundo desde el primer mensaje. No hay cálculo de tasa de mensajes ni sondeo periódico. Mientras el clasificador procesa o reintenta un lote, los mensajes entrantes se descartan; al terminar, comienza a reunir el siguiente lote. La cola tiene capacidad para 16 y cualquier mensaje que no quepa también se descarta. El backend registra cuántos mensajes se descartaron. Las emociones se promedian solo dentro del lote analizado: ningún lote anterior influye en el siguiente. El resumen grupal actual se conserva únicamente para enviarlo a clientes nuevos. El texto se envía al servicio de modelo local, no a una API externa. Si el clasificador falla, el backend lo registra y reintenta el mismo lote.

El modelo predeterminado es [`pysentimiento/robertuito-emotion-analysis`](https://huggingface.co/pysentimiento/robertuito-emotion-analysis), un clasificador de emociones en español entrenado con texto de redes sociales. Descarga sus pesos en el primer arranque y ejecuta inferencia en CPU; no requiere una API con cuota. Las capas lineales se cuantizan dinámicamente a INT8 al iniciar (`EMOTION_DYNAMIC_INT8=0` desactiva la cuantización). En pruebas locales con lotes de 1, 4 y 16 mensajes, INT8 conservó la emoción dominante en los 16 ejemplos evaluados, redujo alrededor de 75 % la memoria de los pesos lineales y bajó la latencia media de 51.8/64.4/137.4 ms a 22.4/34.4/73.0 ms, respectivamente; la cuantización sí puede cambiar las probabilidades/confianzas. Repite estas mediciones en el hardware de despliegue. La tarjeta del modelo reporta F1 macro de alrededor de 0.56 en emociones: puede equivocarse con ironía, ambigüedades o slang distinto a sus datos. Clasifica tono textual, no el estado real de las personas. La tarjeta del modelo no declara claramente una licencia; revisa sus términos antes de despliegue público o comercial.

El servicio del modelo escucha solo en `127.0.0.1:8090`. Usa un servidor HTTP de un solo hilo porque el backend clasifica secuencialmente; PyTorch limita los hilos intraoperación a 2 e interoperación a 1 por defecto. `EMOTION_HOST`, `EMOTION_PORT`, `EMOTION_MODEL_ID`, `EMOTION_CPU_THREADS` y `EMOTION_DYNAMIC_INT8` permiten ajustar el servicio; `EMOTION_SERVICE_URL` cambia el endpoint que usa Go. Despliega ambos procesos en el mismo servidor y permite que el backend alcance la dirección local configurada.

Los benchmarks del backend están en `hub_benchmark_test.go`; ejecútalos con `go test -run '^$' -bench 'BenchmarkHubOnChat|BenchmarkApplyEmotionResults' -benchmem`. En el Ryzen 7 5700X de desarrollo, poner un mensaje en la cola tardó unos 46 ns y no asignó memoria. Agregar probabilidades para un lote de 16 bajó de unos 5.5 µs/753 B a 4.0 µs/496 B tras evitar mapas intermedios. Estas cifras miden backend, no el costo dominante de inferencia del modelo.

El fondo del canvas está configurado en verde `#00ff00` para pruebas con chroma key en OBS.

## Organización del frontend

- `src/main.js`: inicializa PixiJS y coordina la conexión WebSocket.
- `src/score.js`: valida los colores configurados y selecciona la emoción dominante.
- `src/shapes.js`: define y dibuja la figura y su halo animado.
- `src/figures.js`: muestra y anima la figura con el tono grupal enviado por el backend.
