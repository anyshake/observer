import SpectrumWorker from './computeSpectrum.worker.ts?worker';
import type { SpectrumPoints } from './getAmplitudeSpectrum';

interface ISpectrumTask {
    readonly resolve: (points: SpectrumPoints) => void;
    readonly reject: (reason: Error) => void;
}

interface ISpectrumResult {
    readonly id: number;
    readonly points: SpectrumPoints;
}

export class SpectrumExecutor {
    private worker: Worker | null = null;
    private nextId = 0;
    private readonly tasks = new Map<number, ISpectrumTask>();

    compute(samples: Float64Array, sampleRate: number, maxPoints: number, minFrequency: number) {
        const worker = this.getWorker();
        const id = ++this.nextId;
        return new Promise<SpectrumPoints>((resolve, reject) => {
            this.tasks.set(id, { resolve, reject });
            worker.postMessage({ id, samples, sampleRate, maxPoints, minFrequency });
        });
    }

    dispose(reason = 'spectrum executor disposed') {
        this.worker?.terminate();
        this.worker = null;
        this.tasks.forEach(({ reject }) => reject(new Error(reason)));
        this.tasks.clear();
    }

    private getWorker() {
        if (this.worker) {
            return this.worker;
        }
        const worker = new SpectrumWorker();
        worker.onmessage = ({ data: { id, points } }: MessageEvent<ISpectrumResult>) => {
            this.tasks.get(id)?.resolve(points);
            this.tasks.delete(id);
        };
        worker.onerror = (event) => this.dispose(event.message || 'spectrum worker error');
        this.worker = worker;
        return worker;
    }
}
