export default class TimeSeriesBuffer {
    private maxDuration: number;
    private buffer: Array<[number, number | null]>;
    private head = 0;
    private lastSampleRate: number | null = null;
    private sorted = true;

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

        const count = values.length;
        if (count > 0) {
            const interval = 1000 / sampleRate;
            const live = this.buffer.length - this.head;
            if (live > 0 && (interval < 0 || recordTime < this.buffer[this.buffer.length - 1][0])) {
                this.sorted = false;
            }

            for (let i = 0; i < count; i++) {
                this.buffer.push([recordTime + i * interval, values[i]]);
            }
        }

        this.cleanup(currentTime);
    }

    getData(): Array<[number, number | null]> {
        const end = this.buffer.length;
        if (this.head >= end) {
            return [];
        }

        this.ensureSorted();

        const processedData: Array<[number, number | null]> = [];
        let lastTime = this.buffer[this.head][0];

        for (let i = this.head; i < this.buffer.length; i++) {
            const point = this.buffer[i];
            const timestamp = point[0];
            if (timestamp - lastTime > 2000) {
                processedData.push([lastTime + 1, null]);
            }
            processedData.push([timestamp, point[1]]);
            lastTime = timestamp;
        }

        return processedData;
    }

    getSamples(): Array<[number, number]> {
        const end = this.buffer.length;
        if (this.head >= end) {
            return [];
        }

        this.ensureSorted();

        const samples: Array<[number, number]> = [];
        for (let i = this.head; i < this.buffer.length; i++) {
            const point = this.buffer[i];
            if (point[1] !== null) {
                samples.push([point[0], point[1]]);
            }
        }
        return samples;
    }

    getStartTime() {
        return this.buffer.length > this.head ? this.buffer[this.head][0] : 0;
    }

    getEndTime() {
        return this.buffer.length > this.head ? this.buffer[this.buffer.length - 1][0] : 0;
    }

    clear() {
        this.buffer = [];
        this.head = 0;
        this.lastSampleRate = null;
        this.sorted = true;
    }

    private ensureSorted() {
        if (this.sorted || this.head >= this.buffer.length) {
            return;
        }
        this.compact();
        this.buffer.sort((a, b) => a[0] - b[0]);
        this.sorted = true;
    }

    private compact() {
        if (this.head === 0) {
            return;
        }
        this.buffer = this.buffer.slice(this.head);
        this.head = 0;
    }

    private cleanup(currentTime: number) {
        const end = this.buffer.length;
        if (this.head >= end) {
            return;
        }

        const cutoff = currentTime - this.maxDuration;
        let dropTo = this.head;

        if (this.sorted) {
            let lo = this.head;
            let hi = end;
            while (lo < hi) {
                const mid = (lo + hi) >> 1;
                if (this.buffer[mid][0] < cutoff) {
                    lo = mid + 1;
                } else {
                    hi = mid;
                }
            }
            if (lo === end) {
                return;
            }
            dropTo = lo;
        } else {
            let found = -1;
            for (let i = this.head; i < end; i++) {
                if (this.buffer[i][0] >= cutoff) {
                    found = i;
                    break;
                }
            }
            if (found < 0) {
                return;
            }
            dropTo = found;
        }

        if (dropTo <= this.head) {
            return;
        }

        this.head = dropTo;
        const live = this.buffer.length - this.head;
        if (this.head > 0 && this.head >= live) {
            this.compact();
        }
    }
}
