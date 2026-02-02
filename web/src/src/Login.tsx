import {
    mdiAccount,
    mdiChevronUp,
    mdiEarth,
    mdiKey,
    mdiRefreshCircle,
    mdiShieldCheck
} from '@mdi/js';
import Icon from '@mdi/react';
import { Buffer } from 'buffer';
import { ErrorMessage, Field, Form, Formik } from 'formik';
import { md, pki, util } from 'node-forge';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { globalConfig } from './config/global';
import { localeConfig } from './config/locale';
import { sendPromiseAlert } from './helpers/alert/sendPromiseAlert';
import { getRestfulApiUrl } from './helpers/app/getRestfulApiUrl';
import { ApiClient } from './helpers/request/ApiClient';
import { useCredentialStore } from './stores/credential';

interface ILogin {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly locales: Record<string, string>;
    readonly onSwitchLocale: (newLocale: string) => void;
}

export const Login = ({ currentLocale, locales, onSwitchLocale }: ILogin) => {
    useEffect(() => {
        document.title = globalConfig.name[currentLocale];
    }, [currentLocale]);

    // State for pre-authentication data (captcha, encrypt key, etc.)
    const [preAuthTTL, setPreAuthTTL] = useState(0);
    const [preAuthData, setPreAuthData] = useState({
        public_key: '',
        captcha_id: '',
        captcha_img: '',
        error: false
    });
    const { t } = useTranslation();
    const getPreAuthData = useCallback(
        async (notify: boolean) => {
            setPreAuthData({ public_key: '', captcha_id: '', captcha_img: '', error: false });
            const requestFn = async (throwError: boolean) => {
                const result = await ApiClient.request<{
                    ttl: number;
                    public_key: string;
                    captcha_id: string;
                    captcha_img: string;
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
                const { ttl, public_key, captcha_id, captcha_img } = res.data;
                setPreAuthData({
                    error: false,
                    public_key,
                    captcha_id,
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
        const encrypt = async (
            secret: Uint8Array<ArrayBuffer>,
            data: Uint8Array<ArrayBuffer>,
            aad: Uint8Array<ArrayBuffer>
        ): Promise<string> => {
            const keyMaterial = await crypto.subtle.importKey(
                'raw',
                secret,
                { name: 'HKDF' },
                false,
                ['deriveKey']
            );
            const key = await crypto.subtle.deriveKey(
                {
                    name: 'HKDF',
                    hash: 'SHA-512',
                    salt: new Uint8Array([]),
                    info: new Uint8Array([])
                },
                keyMaterial,
                { name: 'AES-GCM', length: 256 },
                false,
                ['encrypt']
            );

            const iv = crypto.getRandomValues(new Uint8Array(12));
            const ciphertextBuffer = await crypto.subtle.encrypt(
                { name: 'AES-GCM', iv, additionalData: aad, tagLength: 128 },
                key,
                data
            );

            const ciphertext = new Uint8Array(ciphertextBuffer);
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
            Buffer.from(
                JSON.stringify({
                    timestamp: Date.now(), // currently not used in backend
                    username,
                    password,
                    captcha_solution: captcha,
                    captcha_id: preAuthData.captcha_id
                })
            ),
            sessionIdBuffer
        );

        const res = await ApiClient.request<{ token: string; life_time: number }>({
            url: getRestfulApiUrl('/auth'),
            method: 'post',
            data: {
                action: 'login',
                session: sessionId,
                nonce: encryptedNonce,
                secret: encryptedSecret,
                payload: encryptedPayload
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
        <div className="animate-fade animate-duration-500 animate-delay-300 flex min-h-screen flex-col bg-gradient-to-br from-purple-300 via-purple-200 to-purple-300 p-20 px-4">
            <div className="m-auto w-full max-w-md rounded-lg bg-white p-12 shadow-xl md:max-w-xl">
                <img
                    src={globalConfig.logo}
                    alt="Login"
                    className="mx-auto mb-6 size-24 md:size-32"
                />
                <Formik
                    enableReinitialize
                    initialValues={{ username: '', password: '', captcha: '' }}
                    onSubmit={async ({ username, password, captcha }, { setSubmitting }) => {
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
                        <Form className="space-y-4">
                            <div>
                                <label
                                    htmlFor="username"
                                    className="flex items-center text-gray-700"
                                >
                                    <Icon
                                        className="mr-2 flex-shrink-0"
                                        path={mdiAccount}
                                        size={0.8}
                                    />
                                    {t('Login.username.label')}
                                </label>
                                <ErrorMessage
                                    className="text-sm text-red-500"
                                    name="username"
                                    component="div"
                                />
                                <Field
                                    required
                                    id="username"
                                    className="mt-1 w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none disabled:bg-gray-100"
                                    type="text"
                                    name="username"
                                    placeholder={t('Login.username.placeholder')}
                                />
                            </div>

                            <div>
                                <label
                                    htmlFor="password"
                                    className="flex items-center text-gray-700"
                                >
                                    <Icon className="mr-2 flex-shrink-0" path={mdiKey} size={0.8} />
                                    {t('Login.password.label')}
                                </label>
                                <ErrorMessage
                                    className="text-sm text-red-500"
                                    name="password"
                                    component="div"
                                />
                                <Field
                                    required
                                    id="password"
                                    className="mt-1 w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none disabled:bg-gray-100"
                                    type="password"
                                    name="password"
                                    placeholder={t('Login.password.placeholder')}
                                />
                            </div>

                            <div>
                                <label
                                    htmlFor="captcha"
                                    className="flex items-center text-gray-700"
                                >
                                    <Icon
                                        className="mr-2 flex-shrink-0"
                                        path={mdiShieldCheck}
                                        size={0.8}
                                    />
                                    {t('Login.captcha.label')}
                                </label>
                                <ErrorMessage
                                    className="text-sm text-red-500"
                                    name="captcha"
                                    component="div"
                                />
                                <div className="mt-1 flex items-center justify-between space-x-2">
                                    <Field
                                        required
                                        autoComplete="off"
                                        id="captcha"
                                        className="w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none disabled:bg-gray-100"
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
                                    <div
                                        className="flex w-24 cursor-pointer items-center justify-center rounded-md border border-gray-300 py-2 transition-all hover:border-gray-400 md:w-32"
                                        onClick={() => {
                                            getPreAuthData(true);
                                        }}
                                    >
                                        {preAuthData.captcha_img.length ? (
                                            <img
                                                className="h-6"
                                                src={preAuthData.captcha_img}
                                                alt=""
                                            />
                                        ) : preAuthData.error ? (
                                            <Icon
                                                className="size-6 flex-shrink-0 text-red-400"
                                                path={mdiRefreshCircle}
                                            />
                                        ) : (
                                            <span className="loading loading-dots loading-sm size-6 bg-gray-500" />
                                        )}
                                    </div>
                                </div>
                            </div>

                            <button
                                className="btn mt-4 mb-8 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                                type="submit"
                                disabled={
                                    isSubmitting ||
                                    !preAuthData.captcha_img.length ||
                                    credential.token.length > 0
                                }
                            >
                                {t('Login.signin.button')}
                            </button>

                            <div className="flex items-center justify-center text-xs">
                                <div className="dropdown dropdown-top dropdown-end">
                                    <summary
                                        tabIndex={0}
                                        className="flex cursor-pointer items-center space-x-2 text-gray-500 select-none hover:text-gray-700"
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiEarth}
                                            size={0.6}
                                        />
                                        <span>{locales[currentLocale]}</span>
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiChevronUp}
                                            size={0.6}
                                        />
                                    </summary>
                                    <ul
                                        tabIndex={0}
                                        className="menu dropdown-content bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
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
                        </Form>
                    )}
                </Formik>
            </div>
        </div>
    );
};
