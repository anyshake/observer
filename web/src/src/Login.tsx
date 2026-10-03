import {
    mdiAccount,
    mdiChevronDown,
    mdiEarth,
    mdiKey,
    mdiRefreshCircle,
    mdiShieldCheck
} from '@mdi/js';
import Icon from '@mdi/react';
import { gcm } from '@noble/ciphers/aes.js';
import { hkdf } from '@noble/hashes/hkdf.js';
import { sha512 } from '@noble/hashes/sha2.js';
import { Buffer } from 'buffer';
import { Field, Form, Formik } from 'formik';
import { md, pki, util } from 'node-forge';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import loginBackground from '/login_bg.webp?url';

import { globalConfig } from './config/global';
import { localeConfig } from './config/locale';
import { sendPromiseAlert } from './helpers/alert/sendPromiseAlert';
import { getRestfulApiUrl } from './helpers/app/getRestfulApiUrl';
import { solvePoWChallenge } from './helpers/app/solvePoWChallenge';
import { ApiClient } from './helpers/request/ApiClient';
import { useCredentialStore } from './stores/credential';

interface ILogin {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly locales: Record<string, string>;
    readonly onSwitchLocale: (newLocale: string) => void;
}

export const Login = ({ currentLocale, locales, onSwitchLocale }: ILogin) => {
    const { t } = useTranslation();
    useEffect(() => {
        document.title = t(globalConfig.name);
    }, [currentLocale, t]);

    // State for pre-authentication data (captcha, encrypt key, etc.)
    const [preAuthTTL, setPreAuthTTL] = useState(0);
    const [preAuthData, setPreAuthData] = useState({
        public_key: '',
        captcha_id: '',
        challenge_id: '',
        challenge_seed: '',
        captcha_img: '',
        error: false
    });
    const getPreAuthData = useCallback(
        async (notify: boolean) => {
            setPreAuthData({
                public_key: '',
                captcha_id: '',
                challenge_id: '',
                challenge_seed: '',
                captcha_img: '',
                error: false
            });
            const requestFn = async (throwError: boolean) => {
                const result = await ApiClient.request<{
                    ttl: number;
                    public_key: string;
                    captcha_id: string;
                    captcha_img: string;
                    challenge_id: string;
                    challenge_seed: string;
                }>({
                    url: getRestfulApiUrl('/auth'),
                    method: 'post',
                    ignoreErrors: true,
                    data: { action: 'preauth', nonce: '', credential: '' }
                });
                if (result.error) {
                    setPreAuthData((preAuthData) => ({ ...preAuthData, error: true }));
                    if (throwError) {
                        throw new Error(result.message);
                    }
                }
                return result;
            };
            const res = (
                notify
                    ? await sendPromiseAlert(
                          requestFn(true),
                          t('Login.captcha.refreshing'),
                          t('Login.captcha.refresh_success'),
                          t('Login.captcha.refresh_error')
                      )
                    : await requestFn(false)
            )!;
            if (res?.data) {
                const { ttl, public_key, captcha_id, challenge_id, challenge_seed, captcha_img } =
                    res.data;
                setPreAuthData({
                    error: false,
                    public_key,
                    captcha_id,
                    challenge_id,
                    challenge_seed,
                    captcha_img: `data:image/png;base64,${captcha_img}`
                });
                setPreAuthTTL(ttl);
            }
        },
        [t]
    );
    useEffect(() => {
        // Refresh captcha and encrypt key after TTL
        if (preAuthTTL) {
            const interval = setInterval(() => {
                getPreAuthData(true);
            }, preAuthTTL);
            return () => clearInterval(interval);
        }
        getPreAuthData(false);
    }, [preAuthTTL, getPreAuthData]);

    const { setCredential, credential } = useCredentialStore();
    const handleLoginSubmit = async (username: string, password: string, captcha: string) => {
        const encrypt = async (secret: Uint8Array, data: Uint8Array, aad: Uint8Array) => {
            const key = hkdf(sha512, secret, new Uint8Array([]), new Uint8Array([]), 32);
            const iv = crypto.getRandomValues(new Uint8Array(12));
            const cipher = gcm(key, iv, aad);
            const ciphertext = cipher.encrypt(data);

            const payload = new Uint8Array(iv.length + ciphertext.length);
            payload.set(iv, 0);
            payload.set(ciphertext, iv.length);

            return Buffer.from(payload).toString('base64');
        };

        const publicKey = util.decode64(preAuthData.public_key);
        const sessionId = md.sha512.create().update(publicKey).digest().toHex();
        const sessionIdBuffer = Buffer.from(sessionId);

        const secret = crypto.getRandomValues(new Uint8Array(32));
        const encryptedSecret = util.encode64(
            pki
                .publicKeyFromPem(publicKey)
                .encrypt(Buffer.from(secret).toString('base64'), 'RSA-OAEP')
        );

        const encryptedNonce = await encrypt(
            secret,
            crypto.getRandomValues(new Uint8Array(16)),
            sessionIdBuffer
        );
        const encryptedPayload = await encrypt(
            secret,
            Buffer.from(JSON.stringify({ username, password })),
            sessionIdBuffer
        );

        const { nonce, hash } = await solvePoWChallenge(preAuthData.challenge_seed);
        const res = await ApiClient.request<{ token: string; life_time: number }>({
            url: getRestfulApiUrl('/auth'),
            method: 'post',
            data: {
                action: 'login',
                session: sessionId,
                nonce: encryptedNonce,
                secret: encryptedSecret,
                payload: encryptedPayload,
                captcha_val: captcha,
                captcha_id: preAuthData.captcha_id,
                challenge_id: preAuthData.challenge_id,
                challenge_solution: `${nonce}:${hash}`
            }
        });
        if (res.error) {
            throw new Error(res.message);
        }
        if (res.data) {
            const { token, life_time } = res.data!;
            setCredential(token, life_time);
        }
    };

    return (
        <div className="bg-base-200 animate-fade animate-duration-500 animate-delay-300 flex min-h-screen flex-col text-gray-700">
            <main className="flex flex-1 items-center justify-center px-4 py-10 sm:px-6 lg:py-14">
                <div className="bg-base-100 grid w-full max-w-4xl overflow-hidden rounded-xl shadow-xl lg:grid-cols-[0.9fr_1.1fr]">
                    <section className="relative hidden min-h-[590px] overflow-hidden bg-gray-800 lg:block">
                        <img
                            src={loginBackground}
                            alt=""
                            className="absolute -inset-3 h-[calc(100%+1.5rem)] w-[calc(100%+1.5rem)] scale-105 object-cover brightness-50"
                        />
                    </section>

                    <section className="p-7 sm:p-10 lg:p-12">
                        <div className="mb-8">
                            <div className="flex items-center justify-between gap-4">
                                <img
                                    src={globalConfig.logo}
                                    alt=""
                                    className="flex size-18 items-center gap-4"
                                />

                                <div className="dropdown dropdown-end">
                                    <button
                                        type="button"
                                        tabIndex={0}
                                        className="btn btn-sm btn-ghost gap-2 text-gray-500"
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiEarth}
                                            size={0.8}
                                        />
                                        <span>{locales[currentLocale]}</span>
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiChevronDown}
                                            size={0.6}
                                        />
                                    </button>
                                    <ul
                                        tabIndex={0}
                                        className="menu dropdown-content bg-base-100 rounded-box z-20 mt-2 w-52 p-2 shadow-md"
                                    >
                                        {Object.entries(locales).map(([key, value]) => (
                                            <li
                                                className={`text-gray-700 ${key === currentLocale ? 'font-bold' : ''}`}
                                                key={key}
                                            >
                                                <a onClick={() => onSwitchLocale(key)}>{value}</a>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            </div>
                            <hr className="mt-6 text-gray-200" />
                        </div>

                        <Formik
                            enableReinitialize
                            initialValues={{ username: '', password: '', captcha: '' }}
                            onSubmit={async (
                                { username, password, captcha },
                                { setSubmitting }
                            ) => {
                                try {
                                    await sendPromiseAlert(
                                        handleLoginSubmit(username, password, captcha),
                                        t('Login.signin.loading'),
                                        t('Login.signin.success'),
                                        () => {
                                            setSubmitting(false);
                                            return t('Login.signin.error');
                                        },
                                        false
                                    );
                                } catch {
                                    getPreAuthData(false);
                                }
                            }}
                        >
                            {({ isSubmitting }) => (
                                <Form className="space-y-5">
                                    <div>
                                        <label
                                            htmlFor="username"
                                            className="flex items-center text-sm font-medium text-gray-700"
                                        >
                                            <Icon
                                                className="mr-2 flex-shrink-0"
                                                path={mdiAccount}
                                                size={0.8}
                                            />
                                            {t('Login.username.label')}
                                        </label>
                                        <Field
                                            required
                                            id="username"
                                            autoComplete="username"
                                            className="input mt-2 w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                            type="text"
                                            name="username"
                                            placeholder={t('Login.username.placeholder')}
                                        />
                                    </div>

                                    <div>
                                        <label
                                            htmlFor="password"
                                            className="flex items-center text-sm font-medium text-gray-700"
                                        >
                                            <Icon
                                                className="mr-2 flex-shrink-0"
                                                path={mdiKey}
                                                size={0.8}
                                            />
                                            {t('Login.password.label')}
                                        </label>
                                        <Field
                                            required
                                            id="password"
                                            autoComplete="current-password"
                                            className="input mt-2 w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                            type="password"
                                            name="password"
                                            placeholder={t('Login.password.placeholder')}
                                        />
                                    </div>

                                    <div>
                                        <label
                                            htmlFor="captcha"
                                            className="flex items-center text-sm font-medium text-gray-700"
                                        >
                                            <Icon
                                                className="mr-2 flex-shrink-0"
                                                path={mdiShieldCheck}
                                                size={0.8}
                                            />
                                            {t('Login.captcha.label')}
                                        </label>
                                        <div className="mt-2 flex items-stretch justify-between space-x-2">
                                            <Field
                                                required
                                                autoComplete="off"
                                                id="captcha"
                                                className="input min-w-0 flex-1 border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                                type="text"
                                                name="captcha"
                                                disabled={!preAuthData.captcha_img.length}
                                                placeholder={t(
                                                    preAuthData.captcha_img.length
                                                        ? 'Login.captcha.placeholder'
                                                        : preAuthData.error
                                                          ? 'Login.captcha.error'
                                                          : 'Login.captcha.loading'
                                                )}
                                            />
                                            <button
                                                type="button"
                                                className="btn w-28 flex-shrink-0 border border-gray-300 shadow-sm transition-all hover:ring sm:w-32"
                                                disabled={
                                                    !preAuthData.captcha_img.length &&
                                                    !preAuthData.error
                                                }
                                                onClick={() => {
                                                    if (
                                                        preAuthData.captcha_img.length ||
                                                        preAuthData.error
                                                    ) {
                                                        getPreAuthData(true);
                                                    }
                                                }}
                                            >
                                                {preAuthData.captcha_img.length ? (
                                                    <img
                                                        className="h-6"
                                                        src={preAuthData.captcha_img}
                                                        alt={t('Login.captcha.label')}
                                                    />
                                                ) : preAuthData.error ? (
                                                    <Icon
                                                        className="size-6 flex-shrink-0 text-red-400"
                                                        path={mdiRefreshCircle}
                                                    />
                                                ) : (
                                                    <span className="loading loading-dots loading-sm size-6 bg-gray-500" />
                                                )}
                                            </button>
                                        </div>
                                    </div>

                                    <button
                                        className="btn mt-3 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                                        type="submit"
                                        disabled={
                                            isSubmitting ||
                                            !preAuthData.captcha_img.length ||
                                            credential.token.length > 0
                                        }
                                    >
                                        {t('Login.signin.button')}
                                    </button>
                                </Form>
                            )}
                        </Formik>
                    </section>
                </div>
            </main>
        </div>
    );
};
