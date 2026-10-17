import { selectThemeMode, ThemeMode, useThemeStore } from './stores/theme';

const applyTheme = (theme: ThemeMode) => {
    document.documentElement.dataset.theme = theme;
    document
        .querySelector('meta[name="theme-color"]')
        ?.setAttribute(
            'content',
            getComputedStyle(document.documentElement).getPropertyValue('--app-page').trim()
        );
};

const systemTheme = window.matchMedia('(prefers-color-scheme: dark)');
const syncSystemTheme = () => {
    useThemeStore.setState({ systemTheme: systemTheme.matches ? 'dark' : 'light' });
};

syncSystemTheme();
applyTheme(selectThemeMode(useThemeStore.getState()));
const unsubscribe = useThemeStore.subscribe((state) => applyTheme(selectThemeMode(state)));
systemTheme.addEventListener('change', syncSystemTheme);

if (import.meta.hot) {
    import.meta.hot.dispose(() => {
        unsubscribe();
        systemTheme.removeEventListener('change', syncSystemTheme);
    });
}
