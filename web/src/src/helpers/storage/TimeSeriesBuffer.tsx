export default class TimeSeriesBuffer {
    private maxDuration: number;
    private buffer: Array<[number, number | null]>; // [timestamp, value]
    private lastSampleRate: number | null = null;

    constructor(maxDuration: number) {
        this.maxDuration = maxDuration * 1000;
        this.buffer = [];
    }

    addData(
        values: (number | null)[],
        recordTime: number,
        currentTime: number,
        sampleRate: number
    ) {
        if (this.lastSampleRate !== null && this.lastSampleRate !== sampleRate) {
            this.clear();
        }

        this.lastSampleRate = sampleRate;

        const interval = 1000 / sampleRate;
        for (let i = 0; i < values.length; i++) {
            this.buffer.push([recordTime + i * interval, values[i]]);
        }
        this.cleanup(currentTime);
    }

    getData(): Array<[number, number | null]> {
        if (this.buffer.length === 0) {
            return [];
        }

        this.buffer.sort((a, b) => a[0] - b[0]);

        const processedData: Array<[number, number | null]> = [];
        let lastTime = this.buffer[0][0];

        for (const [timestamp, value] of this.buffer) {
            if (timestamp - lastTime > 2000) {
                processedData.push([lastTime + 1, null]);
            }

            processedData.push([timestamp, value]);
            lastTime = timestamp;
        }

        return processedData;
    }

    getStartTime() {
        return this.buffer.length > 0 ? this.buffer[0][0] : 0;
    }

    getEndTime() {
        return this.buffer.length > 0 ? this.buffer[this.buffer.length - 1][0] : 0;
    }

    clear() {
        this.buffer = [];
        this.lastSampleRate = null;
    }

    private cleanup(currentTime: number) {
        const cutoff = currentTime - this.maxDuration;
        const index = this.buffer.findIndex(([timestamp]) => timestamp >= cutoff);
        if (index > 0) {
            this.buffer = this.buffer.slice(index);
        }
    }
}
