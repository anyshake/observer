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
import type { ColorMapName, FFTExecutor } from 'spectrogram-js';
import { Spectrogram as SpectrogramCore } from 'spectrogram-js';

import TimeSeriesBuffer from '../../helpers/storage/TimeSeriesBuffer';
import { DEFAULT_SPECTROGRAM_COLOR_MAP, SPECTROGRAM_COLOR_MAPS } from './spectrogramColorMaps';

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
    readonly colorMap?: ColorMapName;
    readonly fftExecutor?: FFTExecutor;
    readonly renderFPS?: number;
    readonly paused?: boolean;
    readonly onSpectrogramUpdate?: (minDB: number, maxDB: number, colorMap: ColorMapName) => void;
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
                colorMap = DEFAULT_SPECTROGRAM_COLOR_MAP,
                fftExecutor,
                renderFPS = 2,
                paused = false,
                onSpectrogramUpdate
            },
            ref
        ) => {
            const { t } = useTranslation();

            const [showSettings, setShowSettings] = useState(false);
            const [minDBState, setMinDBState] = useState(minDB);
            const [maxDBState, setMaxDBState] = useState(maxDB);
            const [colormap, setColormap] = useState<ColorMapName>(colorMap);

            const bufferRef = useRef<TimeSeriesBuffer>(new TimeSeriesBuffer(duration));
            const needsUpdateRef = useRef(true);
            const pausedRef = useRef(paused);
            pausedRef.current = paused;
            const [initializedSpectrogram, setInitializedSpectrogram] =
                useState<SpectrogramCore | null>(null);
            const spectrogramRef = useRef<SpectrogramCore | null>(null);
            useEffect(() => {
                setInitializedSpectrogram(null);
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
                    if (spectrogramRef.current !== sp) {
                        return;
                    }
                    needsUpdateRef.current = true;
                    setInitializedSpectrogram(sp);
                });
                return () => {
                    sp.destroy();
                    if (spectrogramRef.current === sp) {
                        spectrogramRef.current = null;
                    }
                };
            }, [fftExecutor, maxDB, minDB, overlap, sampleRate, windowSize]);

            useEffect(() => {
                setColormap(colorMap);
            }, [colorMap]);

            useEffect(() => {
                bufferRef.current = new TimeSeriesBuffer(duration);
                needsUpdateRef.current = true;
            }, [duration]);

            const canvasRef = useRef<HTMLCanvasElement>(null);
            const sizeRef = useRef({ width: 0, height: 0 });
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
                        sizeRef.current.width = Math.max(0, Math.floor(width));
                        sizeRef.current.height = Math.max(0, Math.floor(height));
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
            const needsBitmapFollowupRef = useRef(false);
            const bitmapRenderRef = useRef<number | null>(null);
            const drawSpectrogram = useCallback(() => {
                const canvas = canvasRef.current;
                const sp = spectrogramRef.current;
                if (!canvas || !sp || initializedSpectrogram !== sp) {
                    return;
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
            }, [freqRange, initializedSpectrogram]);

            const renderSpectrogram = useCallback(() => {
                if (pausedRef.current) {
                    return;
                }

                const canvas = canvasRef.current;
                const sp = spectrogramRef.current;
                const { width, height } = sizeRef.current;
                if (!canvas || !sp || initializedSpectrogram !== sp || !width || !height) {
                    return;
                }

                let needsBitmapFollowup = needsBitmapFollowupRef.current;
                needsBitmapFollowupRef.current = false;

                if (needsUpdateRef.current) {
                    needsUpdateRef.current = false;
                    const bufData = bufferRef.current.getSamples();
                    sp.setData(bufData);
                    if (bufData.length > 0) {
                        const end = sp.getDuration();
                        timeRangeRef.current = [end - duration, end];
                    }
                    needsBitmapFollowup = true;
                }

                drawSpectrogram();
                if (needsBitmapFollowup) {
                    if (bitmapRenderRef.current !== null) {
                        cancelAnimationFrame(bitmapRenderRef.current);
                    }
                    bitmapRenderRef.current = requestAnimationFrame(() => {
                        drawSpectrogram();
                        bitmapRenderRef.current = requestAnimationFrame(() => {
                            bitmapRenderRef.current = null;
                            drawSpectrogram();
                        });
                    });
                }
            }, [drawSpectrogram, duration, initializedSpectrogram]);

            const pendingRenderRef = useRef<number | null>(null);
            const requestSpectrogramRender = useCallback(() => {
                if (pendingRenderRef.current !== null) {
                    return;
                }
                pendingRenderRef.current = requestAnimationFrame(() => {
                    pendingRenderRef.current = null;
                    renderSpectrogram();
                });
            }, [renderSpectrogram]);

            useEffect(
                () => () => {
                    if (pendingRenderRef.current !== null) {
                        cancelAnimationFrame(pendingRenderRef.current);
                        pendingRenderRef.current = null;
                    }
                    if (bitmapRenderRef.current !== null) {
                        cancelAnimationFrame(bitmapRenderRef.current);
                        bitmapRenderRef.current = null;
                    }
                },
                [requestSpectrogramRender]
            );

            useEffect(() => {
                if (!initializedSpectrogram || paused) {
                    if (pendingRenderRef.current !== null) {
                        cancelAnimationFrame(pendingRenderRef.current);
                        pendingRenderRef.current = null;
                    }
                    return;
                }

                const frameInterval = 1000 / renderFPS;
                const intervalId = window.setInterval(requestSpectrogramRender, frameInterval);
                requestSpectrogramRender();
                return () => window.clearInterval(intervalId);
            }, [initializedSpectrogram, paused, renderFPS, requestSpectrogramRender]);

            useEffect(() => {
                if (
                    paused ||
                    !initializedSpectrogram ||
                    spectrogramRef.current !== initializedSpectrogram
                ) {
                    return;
                }

                initializedSpectrogram.setColormap(colormap);
                needsBitmapFollowupRef.current = true;
                requestSpectrogramRender();
            }, [colormap, initializedSpectrogram, paused, requestSpectrogramRender]);

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
                    onSpectrogramUpdate?.(value, Math.max(value, maxDBState), colormap);
                },
                [colormap, maxDBState, onSpectrogramUpdate, spectrogramRef]
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
                    onSpectrogramUpdate?.(Math.min(value, minDBState), value, colormap);
                },
                [colormap, minDBState, onSpectrogramUpdate, spectrogramRef]
            );

            const handleToggleColormap = useCallback(() => {
                const currentIndex = SPECTROGRAM_COLOR_MAPS.indexOf(colormap);
                const nextIndex = (currentIndex + 1) % SPECTROGRAM_COLOR_MAPS.length;
                const next = SPECTROGRAM_COLOR_MAPS[nextIndex];
                setColormap(next);
                onSpectrogramUpdate?.(minDBState, maxDBState, next);
            }, [colormap, maxDBState, minDBState, onSpectrogramUpdate]);

            return (
                <div className="relative h-full w-full">
                    <canvas ref={canvasRef} className="block h-full w-full" />

                    <div className="absolute top-5 left-15 flex items-center gap-1">
                        {title && (
                            <div className="flex h-6 items-center rounded bg-black/50 px-3 text-sm font-bold text-white select-none">
                                {title}
                            </div>
                        )}

                        <button
                            className="flex size-6 flex-shrink-0 cursor-pointer items-center justify-center rounded bg-black/50 text-white opacity-50 transition-all hover:opacity-100"
                            onClick={() => setShowSettings((v) => !v)}
                        >
                            <Icon path={mdiCog} size={0.7} />
                        </button>

                        <button
                            className="flex h-6 min-w-12 flex-shrink-0 cursor-pointer items-center justify-center rounded bg-black/50 px-2 text-xs font-semibold text-white opacity-50 transition-all hover:opacity-100"
                            onClick={handleToggleColormap}
                        >
                            {colormap}
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
                                            onMouseUp={(e) => {
                                                handleApplyMinDB(minDBState);
                                                e.stopPropagation();
                                            }}
                                            onMouseDown={(e) => {
                                                handleApplyMinDB(minDBState);
                                                e.stopPropagation();
                                            }}
                                            onTouchStart={(e) => {
                                                handleApplyMinDB(minDBState);
                                                e.stopPropagation();
                                            }}
                                            onTouchEnd={(e) => {
                                                handleApplyMinDB(minDBState);
                                                e.stopPropagation();
                                            }}
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
                                            onMouseUp={(e) => {
                                                handleApplyMaxDB(maxDBState);
                                                e.stopPropagation();
                                            }}
                                            onMouseDown={(e) => {
                                                handleApplyMaxDB(maxDBState);
                                                e.stopPropagation();
                                            }}
                                            onTouchStart={(e) => {
                                                handleApplyMaxDB(maxDBState);
                                                e.stopPropagation();
                                            }}
                                            onTouchEnd={(e) => {
                                                handleApplyMaxDB(maxDBState);
                                                e.stopPropagation();
                                            }}
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
