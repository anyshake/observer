import { ThemeProvider } from '@mui/material/styles';
import { ReactNode, useMemo } from 'react';

import { selectThemeMode, useThemeStore } from '../../stores/theme';
import { createAppTheme } from '../../styles/timePickerTheme';

export const AppThemeProvider = ({ children }: { readonly children: ReactNode }) => {
    const mode = useThemeStore(selectThemeMode);
    const theme = useMemo(() => createAppTheme(mode), [mode]);

    return <ThemeProvider theme={theme}>{children}</ThemeProvider>;
};
