import { memo, useEffect, useState } from 'react';

import { SpectrumConstraints } from '../../config/constraints';
import type { SpectrumPoints } from '../../helpers/dsp/getAmplitudeSpectrum';
import { SpectrumExecutor } from '../../helpers/dsp/SpectrumExecutor';
import { LogAxisChart } from './LogAxisChart';
import type { IPlotBounds } from './usePlotBounds';

interface ISpectrumChart {
    readonly samples: Float64Array | null;
    readonly sampleRate: number;
    readonly unit?: string;
    readonly executor: SpectrumExecutor;
    readonly bounds: IPlotBounds | null;
}

const EMPTY_POINTS: SpectrumPoints = [];

export const SpectrumChart = memo(function SpectrumChart({
    samples,
    sampleRate,
    unit,
    executor,
    bounds
}: ISpectrumChart) {
    const [points, setPoints] = useState(EMPTY_POINTS);

    useEffect(() => {
        if (!samples || samples.length < SpectrumConstraints.minSamples || sampleRate <= 0) {
            setPoints(EMPTY_POINTS);
            return;
        }
        let stale = false;
        executor
            .compute(
                samples,
                sampleRate,
                SpectrumConstraints.maxPoints,
                SpectrumConstraints.minFrequency
            )
            .then((result) => {
                if (!stale) {
                    setPoints(result);
                }
            })
            .catch(() => {
                if (!stale) {
                    setPoints(EMPTY_POINTS);
                }
            });
        return () => {
            stale = true;
        };
    }, [executor, samples, sampleRate]);

    return (
        <LogAxisChart
            height="100%"
            data={points}
            xName="Frequency (Hz)"
            yName={unit}
            xMin={SpectrumConstraints.minFrequency}
            yLog
            yPosition="right"
            bounds={bounds}
        />
    );
});
