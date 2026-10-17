import type { AxisOptions } from 'spectrogram-js';

import type { ThemeColors } from './useThemeColors';

export const getSpectrogramAxis = (colors: ThemeColors): AxisOptions => ({
    color: colors.axis,
    labelStyle: {
        color: colors.label,
        // Keep the embedded font and metrics used by spectrogram-js 1.1.2.
        fontFamily: 'AxisFont',
        fontSize: 11,
        fontWeight: 'normal'
    }
});
