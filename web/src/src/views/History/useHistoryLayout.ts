import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import type { ColorMapName } from 'spectrogram-js';

import { HistoryConstraints } from '../../config/constraints';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { useThrottleFnTrailing } from '../../helpers/utils/useThrottleFnTrailing';
import { useLayoutStore } from '../../stores/layout';
import type { useSeismicRecords } from './useSeismicRecords';

export const useHistoryLayout = (
    activeChannels: ReturnType<typeof useSeismicRecords>['activeChannels']
) => {
    const { t } = useTranslation();
    const { config, locks, toggleLock, setLayoutConfig, resetLayoutConfig } = useLayoutStore();
    const getInitialLayout = useCallback(
        (id: string, index: number) => {
            if (config[id]?.position && config[id]?.size && config[id]?.spectrogram) {
                return config[id];
            }

            let x = 20;
            let y = 50;
            for (let i = 0; i < id.length; i++) {
                x = (x + id.charCodeAt(i) * (i + 1)) % 450;
                y = (y + id.charCodeAt(i) * (i + 2)) % 450;
            }
            x = Math.max(20, Math.min(x, 500));
            y = Math.max(50, Math.min(y, 500));

            return {
                position: { x, y },
                size: {
                    width:
                        document.documentElement.clientWidth > 768
                            ? HistoryConstraints.minWidth * 2
                            : HistoryConstraints.minWidth,
                    height:
                        document.documentElement.clientWidth > 768
                            ? HistoryConstraints.minWidth * 2
                            : HistoryConstraints.minHeight
                },
                spectrogram: {
                    ...HistoryConstraints.getDynamicDB(index),
                    colorMap: 'jet' as const
                }
            };
        },
        [config]
    );

    const handleResetLayout = useCallback(() => {
        sendUserConfirm(t('views.History.reset_layout.confirm_message'), {
            title: t('views.History.reset_layout.confirm_title'),
            cancelBtnText: t('views.History.reset_layout.cancel_button'),
            confirmBtnText: t('views.History.reset_layout.confirm_button'),
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
