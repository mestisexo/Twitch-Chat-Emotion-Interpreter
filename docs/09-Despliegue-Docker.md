---
title: Despliegue con Docker
tags: [despliegue, docker]
status: borrador
updated: 2026-10-09
---

# Despliegue con Docker

Contexto: [[Arquitectura]]. Pipeline que construye las imágenes: [[GitHub]].

## Imágenes
| Imagen | Contenido | Base sugerida |
|---|---|---|
| `app` | Binario Go + frontend compilado embebido | Imagen mínima (distroless o similar), usuario no root |
| `classifier` | Servicio Python + modelo fijado | Python slim, solo CPU por defecto |

Un solo contenedor `app` sirve API y frontend: menos piezas para quien instala. En desarrollo se usa el servidor de Vite aparte.

## Principios de "imagen fija"
- **Multi-stage builds:** compilar en una etapa, copiar solo el resultado.
- **Fijar versiones:** base images por digest o versión exacta, dependencias bloqueadas (`go.sum`, lockfile de npm, lockfile de Python).
- **Modelo dentro de la imagen** con revisión fijada ([[05-Clasificador-de-Emociones]]).
- **Reproducible:** el mismo commit produce el mismo comportamiento.
- **Sin secretos en la imagen:** todo por variables de entorno o secretos del orquestador.
- Usuario no root, sistema de archivos de solo lectura cuando sea posible.
- `HEALTHCHECK` o sondas apuntando a `/healthz` y `/readyz` ([[06-Contratos-e-Integracion]]).
- Escaneo de vulnerabilidades en CI (ver [[GitHub]]).

## Docker Compose (boceto)
```yaml
services:
  classifier:
    image: ghcr.io/ORG/chat-emotion-meter-classifier:${VERSION:-latest}
    healthcheck:
      test: ["CMD", "python", "-m", "app.healthcheck"]
      interval: 10s
    deploy:
      resources:
        limits: { cpus: "2", memory: 2g }

  app:
    image: ghcr.io/ORG/chat-emotion-meter-app:${VERSION:-latest}
    environment:
      TWITCH_CHANNEL: ${TWITCH_CHANNEL:?define TWITCH_CHANNEL}
      CLASSIFIER_ADDR: classifier:50051
    depends_on:
      classifier:
        condition: service_healthy
    ports:
      - "8080:8080"

  # Opcional: observabilidad (perfil aparte)
  prometheus:
    profiles: ["observability"]
    image: prom/prometheus
  grafana:
    profiles: ["observability"]
    image: grafana/grafana
```

Un archivo `compose.dev.yaml` adicional monta el código y levanta Vite con recarga en caliente.

## Uso para un streamer (objetivo de experiencia)
```
cp .env.example .env      # poner TWITCH_CHANNEL
docker compose up -d
# en OBS: Fuente de navegador -> http://localhost:8080/overlay
```

## Dónde puede correr
- En el mismo PC del streamer (la opción más simple; el consumo de CPU compite con el juego y OBS, hay que medirlo).
- En un VPS propio.
- Más adelante: orquestadores (Kubernetes) si el proyecto crece. No es objetivo de v1.

## Versionado y publicación
- Imágenes en **GitHub Container Registry (GHCR)**.
- Etiquetas: `vX.Y.Z`, `vX.Y`, `latest` (solo desde releases) y `sha-<commit>` para pruebas.
- Plataformas: `linux/amd64` y `linux/arm64`.

## Recursos
Documentar CPU/RAM mínimos medidos para cada imagen en el README una vez que existan benchmarks ([[07-Latencia-y-Rendimiento]]).
