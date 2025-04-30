export default class FIRFilter {
    kernel: number[];
    state: number[];

    constructor(kernel: number[]) {
        this.kernel = kernel;
        this.state = new Array(kernel.length - 1).fill(0);
    }

    apply(input: number[]): number[] {
        const numTaps = this.kernel.length;
        const extendedInput = this.state.concat(input);
        const output = new Array(input.length).fill(0);

        for (let i = 0; i < input.length; i++) {
            for (let j = 0; j < numTaps; j++) {
                output[i] += extendedInput[i + j] * this.kernel[j];
            }
        }

        this.state = extendedInput.slice(extendedInput.length - (numTaps - 1));
        return output;
    }
}
