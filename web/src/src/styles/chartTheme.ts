import type { ThemeColors } from './useThemeColors';

export const getChartAxis = (colors: ThemeColors) => ({
    axisLine: { lineStyle: { color: colors.axis } },
    axisTick: { lineStyle: { color: colors.axis } },
    axisLabel: {
        color: colors.label,
        rich: { normal: { color: colors.label } }
    },
    splitLine: { lineStyle: { color: colors.grid } }
});
