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

export const getChartTitle = (title: string | undefined, colors: ThemeColors) =>
    title
        ? {
              text: title,
              left: 'left',
              textStyle: { color: colors.label, fontSize: 13, fontWeight: 'bold' },
              backgroundColor: colors.titleBackground,
              padding: [6, 12],
              borderRadius: 3
          }
        : {};
