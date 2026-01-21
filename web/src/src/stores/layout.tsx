import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';

const layoutStoreKey = 'layout-store';

type LayoutConfig = {
    readonly position: { x: number; y: number };
    readonly size: { width: number; height: number };
    readonly spectrogram: { maxDB: number; minDB: number };
};

interface ILayout {
    locks: Record<string, boolean>;
    toggleLock: (id: string) => void;
    setLock: (id: string, locked: boolean) => void;
    config: Record<string, LayoutConfig>;
    resetLayoutConfig: (id?: string) => void;
    setLayoutConfig: (id: string, config: LayoutConfig) => void;
}

export const useLayoutStore = create<ILayout>()(
    persist(
        (set) => ({
            locks: getLocalStorage(layoutStoreKey, {}),
            setLock: (id: string, locked: boolean) => {
                set((state) => {
                    state.locks[id] = locked;
                    return state;
                });
            },
            toggleLock: (id: string) => {
                set((state) => {
                    if (!state.locks[id]) {
                        state.locks[id] = false;
                    }
                    state.locks[id] = !state.locks[id];
                    return { locks: state.locks };
                });
            },
            config: getLocalStorage(layoutStoreKey, {}),
            resetLayoutConfig: (id?: string) => {
                set((state) => {
                    if (id) {
                        delete state.config[id];
                    } else {
                        state.config = {};
                    }
                    return { config: state.config };
                });
            },
            setLayoutConfig: (id: string, config: LayoutConfig) => {
                set((state) => {
                    state.config[id] = config;
                    return { config: state.config };
                });
            }
        }),
        {
            name: layoutStoreKey,
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            }
        }
    )
);
