import { mdiThemeLightDark, mdiWeatherNight, mdiWhiteBalanceSunny } from '@mdi/js';
import Icon from '@mdi/react';
import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';

import { useThemeStore } from '../../stores/theme';

const themeOptions = [
    { value: 'light', icon: mdiWhiteBalanceSunny },
    { value: 'dark', icon: mdiWeatherNight },
    { value: 'system', icon: mdiThemeLightDark }
] as const;

export const ThemeToggle = () => {
    const { t } = useTranslation();
    const theme = useThemeStore((state) => state.theme);
    const setTheme = useThemeStore((state) => state.setTheme);
    const detailsRef = useRef<HTMLDetailsElement>(null);
    const label = t('components.ThemeToggle.label');
    const currentLabel = `${label}: ${t(`components.ThemeToggle.${theme}`)}`;

    useEffect(() => {
        const closeOnOutsideClick = (event: PointerEvent) => {
            const details = detailsRef.current;
            if (details && !details.contains(event.target as Node)) {
                details.open = false;
            }
        };
        const closeOnEscape = (event: KeyboardEvent) => {
            const details = detailsRef.current;
            if (event.key === 'Escape' && details?.open) {
                details.open = false;
                details.querySelector('summary')?.focus();
            }
        };

        document.addEventListener('pointerdown', closeOnOutsideClick);
        document.addEventListener('keydown', closeOnEscape);
        return () => {
            document.removeEventListener('pointerdown', closeOnOutsideClick);
            document.removeEventListener('keydown', closeOnEscape);
        };
    }, []);

    return (
        <details ref={detailsRef} className="dropdown dropdown-end">
            <summary
                className="btn btn-sm btn-ghost app-theme-toggle"
                aria-label={currentLabel}
                title={currentLabel}
            >
                <Icon
                    className="app-icon"
                    path={themeOptions.find(({ value }) => value === theme)!.icon}
                    size={0.8}
                />
            </summary>
            <ul
                className="dropdown-content menu rounded-box app-dropdown-menu app-theme-options"
                aria-label={label}
            >
                {themeOptions.map(({ value, icon }) => (
                    <li key={value}>
                        <button
                            type="button"
                            aria-pressed={theme === value}
                            onClick={() => setTheme(value)}
                        >
                            <Icon className="app-icon" path={icon} size={0.7} />
                            <span className="app-grow">{t(`components.ThemeToggle.${value}`)}</span>
                        </button>
                    </li>
                ))}
            </ul>
        </details>
    );
};
