export type SpectrumPoints = [number, number | null][];

const nextPowerOfTwo = (value: number) => 2 ** Math.ceil(Math.log2(Math.max(2, value)));

const fft = (re: Float64Array, im: Float64Array) => {
    const size = re.length;
    for (let i = 1, j = 0; i < size; i++) {
        let bit = size >> 1;
        for (; j & bit; bit >>= 1) {
            j ^= bit;
        }
        j ^= bit;
        if (i < j) {
            [re[i], re[j]] = [re[j], re[i]];
            [im[i], im[j]] = [im[j], im[i]];
        }
    }

    const cos = new Float64Array(size / 2);
    const sin = new Float64Array(size / 2);
    for (let k = 0; k < size / 2; k++) {
        cos[k] = Math.cos((2 * Math.PI * k) / size);
        sin[k] = -Math.sin((2 * Math.PI * k) / size);
    }
    for (let length = 2; length <= size; length <<= 1) {
        const half = length >> 1;
        const step = size / length;
        for (let start = 0; start < size; start += length) {
            for (let k = 0; k < half; k++) {
                const a = start + k;
                const b = a + half;
                const wr = cos[k * step];
                const wi = sin[k * step];
                const tr = re[b] * wr - im[b] * wi;
                const ti = re[b] * wi + im[b] * wr;
                re[b] = re[a] - tr;
                im[b] = im[a] - ti;
                re[a] += tr;
                im[a] += ti;
            }
        }
    }
};

const hann = (index: number, length: number) =>
    length > 1 ? 0.5 - 0.5 * Math.cos((2 * Math.PI * index) / (length - 1)) : 1;

const toLogPoints = (frequencies: Float64Array, amplitudes: Float64Array, maxPoints: number) => {
    const points: SpectrumPoints = [];
    const push = (frequency: number, amplitude: number) =>
        points.push([frequency, amplitude > 0 ? amplitude : null]);
    if (frequencies.length <= maxPoints) {
        frequencies.forEach((frequency, index) => push(frequency, amplitudes[index]));
        return points;
    }

    const first = Math.log10(frequencies[0]);
    const span = Math.log10(frequencies[frequencies.length - 1]) - first;
    let bin = -1;
    let peak = -1;
    let peakFrequency = 0;
    frequencies.forEach((frequency, index) => {
        const current = Math.min(
            maxPoints - 1,
            Math.floor(((Math.log10(frequency) - first) / span) * maxPoints)
        );
        if (current !== bin && bin >= 0) {
            push(peakFrequency, peak);
            peak = -1;
        }
        bin = current;
        if (amplitudes[index] > peak) {
            peak = amplitudes[index];
            peakFrequency = frequency;
        }
    });
    push(peakFrequency, peak);
    const last = frequencies.length - 1;
    if (peakFrequency !== frequencies[last]) {
        push(frequencies[last], amplitudes[last]);
    }
    return points;
};

export const getAmplitudeSpectrum = (
    samples: Float64Array,
    sampleRate: number,
    maxPoints: number,
    minFrequency: number
): SpectrumPoints => {
    const count = samples.length;
    if (count < 2 || sampleRate <= 0) {
        return [];
    }

    let mean = 0;
    for (let index = 0; index < count; index++) {
        mean += samples[index] / count;
    }
    const size = nextPowerOfTwo(count);
    const re = new Float64Array(size);
    const im = new Float64Array(size);
    let windowSum = 0;
    for (let index = 0; index < count; index++) {
        const weight = hann(index, count);
        re[index] = (samples[index] - mean) * weight;
        windowSum += weight;
    }
    fft(re, im);

    const first = Math.max(1, Math.ceil((minFrequency * size) / sampleRate));
    const last = size / 2;
    if (first > last) {
        return [];
    }
    const frequencies = new Float64Array(last - first + 1);
    const amplitudes = new Float64Array(last - first + 1);
    for (let k = first; k <= last; k++) {
        frequencies[k - first] = (k * sampleRate) / size;
        amplitudes[k - first] = (2 * Math.hypot(re[k], im[k])) / windowSum;
    }
    return toLogPoints(frequencies, amplitudes, maxPoints);
};
