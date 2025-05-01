import { useCallback, useEffect, useMemo, useState } from 'react';
import { Toaster } from 'react-hot-toast';
import { useTranslation } from 'react-i18next';
import { useRegisterSW } from 'virtual:pwa-register/react';

import { globalConfig } from './config/global';
import i18n, { localeConfig } from './config/locale';
import { Entry } from './Entry';
import { sendUserConfirm } from './helpers/alert/sendUserConfirm';
import { getRestfulApiUrl } from './helpers/app/getRestfulApiUrl';
import { hideLoaderAnimation } from './helpers/app/hideLoaderAnimation';
import { getCurrentLocale } from './helpers/locale/getCurrentLocale';
import { setUserLocale } from './helpers/locale/setUserLocale';
import { ApiClient } from './helpers/request/ApiClient';
import { Login } from './Login';
import { useCredentialStore } from './stores/credential';

const App = () => {
    const { t } = useTranslation();

    const {
        needRefresh: [needRefreshApp, setNeedRefreshApp],
        updateServiceWorker
    } = useRegisterSW({
        onRegistered(r) {
            if (r) {
                setInterval(() => {
                    r.update();
                }, globalConfig.update);
            }
        }
    });
    useEffect(() => {
        if (needRefreshApp) {
            sendUserConfirm(t('App.app_update.confirm_message'), {
                title: t('App.app_update.confirm_title'),
                cancelBtnText: t('App.app_update.cancel_button'),
                confirmBtnText: t('App.app_update.confirm_button'),
                onConfirmed: () => updateServiceWorker(),
                onCancelled: () => setNeedRefreshApp(false)
            });
        }
    }, [t, needRefreshApp, setNeedRefreshApp, updateServiceWorker]);

    const [hasLoggedIn, setHasLoggedIn] = useState(false);
    const { needRefresh, credential, setCredential } = useCredentialStore();
    const getUserLoginStatus = useCallback(async () => {
        const { error, code } = await ApiClient.request({
            url: getRestfulApiUrl('/auth'),
            method: 'get',
            ignoreErrors: true
        });
        if (!error && code === 200) {
            setHasLoggedIn(true);
        }
        hideLoaderAnimation();
    }, []);
    useEffect(() => {
        const { token, lifeTime } = credential;
        if (token.length > 0 && lifeTime > 0) {
            getUserLoginStatus();
        } else {
            setHasLoggedIn(false);
        }
    }, [credential, getUserLoginStatus]);

    const refreshToken = useCallback(async () => {
        const { data } = await ApiClient.request<{ token: string; life_time: number }>({
            url: getRestfulApiUrl('/auth'),
            method: 'post',
            data: { action: 'refresh' },
            ignoreErrors: true
        });
        if (data && data.token.length && data.life_time > 0) {
            setCredential(data.token, data.life_time);
        }
    }, [setCredential]);
    useEffect(() => {
        if (needRefresh) {
            refreshToken();
        }
    }, [needRefresh, refreshToken]);

    const localeMap = useMemo(
        () =>
            Object.entries(localeConfig.resources).reduce(
                (acc, [key, value]) => {
                    acc[key] = value.label;
                    return acc;
                },
                {} as Record<string, string>
            ),
        []
    );
    const [currentLocale, setCurrentLocale] = useState(localeConfig.fallback);
    const handleSwitchLocale = useCallback(async (newLocale: string) => {
        await setUserLocale(i18n, newLocale);
    }, []);
    useEffect(() => {
        (async () => {
            setCurrentLocale(await getCurrentLocale(i18n));
        })();
    }, [t]);

    return (
        <div>
            {hasLoggedIn ? (
                <Entry
                    currentLocale={currentLocale}
                    locales={localeMap}
                    onSwitchLocale={handleSwitchLocale}
                />
            ) : (
                <Login
                    currentLocale={currentLocale}
                    locales={localeMap}
                    onSwitchLocale={handleSwitchLocale}
                />
            )}
            <Toaster />
        </div>
    );
};

export default App;
