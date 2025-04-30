import { mdiAccount, mdiChevronUp, mdiEarth, mdiKey, mdiShieldCheck } from '@mdi/js';
import Icon from '@mdi/react';
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
        encrypt_key: '',
        captcha_id: '',
        captcha_img: ''
    });
    const { t } = useTranslation();
    const getPreAuthData = useCallback(
        async (notify: boolean) => {
            setPreAuthData({ encrypt_key: '', captcha_id: '', captcha_img: '' });
            const requestFn = (throwError: boolean) => {
                return ApiClient.request<{
                    ttl: number;
                    encrypt_key: string;
                    captcha_id: string;
                    captcha_img: string;
                }>({
                    url: getRestfulApiUrl('/auth'),
                    method: 'post',
                    ignoreErrors: !throwError,
                    data: { action: 'preauth', nonce: '', credential: '' }
                });
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
                const { ttl, encrypt_key, captcha_id, captcha_img } = res.data;
                setPreAuthData({
                    encrypt_key,
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

    const { setCredential } = useCredentialStore();
    const handleLoginSubmit = async (username: string, password: string, captcha: string) => {
        // Load public key and encrypt credential
        const publicKey = util.decode64(preAuthData.encrypt_key);
        const credential = pki.publicKeyFromPem(publicKey).encrypt(
            JSON.stringify({
                username,
                password,
                timestamp: Date.now(),
                captcha_solution: captcha,
                captcha_id: preAuthData.captcha_id
            }),
            'RSA-OAEP'
        );

        const res = await ApiClient.request<{ token: string; life_time: number }>({
            url: getRestfulApiUrl('/auth'),
            method: 'post',
            data: {
                action: 'login',
                credential: util.encode64(credential),
                nonce: md.sha1.create().update(publicKey).digest().toHex()
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
        <div className="animate-fade animate-duration-500 animate-delay-300 flex min-h-screen flex-col bg-gradient-to-br from-purple-500 to-blue-500 p-20 px-4">
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
                                t('Login.signin.error'),
                                false
                            );
                        } catch {
                            getPreAuthData(false);
                        } finally {
                            setSubmitting(false);
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
                                <div className="flex justify-between space-x-2">
                                    <Field
                                        required
                                        autoComplete="off"
                                        id="captcha"
                                        className="mt-1 w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none disabled:bg-gray-100"
                                        type="text"
                                        name="captcha"
                                        disabled={!preAuthData.captcha_img.length}
                                        placeholder={t(
                                            preAuthData.captcha_img.length
                                                ? 'Login.captcha.placeholder'
                                                : 'Login.captcha.loading'
                                        )}
                                    />
                                    <div
                                        className="flex w-24 cursor-pointer items-center justify-center md:w-32"
                                        onClick={() => {
                                            getPreAuthData(true);
                                        }}
                                    >
                                        {preAuthData.captcha_img.length ? (
                                            <img
                                                className="text-gray-700"
                                                src={preAuthData.captcha_img}
                                                alt=""
                                            />
                                        ) : (
                                            <span className="loading loading-dots loading-sm bg-gray-500" />
                                        )}
                                    </div>
                                </div>
                            </div>

                            <button
                                className="btn mt-4 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                                type="submit"
                                disabled={isSubmitting || !preAuthData.captcha_img.length}
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
