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

import { ThemeToggle } from './components/ui/ThemeToggle';
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
    const [preAuthExpiresAt, setPreAuthExpiresAt] = useState(0);
    const [captchaSecondsRemaining, setCaptchaSecondsRemaining] = useState(0);
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
            setPreAuthExpiresAt(0);
            setCaptchaSecondsRemaining(0);
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
                setPreAuthExpiresAt(Date.now() + ttl);
                setCaptchaSecondsRemaining(Math.ceil(ttl / 1000));
            }
        },
        [t]
    );
    useEffect(() => {
        getPreAuthData(false);
    }, [getPreAuthData]);

    useEffect(() => {
        if (!preAuthExpiresAt) {
            return;
        }

        const updateCountdown = () => {
            setCaptchaSecondsRemaining(
                Math.max(0, Math.ceil((preAuthExpiresAt - Date.now()) / 1000))
            );
        };
        updateCountdown();
        const interval = window.setInterval(updateCountdown, 1000);
        // Each new captcha gets its own deadline, including manual refreshes.
        const timeout = window.setTimeout(
            () => getPreAuthData(true),
            Math.max(0, preAuthExpiresAt - Date.now())
        );
        return () => {
            window.clearInterval(interval);
            window.clearTimeout(timeout);
        };
    }, [preAuthExpiresAt, getPreAuthData]);

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
        <div className="app-login">
            <main className="app-login-main">
                <div className="app-login-panel">
                    <section className="app-login-artwork">
                        <img src={loginBackground} alt="" className="app-login-background" />
                    </section>

                    <section className="app-login-content">
                        <div className="app-login-heading">
                            <div className="app-login-brand-row">
                                <img src={globalConfig.logo} alt="" className="app-login-logo" />

                                <div className="app-login-actions">
                                    <ThemeToggle />
                                    <div className="dropdown dropdown-end">
                                        <button
                                            type="button"
                                            tabIndex={0}
                                            className="btn btn-sm btn-ghost app-login-locale-button"
                                        >
                                            <Icon className="app-icon" path={mdiEarth} size={0.8} />
                                            <span>{locales[currentLocale]}</span>
                                            <Icon
                                                className="app-icon"
                                                path={mdiChevronDown}
                                                size={0.6}
                                            />
                                        </button>
                                        <ul
                                            tabIndex={0}
                                            className="menu dropdown-content rounded-box app-dropdown-menu app-login-locales"
                                        >
                                            {Object.entries(locales).map(([key, value]) => (
                                                <li
                                                    className={`app-text-body ${key === currentLocale ? 'app-text-bold' : ''}`}
                                                    key={key}
                                                >
                                                    <button
                                                        type="button"
                                                        aria-pressed={key === currentLocale}
                                                        onClick={() => onSwitchLocale(key)}
                                                    >
                                                        {value}
                                                    </button>
                                                </li>
                                            ))}
                                        </ul>
                                    </div>
                                </div>
                            </div>
                            <hr className="app-login-divider" />
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
                                <Form className="app-stack-lg">
                                    <div>
                                        <label htmlFor="username" className="app-login-label">
                                            <Icon
                                                className="app-icon-leading"
                                                path={mdiAccount}
                                                size={0.8}
                                            />
                                            {t('Login.username.label')}
                                        </label>
                                        <Field
                                            required
                                            id="username"
                                            autoComplete="username"
                                            className="input app-input app-login-input"
                                            type="text"
                                            name="username"
                                            placeholder={t('Login.username.placeholder')}
                                        />
                                    </div>

                                    <div>
                                        <label htmlFor="password" className="app-login-label">
                                            <Icon
                                                className="app-icon-leading"
                                                path={mdiKey}
                                                size={0.8}
                                            />
                                            {t('Login.password.label')}
                                        </label>
                                        <Field
                                            required
                                            id="password"
                                            autoComplete="current-password"
                                            className="input app-input app-login-input"
                                            type="password"
                                            name="password"
                                            placeholder={t('Login.password.placeholder')}
                                        />
                                    </div>

                                    <div>
                                        <label htmlFor="captcha" className="app-login-label">
                                            <Icon
                                                className="app-icon-leading"
                                                path={mdiShieldCheck}
                                                size={0.8}
                                            />
                                            {t('Login.captcha.label')}
                                            {preAuthData.captcha_img.length > 0 &&
                                                captchaSecondsRemaining > 0 && (
                                                    <span className="app-login-captcha-countdown">
                                                        {t('Login.captcha.countdown', {
                                                            seconds: captchaSecondsRemaining
                                                        })}
                                                    </span>
                                                )}
                                        </label>
                                        <div className="app-login-captcha">
                                            <Field
                                                required
                                                autoComplete="off"
                                                id="captcha"
                                                className="input app-input app-login-captcha-input"
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
                                                className="btn app-login-captcha-button"
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
                                                        className="app-login-captcha-image"
                                                        src={preAuthData.captcha_img}
                                                        alt={t('Login.captcha.label')}
                                                    />
                                                ) : preAuthData.error ? (
                                                    <Icon
                                                        className="app-login-captcha-error"
                                                        path={mdiRefreshCircle}
                                                    />
                                                ) : (
                                                    <span className="loading loading-dots loading-sm app-login-captcha-loading" />
                                                )}
                                            </button>
                                        </div>
                                    </div>

                                    <button
                                        className="btn app-button-primary app-button-login"
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
