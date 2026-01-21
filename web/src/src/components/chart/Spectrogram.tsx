import { mdiClose, mdiCog } from '@mdi/js';
import Icon from '@mdi/react';
import { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { FFTExecutor, Spectrogram as SpectrogramCore } from 'spectrogram-js';

interface ISpectrogram {
    readonly title?: string;
    readonly sampleRate: number;
    readonly freqRange: [number, number];
    readonly duration: number;
    readonly minDB: number;
    readonly maxDB: number;
    readonly windowSize: number;
    readonly overlap: number;
    readonly data: [number, number][];
    readonly fftExecutor?: FFTExecutor;
    readonly renderFPS?: number;
    readonly onSpectrogramUpdate?: (minDB: number, maxDB: number) => void;
}

export const Spectrogram = memo(
    ({
        title,
        sampleRate,
        minDB,
        maxDB,
        freqRange,
        duration,
        windowSize,
        overlap,
        data,
        fftExecutor,
        renderFPS = 2,
        onSpectrogramUpdate
    }: ISpectrogram) => {
        const { t } = useTranslation();

        const [showSettings, setShowSettings] = useState(false);
        const [minDBState, setMinDBState] = useState(minDB);
        const [maxDBState, setMaxDBState] = useState(maxDB);

        const [initialized, setInitialized] = useState(false);
        const spectrogramRef = useRef<SpectrogramCore | null>(null);
        useEffect(() => {
            const sp = new SpectrogramCore({
                sampleRate,
                windowSize,
                overlap,
                fftExecutor,
                minDb: minDB,
                maxDb: maxDB,
                windowType: 'hann'
            });
            spectrogramRef.current = sp;
            sp.init().then(() => {
                sp.setColormap('jet');
                setInitialized(true);
            });
            return () => {
                sp.destroy();
                spectrogramRef.current = null;
            };
        }, [fftExecutor, maxDB, minDB, overlap, sampleRate, windowSize]);

        const canvasRef = useRef<HTMLCanvasElement>(null);
        const sizeRef = useRef({ width: 1, height: 1 });
        useEffect(() => {
            const canvas = canvasRef.current;
            if (!canvas || !canvas.parentElement) {
                return;
            }

            const ro = new ResizeObserver(([entry]) => {
                const { width, height } = entry.contentRect;
                sizeRef.current.width = Math.max(1, Math.floor(width));
                sizeRef.current.height = Math.max(1, Math.floor(height));
            });

            ro.observe(canvas.parentElement);
            return () => ro.disconnect();
        }, []);

        const [timePercent, setTimePercent] = useState(0);
        const dataDuration = useMemo<number | null>(() => {
            if (data.length === 0) {
                return null;
            }
            return (data[data.length - 1][0] - data[0][0]) / 1000;
        }, [data]);

        useEffect(() => {
            if (!initialized) {
                return;
            }

            const id = setInterval(() => {
                const canvas = canvasRef.current;
                const sp = spectrogramRef.current;
                if (!canvas || !sp) {
                    return;
                }

                const { width, height } = sizeRef.current;
                if (width <= 0 || height <= 0) {
                    return;
                }

                const total = dataDuration ?? 0;
                const window = duration;
                const maxStart = Math.max(0, total - window);
                const start = maxStart * timePercent;
                const end = start + window;
                sp.render({
                    canvas,
                    width,
                    height,
                    timeRange: [start, end],
                    freqRange
                });
            }, 1000 / renderFPS);

            return () => clearInterval(id);
        }, [data, dataDuration, duration, freqRange, initialized, renderFPS, timePercent]);

        useEffect(() => {
            if (initialized) {
                spectrogramRef.current?.setData(data);
            }
        }, [data, initialized]);

        const handlePreviewMinDB = useCallback((value: number) => {
            setMinDBState(value);
            setMaxDBState((prev) => (value > prev ? value : prev));
        }, []);

        const handleApplyMinDB = useCallback(
            (value: number) => {
                spectrogramRef.current?.updateConfig({
                    minDb: value,
                    maxDb: Math.max(value, maxDBState)
                });
                onSpectrogramUpdate?.(value, Math.max(value, maxDBState));
            },
            [maxDBState, onSpectrogramUpdate, spectrogramRef]
        );

        const handlePreviewMaxDB = useCallback((value: number) => {
            setMaxDBState(value);
            setMinDBState((prev) => (value < prev ? value : prev));
        }, []);

        const handleApplyMaxDB = useCallback(
            (value: number) => {
                spectrogramRef.current?.updateConfig({
                    minDb: Math.min(value, minDBState),
                    maxDb: value
                });
                onSpectrogramUpdate?.(Math.min(value, minDBState), value);
            },
            [minDBState, onSpectrogramUpdate, spectrogramRef]
        );

        return (
            <div className="relative flex h-full w-full flex-col">
                <canvas className="block" ref={canvasRef} />

                <div className="absolute top-5 left-15 flex items-center space-x-1">
                    {title && (
                        <div className="flex h-6 items-center rounded bg-black/50 px-3 text-sm font-bold text-white select-none">
                            {title}
                        </div>
                    )}

                    <button
                        className="flex size-6 cursor-pointer items-center justify-center rounded bg-black/50 text-white opacity-50 transition-all hover:opacity-100"
                        onClick={() => setShowSettings((v) => !v)}
                    >
                        <Icon path={mdiCog} size={0.6} />
                    </button>

                    <div className="flex h-6 items-center rounded bg-black/50 opacity-50 transition-all hover:opacity-100">
                        <input
                            type="range"
                            min={0}
                            max={1}
                            step={0.01}
                            value={timePercent}
                            onChange={({ target }) => setTimePercent(Number(target.value))}
                            className="range range-xs rounded py-3 pl-2 text-gray-300 [--range-bg:#fff] [--range-fill:0]"
                        />
                        <span className="px-2 font-mono text-xs text-white select-none">
                            {(timePercent * 100).toFixed(0)}%
                        </span>
                    </div>
                </div>

                {showSettings && (
                    <div
                        className="absolute inset-0 z-10 flex cursor-default items-center justify-center rounded-md bg-black/50"
                        onClick={() => setShowSettings(false)}
                    >
                        <div className="max-h-full w-full overflow-y-auto px-4 py-2">
                            <div
                                className="mx-auto w-64 space-y-4 rounded bg-black/90 p-4 text-xs text-white"
                                onClick={(e) => e.stopPropagation()}
                            >
                                <div className="flex justify-between text-sm font-bold">
                                    <span>{t('components.Spectrogram.settings.title')}</span>
                                    <button onClick={() => setShowSettings(false)}>
                                        <Icon path={mdiClose} size={0.6} />
                                    </button>
                                </div>

                                <>
                                    <div className="mb-1 flex justify-between">
                                        <span>{t('components.Spectrogram.settings.min_db')}</span>
                                        <span>{minDBState}</span>
                                    </div>
                                    <input
                                        type="range"
                                        min={-250}
                                        max={250}
                                        step={2}
                                        value={minDBState}
                                        className="range range-info range-sm w-full"
                                        onChange={({ target }) =>
                                            handlePreviewMinDB(Number(target.value))
                                        }
                                        onMouseUp={() => handleApplyMinDB(minDBState)}
                                        onTouchEnd={() => handleApplyMinDB(minDBState)}
                                    />
                                </>

                                <>
                                    <div className="mb-1 flex justify-between">
                                        <span>{t('components.Spectrogram.settings.max_db')}</span>
                                        <span>{maxDBState}</span>
                                    </div>
                                    <input
                                        type="range"
                                        min={-250}
                                        max={250}
                                        step={2}
                                        value={maxDBState}
                                        className="range range-info range-sm w-full"
                                        onChange={({ target }) =>
                                            handlePreviewMaxDB(Number(target.value))
                                        }
                                        onMouseUp={() => handleApplyMaxDB(maxDBState)}
                                        onTouchEnd={() => handleApplyMaxDB(maxDBState)}
                                    />
                                </>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        );
    }
);
