export const HomeConstraints = {
    lineChartRetention: 100,
    pollInterval: 2000,
    maxGapSeconds: 10
};

export const RealTimeConstraints = {
    id: 'realtime',
    minWidth: 200,
    minHeight: 150,
    maxWidth: 800,
    maxHeight: 600,
    // waveform default options
    minSpanValue: 100,
    lineColor: '#8A3EED',
    // spectrogram default options
    fftSize: 1024,
    windowSize: 512,
    overlap: Math.floor(512 * 0.86),
    freqRange: [0, 25] as [number, number],
    getDynamicDB: (index: number) => {
        if (index > 2 && index < 6) {
            return { minDB: 20, maxDB: 120 };
        }
        return { minDB: 110, maxDB: 170 };
    }
};

export const HistoryConstraints = {
    id: 'history',
    minWidth: 200,
    minHeight: 150,
    maxWidth: 800,
    maxHeight: 600,
    // waveform default options
    minSpanValue: 100,
    lineColor: '#8A3EED',
    // spectrogram default options
    fftSize: 1024,
    windowSize: 512,
    overlap: Math.floor(512 * 0.86),
    freqRange: [0, 25] as [number, number],
    getDynamicDB: (index: number) => {
        if (index > 2 && index < 6) {
            return { minDB: 20, maxDB: 120 };
        }
        return { minDB: 110, maxDB: 170 };
    }
};

export const DownloadConstraints = {};

export const SettingsConstraints = {};
