import { getAmplitudeSpectrum } from './getAmplitudeSpectrum';

self.addEventListener('message', ({ data }) => {
    const { id, samples, sampleRate, maxPoints, minFrequency } = data;
    self.postMessage({
        id,
        points: getAmplitudeSpectrum(samples, sampleRate, maxPoints, minFrequency)
    });
});
