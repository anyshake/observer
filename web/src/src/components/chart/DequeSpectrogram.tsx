import { mdiClose, mdiCog } from '@mdi/js';
import Icon from '@mdi/react';
import {
    forwardRef,
    memo,
    useCallback,
    useEffect,
    useImperativeHandle,
    useRef,
    useState
} from 'react';
import { useTranslation } from 'react-i18next';
import { FFTExecutor, Spectrogram as SpectrogramCore } from 'spectrogram-js';

import TimeSeriesBuffer from '../../helpers/storage/TimeSeriesBuffer';

export interface DequeSpectrogramHandle {
    addData(values: number[], recordTime: number, currentTime: number, sampleRate: number): void;
}

interface ISpectrogramDeque {
    readonly title?: string;
    readonly sampleRate: number;
    readonly duration: number;
    readonly freqRange: [number, number];
    readonly minDB: number;
    readonly maxDB: number;
    readonly windowSize: number;
    readonly overlap: number;
    readonly fftExecutor?: FFTExecutor;
    readonly renderFPS?: number;
    readonly onSpectrogramUpdate?: (minDB: number, maxDB: number) => void;
}

export const DequeSpectrogram = memo(
    forwardRef<DequeSpectrogramHandle, ISpectrogramDeque>(
        (
            {
                title,
                sampleRate,
                duration,
                minDB,
                maxDB,
                freqRange,
                windowSize,
                overlap,
                fftExecutor,
                renderFPS = 2,
                onSpectrogramUpdate
            },
            ref
        ) => {
            const { t } = useTranslation();

            const [showSettings, setShowSettings] = useState(false);
            const [minDBState, setMinDBState] = useState(minDB);
            const [maxDBState, setMaxDBState] = useState(maxDB);

            const [initialized, setInitialized] = useState(false);
            const spectrogramRef = useRef<SpectrogramCore | null>(null);
            useEffect(() => {
                const sp = new SpectrogramCore({
                    overlap,
                    sampleRate,
                    windowSize,
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

            const bufferRef = useRef<TimeSeriesBuffer>(new TimeSeriesBuffer(duration));
            useEffect(() => {
                bufferRef.current = new TimeSeriesBuffer(duration);
            }, [duration]);

            const canvasRef = useRef<HTMLCanvasElement>(null);
            const sizeRef = useRef({ width: 0, height: 0 });
            const needsUpdateRef = useRef(false);
            useEffect(() => {
                const canvas = canvasRef.current;
                if (!canvas || !canvas.parentElement) {
                    return;
                }

                const ro = new ResizeObserver(([entry]) => {
                    if (!entry) {
                        return;
                    }
                    const { width, height } = entry.contentRect;
                    if (width !== sizeRef.current.width || height !== sizeRef.current.height) {
                        sizeRef.current.width = Math.max(1, Math.floor(width));
                        sizeRef.current.height = Math.max(1, Math.floor(height));
                        needsUpdateRef.current = true;
                    }
                });

                ro.observe(canvas.parentElement);
                return () => ro.disconnect();
            }, []);

            const addData = useCallback(
                (values: number[], recordTime: number, currentTime: number, sr: number) => {
                    bufferRef.current.addData(values, recordTime, currentTime, sr);
                    needsUpdateRef.current = true;
                },
                []
            );

            useImperativeHandle(ref, () => ({ addData }), [addData]);

            const timeRangeRef = useRef<[number, number]>([0, 0.001]);
            useEffect(() => {
                let rafId: number;
                let lastRenderTime = 0;
                const frameInterval = 1000 / renderFPS;
                const renderLoop = (now: number) => {
                    rafId = requestAnimationFrame(renderLoop);

                    if (now - lastRenderTime < frameInterval) {
                        return;
                    }
                    lastRenderTime = now;

                    const canvas = canvasRef.current;
                    const sp = spectrogramRef.current;
                    if (!canvas || !sp || !initialized) {
                        return;
                    }
                    if (needsUpdateRef.current) {
                        needsUpdateRef.current = false;
                        const bufData = bufferRef.current
                            .getData()
                            .filter((v): v is [number, number] => v[1] !== null);
                        sp.setData(bufData);
                        if (bufData.length > 0) {
                            const end = sp.getDuration();
                            timeRangeRef.current = [end - duration, end];
                        }
                    }

                    const { width, height } = sizeRef.current;
                    if (!width || !height) {
                        return;
                    }
                    sp.render({
                        timeRange: timeRangeRef.current,
                        canvas,
                        width,
                        height,
                        freqRange
                    });
                };
                rafId = requestAnimationFrame(renderLoop);
                return () => cancelAnimationFrame(rafId);
            }, [duration, freqRange, initialized, renderFPS]);

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
                <div className="relative h-full w-full">
                    <canvas ref={canvasRef} className="block h-full w-full" />

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
                                        <span>
                                            {t('components.DequeSpectrogram.settings.title')}
                                        </span>
                                        <button onClick={() => setShowSettings(false)}>
                                            <Icon path={mdiClose} size={0.6} />
                                        </button>
                                    </div>

                                    <>
                                        <div className="mb-1 flex justify-between">
                                            <span>
                                                {t('components.DequeSpectrogram.settings.min_db')}
                                            </span>
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
                                            <span>
                                                {t('components.DequeSpectrogram.settings.max_db')}
                                            </span>
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
    )
);
