import { useMemo } from 'react';

import { selectThemeMode, ThemeMode, useThemeStore } from '../stores/theme';

const readThemeColors = (mode: ThemeMode) => {
    const styles = getComputedStyle(document.documentElement);
    const color = (name: string) => styles.getPropertyValue(`--app-${name}`).trim();

    return {
        mode,
        primary: color('primary'),
        primaryContent: color('primary-content'),
        secondary: color('picker-secondary'),
        page: color('page'),
        surface: color('surface'),
        text: color('text'),
        muted: color('muted'),
        border: color('border'),
        line: color('chart-line'),
        label: color('chart-label'),
        axis: color('chart-axis'),
        grid: color('chart-grid'),
        titleBackground: color('chart-title'),
        map: {
            background: color('map-background'),
            water: color('map-water'),
            land: color('map-land'),
            ice: color('map-ice'),
            iceShelf: color('map-ice-shelf'),
            playas: color('map-playas'),
            reef: color('map-reef'),
            island: color('map-island'),
            river: color('map-river'),
            coast: color('map-coast'),
            boundary: color('map-boundary'),
            city: color('map-city'),
            label: color('map-label'),
            halo: color('map-halo')
        }
    };
};

export type ThemeColors = ReturnType<typeof readThemeColors>;

export const useThemeColors = () => {
    const theme = useThemeStore(selectThemeMode);
    return useMemo(() => readThemeColors(theme), [theme]);
};
