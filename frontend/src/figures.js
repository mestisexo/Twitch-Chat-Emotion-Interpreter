import { Container } from 'pixi.js';
import { dominantEmotion } from './score.js';
import { createActivityEffect, drawShape, figureRadius, normalizeShape } from './shapes.js';

const visibilityTransitionDuration = 350;

export function createFigureManager(app, scoreConfig) {
  const shape = normalizeShape(scoreConfig.shape);
  const figure = {
    container: new Container(),
    visualContainer: new Container(),
    graphic: null,
    effect: null,
    color: '',
    confidence: 0,
    scores: null,
    phase: Math.random() * Math.PI * 2,
    phaseY: Math.random() * Math.PI * 2,
    frequencyX: 0.35 + Math.random() * 0.45,
    frequencyX2: 0.8 + Math.random() * 0.9,
    frequencyY: 0.8 + Math.random() * 0.8,
    amplitudeY: 18 + Math.random() * 6,
    appearedAt: 0,
  };
  let elapsed = 0;

  figure.container.visible = false;
  figure.container.addChild(figure.visualContainer);
  app.stage.addChild(figure.container);

  setMood({ others: 1 });

  function setMood(scores) {
    if (!scores || typeof scores !== 'object') {
      scores = { others: 1 };
    }

    if (!figure.container.visible) {
      figure.container.visible = true;
      figure.appearedAt = Date.now();
    }
    figure.scores = scores;
    const emotion = dominantEmotion(scores);
    figure.confidence = Math.min(1, Math.max(0, scores[emotion] || 0));
    const color = scoreConfig.emotions[emotion];
    if (figure.color === color) return;

    const colorValue = Number.parseInt(color.slice(1), 16);
    drawShape(figure, shape, colorValue);
    figure.color = color;
    if (!figure.effect) {
      figure.effect = createActivityEffect(shape, colorValue);
      figure.visualContainer.addChild(figure.effect.container);
    } else {
      figure.effect.update(shape, colorValue);
    }
  }

  function applyEvent(event) {
    switch (event.type) {
      case 'snapshot':
      case 'mood':
        setMood(event.scores);
        break;
    }
  }

  const animate = (ticker) => {
    if (!figure.container.visible) return;
    elapsed += ticker.deltaTime / 60;

    const now = Date.now();
    const entranceProgress = Math.min(1, Math.max(0, (now - figure.appearedAt) / visibilityTransitionDuration));
    const easedProgress = entranceProgress * entranceProgress * (3 - 2 * entranceProgress);
    figure.container.alpha = easedProgress;
    figure.container.scale.set(0.82 + easedProgress * 0.18);

    const horizontalPhase = elapsed * figure.frequencyX + figure.phase;
    const secondaryHorizontalPhase = elapsed * figure.frequencyX2 + figure.phaseY;
    const x = app.screen.width / 2
      + Math.sin(horizontalPhase) * app.screen.width * 0.4
      + Math.sin(secondaryHorizontalPhase) * app.screen.width * 0.08;
    const y = app.screen.height - 50
      + Math.sin(elapsed * figure.frequencyY + figure.phaseY) * figure.amplitudeY
      + Math.cos(horizontalPhase) * 6;
    figure.container.x = Math.max(figureRadius, Math.min(app.screen.width - figureRadius, x));
    figure.container.y = Math.max(app.screen.height - 80, Math.min(app.screen.height - 20, y));

    const horizontalVelocity = Math.cos(horizontalPhase) * figure.frequencyX * 0.4
      + Math.cos(secondaryHorizontalPhase) * figure.frequencyX2 * 0.08;
    const stretch = Math.max(-0.12, Math.min(0.12, horizontalVelocity * 0.18));
    const wobble = Math.sin(elapsed * figure.frequencyY * 2.5 + figure.phaseY) * 0.035;
    const pulse = figure.effect?.animate(elapsed, figure.phase, figure.confidence) || 1;
    figure.visualContainer.scale.set(
      (1 + stretch + wobble) * pulse,
      (1 - stretch * 0.7 - wobble * 0.7) * pulse,
    );
    figure.visualContainer.rotation = horizontalVelocity * 0.09
      + Math.sin(elapsed * figure.frequencyX + figure.phase) * 0.025;
  };
  app.ticker.add(animate);

  return {
    applyEvent,
    clear() {
      app.ticker.remove(animate);
      app.stage.removeChild(figure.container);
      figure.container.destroy({ children: true });
    },
  };
}
