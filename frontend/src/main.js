import { Application } from 'pixi.js';
import './style.css';
import { createFigureManager } from './figures.js';
import { loadScoreConfig } from './score.js';

const params = new URLSearchParams(window.location.search);
const wsUrl = params.get('ws') || `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.hostname || 'localhost'}:8080/ws`;
const app = new Application();
let socket;
let retryTimer;
let figureManager;

async function initialize() {
  const scoreConfig = await loadScoreConfig(import.meta.env.BASE_URL);
  await app.init({ resizeTo: window, backgroundColor: 0x00ff00, backgroundAlpha: 1, antialias: true });
  app.ticker.maxFPS = 30;
  document.querySelector('#app').append(app.canvas);
  figureManager = createFigureManager(app, scoreConfig);
  connect();
}

function connect() {
  socket = new WebSocket(wsUrl);
  socket.addEventListener('close', () => {
    retryTimer = window.setTimeout(connect, 2000);
  });
  socket.addEventListener('error', () => socket.close());
  socket.addEventListener('message', ({ data }) => {
    try {
      figureManager.applyEvent(JSON.parse(data));
    } catch (error) {
      console.warn('Evento WebSocket inválido', error);
    }
  });
}

initialize().catch((error) => {
  console.error('No se pudo inicializar la app:', error);
});

window.addEventListener('beforeunload', () => {
  clearTimeout(retryTimer);
  socket?.close();
  figureManager?.clear();
});
