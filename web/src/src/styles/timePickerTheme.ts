import { createTheme } from '@mui/material/styles';

import { ThemeMode } from '../stores/theme';

// MUI's palette API needs resolved color values; CSS owns the palette definitions.
export const createAppTheme = (mode: ThemeMode, locale: object = {}) => {
    const styles = getComputedStyle(document.documentElement);
    const color = (name: string) => styles.getPropertyValue(`--app-${name}`).trim();

    return createTheme(
        {
            palette: {
                mode,
                primary: { main: color('primary'), contrastText: color('primary-content') },
                secondary: { main: color('picker-secondary') },
                background: { default: color('page'), paper: color('surface') },
                text: { primary: color('text'), secondary: color('muted') },
                divider: color('border')
            }
        },
        locale
    );
};
