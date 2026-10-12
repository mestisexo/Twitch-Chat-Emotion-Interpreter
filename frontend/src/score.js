const emotionLabels = ['others', 'joy', 'sadness', 'anger', 'surprise', 'disgust', 'fear'];

export async function loadScoreConfig(baseUrl) {
  const response = await fetch(`${baseUrl}scores.json`);
  if (!response.ok) throw new Error(`No se pudo cargar scores.json (HTTP ${response.status})`);

  const scoreConfig = await response.json();
  if (
    !scoreConfig.emotions
    || emotionLabels.some((label) =>
      typeof scoreConfig.emotions[label] !== 'string'
      || !/^#[\da-f]{6}$/i.test(scoreConfig.emotions[label])
    )
  ) {
    throw new Error('scores.json no contiene colores válidos para todas las emociones');
  }
  return scoreConfig;
}

export function dominantEmotion(scores) {
  if (!scores || typeof scores !== 'object') return 'others';
  return emotionLabels.reduce((dominant, label) =>
    Number.isFinite(scores[label]) && scores[label] > (scores[dominant] || 0) ? label : dominant,
  'others');
}

export function moodColor(scores, scoreConfig) {
  return scoreConfig.emotions[dominantEmotion(scores)];
}

export function moodConfidence(scores) {
  return Math.min(1, Math.max(0, scores?.[dominantEmotion(scores)] || 0));
}
