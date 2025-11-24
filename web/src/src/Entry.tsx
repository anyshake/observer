import { mdiEmailNewsletter } from '@mdi/js';
import Icon from '@mdi/react';
import axios from 'axios';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router-dom';

import { AsideMenu } from './components/AsideMenu';
import { BreadCrumb } from './components/BreadCrumb';
import { DialogModal } from './components/DialogModal';
import { Footer } from './components/Footer';
import { Header } from './components/Header';
import { Markdown } from './components/Markdown';
import { RouterView } from './components/RouterView';
import { Scroller } from './components/Scroller';
import { Skeleton } from './components/Skeleton';
import { globalConfig } from './config/global';
import { localeConfig } from './config/locale';
import { menuConfig } from './config/menu';
import { routerConfig } from './config/router';
import { useGetSoftwareVersionQuery, useIsGenuineProductLazyQuery } from './graphql';
import { sendUserConfirm } from './helpers/alert/sendUserConfirm';
import { getCtaMessageDataUrl } from './helpers/app/getCtaMessageDataUrl';
import { useCredentialStore } from './stores/credential';
import { useCtaMessageStore } from './stores/ctaMessage';

interface IEntry {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly locales: Record<string, string>;
    readonly onSwitchLocale: (newLocale: string) => void;
}

export const Entry = ({ currentLocale, locales, onSwitchLocale }: IEntry) => {
    const { t } = useTranslation();

    const { data: getSoftwareVersionData, loading: getSoftwareVersionLoading } =
        useGetSoftwareVersionQuery();
    useEffect(() => {
        if (!getSoftwareVersionLoading && getSoftwareVersionData?.getSoftwareVersion) {
            // eslint-disable-next-line no-console
            console.log(
                `%c${getSoftwareVersionData.getSoftwareVersion}`,
                'background: #945cff; color: #ffffff; font-weight: bold; font-size: 12px; padding: 2px 4px;'
            );
        }
    }, [getSoftwareVersionData, getSoftwareVersionLoading]);

    const { pathname } = useLocation();
    const [currentTitle, setCurrentTitle] = useState(globalConfig.name[currentLocale]);
    useEffect(() => {
        for (const key in routerConfig.routes) {
            const { uri } = routerConfig.routes[key];
            if (pathname === uri) {
                setCurrentTitle(routerConfig.routes[key].title[currentLocale]);
                return;
            }
        }
        setCurrentTitle(routerConfig.routes.default.title[currentLocale]);
    }, [pathname, currentLocale]);

    const { clearCredential } = useCredentialStore();
    const handleLogoutSubmit = () => {
        sendUserConfirm(t('Entry.signout.confirm_message'), {
            title: t('Entry.signout.confirm_title'),
            cancelBtnText: t('Entry.signout.cancel_button'),
            confirmBtnText: t('Entry.signout.confirm_button'),
            onConfirmed: clearCredential
        });
    };

    const { status: ctaMessageStatus, markAsRead } = useCtaMessageStore();
    const [ctaMessage, setCtaMessage] = useState<{
        title: string;
        subject: string;
        content: string;
    }>();
    const [ctaModalShouldOpen, setCtaModalShouldOpen] = useState(false);
    const getCtaMessage = useCallback(
        async (resource: string) => {
            const requestFn = async (locale: string) => {
                const url = getCtaMessageDataUrl(resource, locale);
                try {
                    const { data } = await axios.get(url);
                    if (data.title && data.subject && data.content) {
                        return data;
                    }
                } catch {
                    /* empty */
                }
                return null;
            };

            let data = await requestFn(currentLocale);
            if (!data && currentLocale !== localeConfig.fallback) {
                data = await requestFn(localeConfig.fallback);
            }
            if (data) {
                setCtaMessage(data);
            }
        },
        [currentLocale]
    );
    useEffect(() => {
        if (ctaModalShouldOpen) {
            getCtaMessage('license-warning');
        }
    }, [ctaModalShouldOpen, getCtaMessage]);
    const [
        isGenuineProduct,
        {
            data: isGenuineProductData,
            loading: isGenuineProductLoading,
            error: isGenuineProductError
        }
    ] = useIsGenuineProductLazyQuery();
    useEffect(() => {
        if (!ctaMessageStatus.read) {
            isGenuineProduct();
        }
    }, [ctaMessageStatus.read, isGenuineProduct]);
    useEffect(() => {
        if (isGenuineProductData && !isGenuineProductLoading && !isGenuineProductError) {
            setCtaModalShouldOpen(!isGenuineProductData.isGenuineProduct);
        }
    }, [isGenuineProductData, isGenuineProductError, isGenuineProductLoading]);

    return (
        <div className="animate-fade animate-duration-500 animate-delay-300">
            <Header
                title={globalConfig.name[currentLocale]}
                onLogout={handleLogoutSubmit}
                currentLocale={currentLocale}
                onSwitchLocale={onSwitchLocale}
                locales={locales}
            />
            <AsideMenu
                title={globalConfig.name[currentLocale]}
                menu={menuConfig}
                currentLocale={currentLocale}
            />

            <div className="ml-10 flex min-h-screen flex-col space-y-4 p-20 px-4">
                <BreadCrumb
                    pathname={pathname}
                    basename={routerConfig.basename}
                    title={currentTitle}
                />
                <RouterView
                    routerProps={{ currentLocale }}
                    currentLocale={currentLocale}
                    appName={globalConfig.name}
                    routes={routerConfig.routes}
                    suspense={<Skeleton />}
                />
            </div>

            <Footer
                copyright={globalConfig.copyright}
                repository={globalConfig.repository}
                currentLocale={currentLocale}
                text={globalConfig.footer}
                homepage={globalConfig.homepage}
            />
            <Scroller threshold={100} />

            {ctaMessage && (
                <DialogModal
                    enlarge
                    heading={
                        <div className="space-y-4 text-gray-800">
                            <h2 className="flex items-center space-x-2 text-lg font-extrabold">
                                <Icon
                                    className="flex-shrink-0"
                                    path={mdiEmailNewsletter}
                                    size={1}
                                />
                                <span>{ctaMessage.title}</span>
                            </h2>
                            <p className="text-sm">{ctaMessage.subject}</p>
                        </div>
                    }
                    open={ctaModalShouldOpen}
                    onClose={() => {
                        markAsRead();
                        setCtaModalShouldOpen(false);
                    }}
                >
                    <Markdown className="h-[350px] overflow-y-scroll rounded-lg border border-dashed border-gray-300 p-4">
                        {ctaMessage.content}
                    </Markdown>
                </DialogModal>
            )}
        </div>
    );
};
