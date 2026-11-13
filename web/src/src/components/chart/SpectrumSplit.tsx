import { ReactNode, useCallback, useState } from 'react';

import { SpectrumExecutor } from '../../helpers/dsp/SpectrumExecutor';
import { SpectrumChart } from './SpectrumChart';
import type { IPlotBounds, PlotBoundsHandler } from './usePlotBounds';
import type { VisibleSamplesHandler } from './useVisibleSamples';

interface ISpectrumSplit {
    readonly enabled: boolean;
    readonly sampleRate: number;
    readonly unit?: string;
    readonly executor: SpectrumExecutor;
    readonly children: (
        onVisibleSamples?: VisibleSamplesHandler,
        onPlotBounds?: PlotBoundsHandler
    ) => ReactNode;
}

export const SpectrumSplit = ({
    enabled,
    sampleRate,
    unit,
    executor,
    children
}: ISpectrumSplit) => {
    const [samples, setSamples] = useState<Float64Array | null>(null);
    const [bounds, setBounds] = useState<IPlotBounds | null>(null);
    const [prevEnabled, setPrevEnabled] = useState(enabled);
    if (prevEnabled !== enabled) {
        setPrevEnabled(enabled);
        setSamples(null);
    }

    const handlePlotBounds = useCallback((next: IPlotBounds) => {
        setBounds((prev) =>
            prev && prev.left === next.left && prev.right === next.right ? prev : next
        );
    }, []);

    return (
        <div className="app-chart-split">
            <div className="app-chart-split-pane">
                {enabled ? children(setSamples, handlePlotBounds) : children()}
            </div>
            {enabled && (
                <div className="app-chart-split-pane app-chart-split-spectrum">
                    <SpectrumChart
                        samples={samples}
                        sampleRate={sampleRate}
                        executor={executor}
                        bounds={bounds}
                        unit={unit}
                    />
                </div>
            )}
        </div>
    );
};
