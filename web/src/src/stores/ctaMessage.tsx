import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';

const ctaMessageStoreKey = 'cta-message-store';
const ctaReadFlagExpiration = 7 * 24 * 60 * 60 * 1000;

interface ICTAMessage {
    readonly status: {
        readonly read: boolean;
        readonly readAt: number;
    };
    readonly markAsRead: () => void;
}

export const useCtaMessageStore = create<ICTAMessage>()(
    persist(
        (set) => ({
            status: getLocalStorage(ctaMessageStoreKey, { read: false, readAt: 0 }),
            markAsRead: () => {
                set({ status: { read: true, readAt: Date.now() } });
            }
        }),
        {
            name: ctaMessageStoreKey,
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            },
            merge: (persisted, current) => {
                const persistedStatus = (persisted as ICTAMessage)?.status;
                let status = persistedStatus;

                if (
                    persistedStatus?.read &&
                    Date.now() - persistedStatus.readAt > ctaReadFlagExpiration
                ) {
                    status = { read: false, readAt: 0 };
                }

                return { ...current, status: status ?? { read: false, readAt: 0 } };
            }
        }
    )
);
