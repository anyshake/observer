import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';

export type ThemeMode = 'light' | 'dark';
export type ThemePreference = ThemeMode | 'system';

interface ITheme {
    readonly theme: ThemePreference;
    readonly systemTheme: ThemeMode;
    readonly setTheme: (theme: ThemePreference) => void;
}

export const selectThemeMode = ({ theme, systemTheme }: ITheme): ThemeMode =>
    theme === 'system' ? systemTheme : theme;

export const useThemeStore = create<ITheme>()(
    persist(
        (set) => ({
            theme: 'system',
            systemTheme: window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light',
            setTheme: (theme) => set({ theme })
        }),
        {
            name: 'theme-store',
            partialize: ({ theme }) => ({ theme }),
            merge: (persisted, current) => {
                const theme = (persisted as { theme?: unknown } | null)?.theme;
                return {
                    ...current,
                    ...(theme === 'light' || theme === 'dark' || theme === 'system'
                        ? { theme }
                        : {})
                };
            },
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            }
        }
    )
);
