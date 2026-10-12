import test from 'node:test';
import assert from 'node:assert/strict';
import { dominantEmotion, moodColor, moodConfidence } from './score.js';

const config = {
  emotions: {
    others: '#B0BEC5',
    joy: '#FDD835',
    sadness: '#42A5F5',
    anger: '#FF8F00',
    surprise: '#26C6DA',
    disgust: '#66BB6A',
    fear: '#9575CD',
  },
};

test('group mood selects the dominant emotion and its configured color', () => {
  const scores = { others: 0.1, joy: 0.15, sadness: 0.05, anger: 0.6, surprise: 0.05, disgust: 0.03, fear: 0.02 };
  assert.equal(dominantEmotion(scores), 'anger');
  assert.equal(moodColor(scores, config), '#FF8F00');
  assert.equal(moodConfidence(scores), 0.6);
});

test('neutral model label maps to a stable neutral color', () => {
  assert.equal(moodColor({ others: 0.8, joy: 0.1 }, config), '#B0BEC5');
});

test('missing scores safely select the neutral class', () => {
  assert.equal(dominantEmotion(null), 'others');
  assert.equal(moodConfidence(null), 0);
});
