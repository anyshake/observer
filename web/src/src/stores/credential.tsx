import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { getLocalStorage } from '../helpers/storage/getLocalStorage';
import { setLocalStorage } from '../helpers/storage/setLocalStorage';

const credentialStoreKey = 'credential-store';
const tokenRefreshThreshold = 0.8;

interface ICredential {
    readonly needRefresh: boolean;
    readonly clearCredential: () => void;
    readonly credential: {
        readonly lifeTime: number;
        readonly token: string;
        readonly createdAt: number;
    };
    readonly getCredential: () => {
        readonly lifeTime: number;
        readonly token: string;
    };
    readonly setCredential: (token: string, lifeTime: number) => void;
}

export const useCredentialStore = create<ICredential>()(
    persist(
        (set, get) => ({
            needRefresh: false,
            credential: getLocalStorage(credentialStoreKey, {
                token: '',
                lifeTime: 0,
                createdAt: 0
            }),
            clearCredential: () => {
                set({
                    credential: { token: '', lifeTime: 0, createdAt: 0 },
                    needRefresh: false
                });
            },
            setCredential: (token: string, lifeTime: number) => {
                set({
                    credential: { token, lifeTime, createdAt: Date.now() },
                    needRefresh: false
                });
            },
            getCredential: () => {
                const { credential } = get();
                const currentTime = Date.now();
                const shouldRefresh =
                    credential.lifeTime > 0 &&
                    currentTime - credential.createdAt >=
                        credential.lifeTime * tokenRefreshThreshold &&
                    currentTime - credential.createdAt < credential.lifeTime;
                if (shouldRefresh && !get().needRefresh) {
                    set({ needRefresh: true });
                }
                const { token, lifeTime } = get().credential;
                return { token, lifeTime };
            }
        }),
        {
            name: credentialStoreKey,
            storage: {
                getItem: (name) => getLocalStorage(name, null),
                removeItem: (name) => setLocalStorage(name, null),
                setItem: (name, value) => setLocalStorage(name, value)
            }
        }
    )
);
