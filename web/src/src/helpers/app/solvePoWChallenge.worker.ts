import { sha512 } from '@noble/hashes/sha2.js';
import { Buffer } from 'buffer';

self.addEventListener('message', ({ data }) => {
    const seed = Buffer.from(data.seed, 'base64');
    const difficulty = seed[0];
    const challenge = seed.subarray(1);

    let nonce = 0;
    let hashHex: string;
    do {
        const text = String.fromCharCode(...challenge) + nonce;
        const hash = sha512(new TextEncoder().encode(text));
        hashHex = Array.from(hash)
            .map((b) => b.toString(16).padStart(2, '0'))
            .join('');
        nonce++;
    } while (hashHex.slice(0, difficulty) !== '0'.repeat(difficulty));

    self.postMessage(hashHex);
});
