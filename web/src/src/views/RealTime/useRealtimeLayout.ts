import { useCallback, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { ColorMapName } from 'spectrogram-js';

import { DEFAULT_SPECTROGRAM_COLOR_MAP } from '../../components/chart/spectrogramColorMaps';
import { RealTimeConstraints } from '../../config/constraints';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { useThrottleFnTrailing } from '../../helpers/utils/useThrottleFnTrailing';
import { type LayoutConfig, useLayoutStore } from '../../stores/layout';
import type { useRealtimeStream } from './useRealtimeStream';

export const useRealtimeLayout = (
    activeChannels: ReturnType<typeof useRealtimeStream>['activeChannels']
) => {
    const { t } = useTranslation();
    const { config, locks, toggleLock, setLayoutConfig, resetLayoutConfig } = useLayoutStore();
    const layoutCacheRef = useRef(new Map<string, LayoutConfig>());
    const getInitialLayout = useCallback(
        (id: string, index: number) => {
            const saved = config[id];
            if (saved?.position && saved?.size && saved?.spectrogram) {
                return saved;
            }

            const wide = document.documentElement.clientWidth > 768;
            const cacheKey = wide ? `${id}:w` : `${id}:n`;
            const cached = layoutCacheRef.current.get(cacheKey);
            if (cached) {
                return cached;
            }

            let x = 20;
            let y = 50;
            for (let i = 0; i < id.length; i++) {
                x = (x + id.charCodeAt(i) * (i + 1)) % 450;
                y = (y + id.charCodeAt(i) * (i + 2)) % 450;
            }
            x = Math.max(20, Math.min(x, 500));
            y = Math.max(50, Math.min(y, 500));

            const layout: LayoutConfig = {
                position: { x, y },
                size: {
                    width: wide ? RealTimeConstraints.minWidth * 2 : RealTimeConstraints.minWidth,
                    height: wide ? RealTimeConstraints.minWidth * 2 : RealTimeConstraints.minHeight
                },
                spectrogram: {
                    ...RealTimeConstraints.getDynamicDB(index),
                    colorMap: DEFAULT_SPECTROGRAM_COLOR_MAP
                }
            };
            layoutCacheRef.current.set(cacheKey, layout);
            return layout;
        },
        [config]
    );

    const handleResetLayout = useCallback(() => {
        sendUserConfirm(t('views.RealTime.reset_layout.confirm_message'), {
            title: t('views.RealTime.reset_layout.confirm_title'),
            cancelBtnText: t('views.RealTime.reset_layout.cancel_button'),
            confirmBtnText: t('views.RealTime.reset_layout.confirm_button'),
            onConfirmed: () => {
                Object.values(activeChannels).forEach(({ id }) => resetLayoutConfig(id));
            }
        });
    }, [t, activeChannels, resetLayoutConfig]);

    const handleDragStop = useCallback(
        (channel: string, index: number, x: number, y: number) => {
            setLayoutConfig(channel, { ...getInitialLayout(channel, index), position: { x, y } });
        },
        [getInitialLayout, setLayoutConfig]
    );

    const handleResizeStop = useCallback(
        (channel: string, index: number, width: number, height: number) => {
            setLayoutConfig(channel, {
                ...getInitialLayout(channel, index),
                size: { width, height }
            });
        },
        [getInitialLayout, setLayoutConfig]
    );

    const handleSpectrogramUpdate = useThrottleFnTrailing(
        useCallback(
            (
                channel: string,
                index: number,
                minDB: number,
                maxDB: number,
                colorMap: ColorMapName
            ) => {
                setLayoutConfig(channel, {
                    ...getInitialLayout(channel, index),
                    spectrogram: { maxDB, minDB, colorMap }
                });
            },
            [getInitialLayout, setLayoutConfig]
        ),
        500
    );

    return {
        locks,
        toggleLock,
        getInitialLayout,
        handleResetLayout,
        handleDragStop,
        handleResizeStop,
        handleSpectrogramUpdate
    };
};
