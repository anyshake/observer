import { mdiClose, mdiCog, mdiMagnifyMinus, mdiMagnifyPlus } from '@mdi/js';
import Icon from '@mdi/react';
import { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { ColorMapName, FFTExecutor } from 'spectrogram-js';
import { Spectrogram as SpectrogramCore } from 'spectrogram-js';

import { getSpectrogramAxis } from '../../styles/spectrogramTheme';
import { useThemeColors } from '../../styles/useThemeColors';
import { DEFAULT_SPECTROGRAM_COLOR_MAP, SPECTROGRAM_COLOR_MAPS } from './spectrogramColorMaps';

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
    readonly colorMap?: ColorMapName;
    readonly fftExecutor?: FFTExecutor;
    readonly renderFPS?: number;
    readonly paused?: boolean;
    readonly zoomStep?: number;
    readonly onSpectrogramUpdate?: (minDB: number, maxDB: number, colorMap: ColorMapName) => void;
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
        colorMap = DEFAULT_SPECTROGRAM_COLOR_MAP,
        fftExecutor,
        renderFPS = 2,
        paused = false,
        zoomStep = 0.2,
        onSpectrogramUpdate
    }: ISpectrogram) => {
        const { t } = useTranslation();
        const colors = useThemeColors();
        const axis = useMemo(() => getSpectrogramAxis(colors), [colors]);

        const [showSettings, setShowSettings] = useState(false);
        const [minDBState, setMinDBState] = useState(minDB);
        const [maxDBState, setMaxDBState] = useState(maxDB);
        const [zoomDuration, setZoomDuration] = useState(duration);
        const [colormap, setColormap] = useState<ColorMapName>(colorMap);

        const [initializedSpectrogram, setInitializedSpectrogram] =
            useState<SpectrogramCore | null>(null);
        const spectrogramRef = useRef<SpectrogramCore | null>(null);
        useEffect(() => {
            setInitializedSpectrogram(null);
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
                if (spectrogramRef.current !== sp) {
                    return;
                }
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

        const canvasRef = useRef<HTMLCanvasElement>(null);
        const pausedRef = useRef(paused);
        const appliedDataRef = useRef<ISpectrogram['data'] | null>(null);
        const appliedCoreRef = useRef<SpectrogramCore | null>(null);
        pausedRef.current = paused;
        const sizeRef = useRef({ width: 1, height: 1 });
        useEffect(() => {
            const canvas = canvasRef.current;
            if (!canvas || !canvas.parentElement) {
                return;
            }

            const ro = new ResizeObserver(([entry]) => {
                const { width, height } = entry.contentRect;
                sizeRef.current.width = Math.max(0, Math.floor(width));
                sizeRef.current.height = Math.max(0, Math.floor(height));
            });

            ro.observe(canvas.parentElement);
            return () => ro.disconnect();
        }, []);

        const [timePercent, setTimePercent] = useState(0.5);
        const dataDuration = useMemo<number | null>(() => {
            if (data.length === 0) {
                return null;
            }
            return (data[data.length - 1][0] - data[0][0]) / 1000;
        }, [data]);

        const renderSpectrogram = useCallback(
            (force = false) => {
                if (pausedRef.current && !force) {
                    return;
                }

                const canvas = canvasRef.current;
                const sp = spectrogramRef.current;
                if (!canvas || !sp || initializedSpectrogram !== sp) {
                    return;
                }

                const { width, height } = sizeRef.current;
                if (width <= 0 || height <= 0) {
                    return;
                }

                const total = dataDuration ?? 0;
                const window = zoomDuration;
                const maxStart = Math.max(0, total - window);
                const start = maxStart * timePercent;
                const end = start + window;
                sp.render({
                    canvas,
                    width,
                    height,
                    timeRange: [start, end],
                    freqRange,
                    xAxis: axis,
                    yAxis: axis
                });
            },
            [axis, dataDuration, freqRange, initializedSpectrogram, timePercent, zoomDuration]
        );

        useEffect(() => {
            renderSpectrogram(true);
        }, [renderSpectrogram]);

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

            if (
                appliedCoreRef.current !== initializedSpectrogram ||
                appliedDataRef.current !== data
            ) {
                initializedSpectrogram.setData(data);
                appliedCoreRef.current = initializedSpectrogram;
                appliedDataRef.current = data;
            }
            requestSpectrogramRender();
        }, [data, initializedSpectrogram, paused, requestSpectrogramRender]);

        useEffect(() => {
            if (
                paused ||
                !initializedSpectrogram ||
                spectrogramRef.current !== initializedSpectrogram
            ) {
                return;
            }

            initializedSpectrogram.setColormap(colormap);
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

        const handleZoomIn = useCallback(() => {
            if (dataDuration) {
                setZoomDuration((prev) => {
                    const next = prev * (1 - zoomStep);
                    return Math.max(next, dataDuration * 0.01);
                });
            }
        }, [dataDuration, zoomStep]);

        const handleZoomOut = useCallback(() => {
            if (dataDuration) {
                setZoomDuration((prev) => {
                    const next = prev * (1 + zoomStep);
                    return Math.min(next, dataDuration);
                });
            }
        }, [dataDuration, zoomStep]);

        return (
            <div className="app-spectrogram">
                <canvas className="app-block" ref={canvasRef} />

                <div className="app-spectrogram-toolbar">
                    {title && <div className="app-spectrogram-title">{title}</div>}

                    <button
                        className="app-spectrogram-button"
                        onClick={() => setShowSettings((v) => !v)}
                    >
                        <Icon path={mdiCog} size={0.7} />
                    </button>

                    <button className="app-spectrogram-mode" onClick={handleToggleColormap}>
                        {colormap}
                    </button>

                    <button className="app-spectrogram-button" onClick={handleZoomIn}>
                        <Icon path={mdiMagnifyPlus} size={0.7} />
                    </button>

                    <button className="app-spectrogram-button" onClick={handleZoomOut}>
                        <Icon path={mdiMagnifyMinus} size={0.7} />
                    </button>

                    <div className="app-spectrogram-speed">
                        <input
                            type="range"
                            min={0}
                            max={1}
                            step={0.01}
                            value={timePercent}
                            onChange={({ target }) => setTimePercent(Number(target.value))}
                            onMouseDown={(e) => e.stopPropagation()}
                            onTouchStart={(e) => e.stopPropagation()}
                            className="range range-xs app-spectrogram-speed-range"
                        />
                        <span className="app-spectrogram-speed-value">
                            {(timePercent * 100).toFixed(0)}%
                        </span>
                    </div>
                </div>

                {showSettings && (
                    <div
                        className="app-spectrogram-settings-overlay"
                        onClick={() => setShowSettings(false)}
                    >
                        <div className="app-spectrogram-settings-scroll">
                            <div
                                className="app-spectrogram-settings"
                                onClick={(e) => e.stopPropagation()}
                            >
                                <div className="app-spectrogram-settings-title">
                                    <span>{t('components.Spectrogram.settings.title')}</span>
                                    <button onClick={() => setShowSettings(false)}>
                                        <Icon path={mdiClose} size={0.6} />
                                    </button>
                                </div>

                                <>
                                    <div className="app-spectrogram-setting-label">
                                        <span>{t('components.Spectrogram.settings.min_db')}</span>
                                        <span>{minDBState}</span>
                                    </div>
                                    <input
                                        type="range"
                                        min={-250}
                                        max={250}
                                        step={2}
                                        value={minDBState}
                                        className="range range-info range-sm app-full-width"
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
                                    <div className="app-spectrogram-setting-label">
                                        <span>{t('components.Spectrogram.settings.max_db')}</span>
                                        <span>{maxDBState}</span>
                                    </div>
                                    <input
                                        type="range"
                                        min={-250}
                                        max={250}
                                        step={2}
                                        value={maxDBState}
                                        className="range range-info range-sm app-full-width"
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
);
