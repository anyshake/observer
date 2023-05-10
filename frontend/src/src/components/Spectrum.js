import { Component } from "react";
import chroma from "chroma-js";
import { fft, util } from "fft-js";
import nextPow2 from "../helpers/utilities/nextPow2";

class Spectrogram {
    constructor(canvas, width, height) {
        this.baseCanvas = canvas;
        canvas.width = width;
        canvas.height = height;
        this.baseCanvasCtx = this.baseCanvas.getContext("2d");
        this.colors = this.setColor(275);
    }

    init() {
        const canvas = document.createElement("canvas");
        canvas.width = this.baseCanvas.width;
        canvas.height = this.baseCanvas.height;

        const canvasContext = canvas.getContext("2d");

        const tempCanvas = document.createElement("canvas");
        tempCanvas.width = canvas.width;
        tempCanvas.height = canvas.height;
        canvasContext._tempContext = tempCanvas.getContext("2d");

        return canvasContext;
    }

    drawGird() {
        const xGridStep = 10;
        const yGridStep = 10;

        const width = this.baseCanvas.width;
        const height = this.baseCanvas.height;

        this.baseCanvasCtx.strokeStyle = "#aaa";
        this.baseCanvasCtx.beginPath();
        for (let x = xGridStep; x < width; x += xGridStep) {
            this.baseCanvasCtx.moveTo(x, 0);
            this.baseCanvasCtx.lineTo(x, height);
        }
        for (let y = yGridStep; y < height; y += yGridStep) {
            this.baseCanvasCtx.moveTo(0, y);
            this.baseCanvasCtx.lineTo(width, y);
        }
        this.baseCanvasCtx.stroke();
    }

    drawAxis() {}

    drawFFT(dataArray, canvasContext) {
        const canvas = canvasContext.canvas;

        const width = canvas.width;
        const height = canvas.height;

        const tempCanvasContext = canvasContext._tempContext;
        const tempCanvas = tempCanvasContext.canvas;

        tempCanvasContext.drawImage(canvas, 0, 0, width, height);

        const N = nextPow2(dataArray);
        const phasors = fft(N);
        const frequencies = util.fftFreq(phasors, N.length);
        const magnitudes = util.fftMag(phasors);
        const both = frequencies.map((f, ix) => {
            return { frequency: f, magnitude: magnitudes[ix] * 50 };
        });

        for (let i = 0; i < Math.max(both.length, height); i++) {
            const data = both[Math.floor(i * (both.length / height))];
            canvasContext.fillStyle = this.getColor(data.magnitude);
            canvasContext.fillRect(width - 1, height - i, 1, 1);
        }

        canvasContext.translate(-1, 0);
        canvasContext.drawImage(
            tempCanvas,
            0,
            0,
            width,
            height,
            0,
            0,
            width,
            height
        );

        canvasContext.setTransform(1, 0, 0, 1, 0, 0);
        this.baseCanvasCtx.drawImage(canvas, 0, 0, width, height);
    }

    setColor(steps) {
        const scale = new chroma.scale(
            [
                [0, 0, 255, 1],
                [0, 255, 255, 1],
                [0, 255, 0, 1],
                [255, 255, 0, 1],
                [255, 0, 0, 1],
                [255, 0, 255, 1],
                [128, 0, 128, 1],
            ],
            [0, 0.1, 0.2, 0.3, 0.4, 0.6, 0.8]
        ).domain([0, steps]);

        return Array.from({ length: steps }, (_, index) => scale(index).hex());
    }

    getColor(index) {
        let color = this.colors[index >> 0];
        return color ? color : this.colors[this.colors.length - 1];
    }
}

export default class Spectrum extends Component {
    constructor(props) {
        super(props);
        this.state = {
            spectrum: {
                instance: null,
                context: null,
                width: 1000,
                height: 350,
            },
        };
    }

    componentDidMount() {
        const instance = new Spectrogram(
            document.getElementById("viewport"),
            this.state.spectrum.width,
            this.state.spectrum.height
        );
        const ctx = instance.init();

        instance.drawGird();
        this.setState({
            spectrum: {
                ...this.state.spectrum,
                instance: instance,
                context: ctx,
            },
        });
    }

    componentDidUpdate(prevProps) {
        if (prevProps.data.timestamp !== this.props.data.timestamp) {
            this.state.spectrum.context &&
                this.state.spectrum.instance.drawFFT(
                    this.props.data.synthesis,
                    this.state.spectrum.context
                );
        }
    }

    render() {
        return (
            <div className="mx-auto w-full">
                <div className="overflow-hidden">
                    <div className="overflow-x-scroll flex">
                        <div className="mx-auto">
                            <canvas className="rounded-lg" id="viewport" />
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
