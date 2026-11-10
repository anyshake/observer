import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { RetentionConstraints } from '../config/constraints';
import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';
import { normalizeRetention } from '../helpers/utils/normalizeRetention';

const retentionStoreKey = 'retention-store';

interface IRetention {
    readonly retention: number;
    readonly setRetention: (retention: number) => void;
}

export const useRetentionStore = create<IRetention>()(
    persist(
        (set) => ({
            retention: getLocalStorage(retentionStoreKey, RetentionConstraints.defaultValue),
            setRetention: (retention: number) => {
                set({ retention: normalizeRetention(retention) });
            }
        }),
        {
            name: retentionStoreKey,
            partialize: ({ retention }) => ({ retention }),
            merge: (persisted, current) => ({
                ...current,
                retention: normalizeRetention(
                    (persisted as { retention?: unknown } | null)?.retention ?? current.retention
                )
            }),
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            }
        }
    )
);
