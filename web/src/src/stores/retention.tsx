import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';

const retentionStoreKey = 'retention-store';

interface IRetention {
    readonly retention: number;
    readonly setRetention: (retention: number) => void;
}

export const useRetentionStore = create<IRetention>()(
    persist(
        (set) => ({
            retention: getLocalStorage(retentionStoreKey, 120),
            setRetention: (retention: number) => {
                set({ retention });
            }
        }),
        {
            name: retentionStoreKey,
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            }
        }
    )
);
