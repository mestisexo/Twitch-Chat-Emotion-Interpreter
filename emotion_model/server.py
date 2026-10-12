import json
import logging
import os
from http.server import BaseHTTPRequestHandler, HTTPServer

import torch
from transformers import AutoModelForSequenceClassification, AutoTokenizer

MODEL_ID = os.getenv("EMOTION_MODEL_ID", "pysentimiento/robertuito-emotion-analysis")
HOST = os.getenv("EMOTION_HOST", "127.0.0.1")
PORT = int(os.getenv("EMOTION_PORT", "8090"))
MAX_BATCH_SIZE = 16
MAX_BODY_SIZE = 128 * 1024
MAX_LENGTH = 128

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
torch.set_num_threads(max(1, int(os.getenv("EMOTION_CPU_THREADS", "2"))))
torch.set_num_interop_threads(1)

logging.info("Cargando clasificador %s; la primera carga descarga los pesos del modelo", MODEL_ID)
tokenizer = AutoTokenizer.from_pretrained(MODEL_ID)
model = AutoModelForSequenceClassification.from_pretrained(MODEL_ID)
model.eval()
if os.getenv("EMOTION_DYNAMIC_INT8", "1").lower() not in {"0", "false", "no"}:
    model = torch.ao.quantization.quantize_dynamic(
        model,
        {torch.nn.Linear},
        dtype=torch.qint8,
    ).eval()
    logging.info("Cuantización dinámica INT8 activada para las capas lineales")
id_to_label = {int(index): label.lower() for index, label in model.config.id2label.items()}
required_labels = {"others", "joy", "sadness", "anger", "surprise", "disgust", "fear"}
if set(id_to_label.values()) != required_labels:
    raise RuntimeError(f"Etiquetas inesperadas en {MODEL_ID}: {id_to_label}")
logging.info("Modelo listo en http://%s:%d", HOST, PORT)


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/health":
            self.send_error(404)
            return
        self.send_json(200, {"status": "ready", "model": MODEL_ID})

    def do_POST(self):
        if self.path != "/classify":
            self.send_error(404)
            return
        content_length = self.headers.get("Content-Length")
        if content_length is None:
            self.send_error(411, "Content-Length requerido")
            return
        try:
            body_size = int(content_length)
        except ValueError:
            self.send_error(400, "Content-Length inválido")
            return
        if body_size < 1 or body_size > MAX_BODY_SIZE:
            self.send_error(413, "Tamaño de lote inválido")
            return
        try:
            payload = json.loads(self.rfile.read(body_size))
        except (json.JSONDecodeError, UnicodeDecodeError):
            self.send_error(400, "JSON inválido")
            return
        messages = payload.get("messages") if isinstance(payload, dict) else None
        if (
            not isinstance(messages, list)
            or not messages
            or len(messages) > MAX_BATCH_SIZE
            or any(not isinstance(message, str) or not message.strip() for message in messages)
        ):
            self.send_error(400, "Se requiere una lista de 1 a 16 mensajes no vacíos")
            return

        try:
            tokens = tokenizer(
                messages,
                padding=True,
                truncation=True,
                max_length=MAX_LENGTH,
                return_tensors="pt",
            )
            with torch.inference_mode():
                logits = model(**tokens).logits
                probabilities = torch.softmax(logits, dim=-1).tolist()
            predictions = [
                {id_to_label[index]: float(score) for index, score in enumerate(row)}
                for row in probabilities
            ]
        except (RuntimeError, ValueError, KeyError) as error:
            logging.exception("Falló la inferencia emocional")
            self.send_json(500, {"error": f"Falló la inferencia: {type(error).__name__}"})
            return
        self.send_json(200, {"predictions": predictions})

    def send_json(self, status, payload):
        encoded = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def log_message(self, format_string, *args):
        logging.info("HTTP %s %s", self.client_address[0], format_string % args)


HTTPServer((HOST, PORT), Handler).serve_forever()
