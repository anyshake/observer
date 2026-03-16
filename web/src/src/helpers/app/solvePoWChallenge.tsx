import PoWWorker from './solvePoWChallenge.worker.ts?worker';

export const solvePoWChallenge = (seed: string): Promise<{ nonce: number; hash: string }> => {
    return new Promise((resolve, reject) => {
        const worker = new PoWWorker();
        worker.onmessage = (event) => {
            worker.terminate();
            resolve(event.data);
        };
        worker.onerror = (err) => {
            worker.terminate();
            reject(err);
        };
        worker.postMessage({ seed });
    });
};
