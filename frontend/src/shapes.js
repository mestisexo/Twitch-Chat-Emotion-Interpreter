import { Container, Graphics } from 'pixi.js';

export const shapeSides = {
  triangle: 3,
  square: 4,
  pentagon: 5,
  hexagon: 6,
  heptagon: 7,
  octagon: 8,
};

export const defaultShape = 'triangle';
export const figureRadius = 20;

export function normalizeShape(shape) {
  return shapeSides[shape] ? shape : defaultShape;
}

export function drawShape(figure, shape, color) {
  const normalizedShape = normalizeShape(shape);
  const points = polygonPoints(normalizedShape, figureRadius);

  const graphic = new Graphics().poly(points).fill(color);
  if (figure.graphic) {
    figure.visualContainer.removeChild(figure.graphic);
    figure.graphic.destroy();
  }
  figure.graphic = graphic;
  figure.visualContainer.addChild(graphic);
}

export function createActivityEffect(shape, color) {
  const effectContainer = new Container();
  const outline = new Graphics();
  effectContainer.addChild(outline);

  function update(nextShape, nextColor) {
    const outlinePoints = polygonPoints(nextShape, figureRadius);
    outline.clear();
    outline.poly(outlinePoints).stroke({ color: nextColor, width: 8, alpha: 0.9 });
  }
  update(shape, color);

  return {
    container: effectContainer,
    update,
    animate(elapsed, phase, confidence) {
      const pulse = (Math.sin(elapsed * 4 + phase) + 1) / 2;
      const shimmer = (Math.sin(elapsed * 17 + phase * 2.3) + 1) / 2;
      const flicker = shimmer > 0.9 ? 0.12 : 0;
      outline.alpha = Math.min(1, 0.25 + confidence * (0.45 + pulse * 0.25 + flicker));
      return 1 + confidence * (pulse * 0.035 + flicker * 0.02);
    },
  };
}

function polygonPoints(shape, radius, offsetX = 0, offsetY = 0) {
  const sides = shapeSides[normalizeShape(shape)];
  const points = [];
  const startAngle = -Math.PI / 2;
  for (let i = 0; i < sides; i++) {
    const angle = startAngle + (i * Math.PI * 2) / sides;
    points.push(offsetX + Math.cos(angle) * radius, offsetY + Math.sin(angle) * radius);
  }
  return points;
}
