export const createLowPassFilter = (
    cutoffFreq: number,
    sampleRate: number,
    numTaps: number
): number[] => {
    const normalizedCutoff = (2 * Math.PI * cutoffFreq) / sampleRate;
    const coeffs = new Array(numTaps);

    for (let i = 0; i < numTaps; i++) {
        if (i === Math.floor(numTaps / 2)) {
            coeffs[i] = normalizedCutoff / Math.PI;
        } else {
            const n = i - Math.floor(numTaps / 2);
            coeffs[i] = Math.sin(normalizedCutoff * n) / (Math.PI * n);
        }
    }

    const window = new Array(numTaps);
    for (let i = 0; i < numTaps; i++) {
        window[i] = 0.54 - 0.46 * Math.cos((2 * Math.PI * i) / (numTaps - 1));
        coeffs[i] *= window[i];
    }

    return coeffs;
};

export const createHighPassFilter = (
    cutoffFreq: number,
    sampleRate: number,
    numTaps: number
): number[] => {
    const lowPass = createLowPassFilter(cutoffFreq, sampleRate, numTaps);
    const highPass = new Array(numTaps);

    for (let i = 0; i < numTaps; i++) {
        highPass[i] = i === Math.floor(numTaps / 2) ? 1 - lowPass[i] : -lowPass[i];
    }

    return highPass;
};

export const createBandPassFilter = (
    lowCutoffFreq: number,
    highCutoffFreq: number,
    sampleRate: number,
    numTaps: number
): number[] => {
    const lowPass = createLowPassFilter(highCutoffFreq, sampleRate, numTaps);
    const highPass = createHighPassFilter(lowCutoffFreq, sampleRate, numTaps);
    const bandPass = new Array(numTaps);

    for (let i = 0; i < numTaps; i++) {
        bandPass[i] = lowPass[i] + highPass[i];
    }

    return bandPass;
};
