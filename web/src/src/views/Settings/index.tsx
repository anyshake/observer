import {
    mdiAccount,
    mdiAlertDecagram,
    mdiBackupRestore,
    mdiCube,
    mdiHeart,
    mdiHomeEdit,
    mdiKey,
    mdiPowerPlug,
    mdiTimelineText
} from '@mdi/js';
import Icon from '@mdi/react';
import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';

import { Skeleton } from '../../components/ui/Skeleton';
import { IRouterComponent } from '../../config/router';
import { useIsCurrentUserAdminQuery } from '../../graphql';
import { AppLogs } from './applogs';
import { Backup } from './backup';
import { Dangerous } from './dangerous';
import { FDSN } from './fdsnws';
import { Metadata } from './metadata';
import { Preference } from './preference';
import { Service } from './service';
import { Station } from './station';
import { Users } from './users';

const Settings = ({ currentLocale }: IRouterComponent) => {
    const { t } = useTranslation();
    const [searchParams, setSearchParams] = useSearchParams();

    const tabs = useMemo(() => {
        return {
            inventory: {
                adminOnly: false,
                icon: mdiCube,
                label: t('views.Settings.Metadata.title'),
                element: <Metadata />
            },
            preference: {
                adminOnly: false,
                icon: mdiHeart,
                label: t('views.Settings.Preference.title'),
                element: <Preference />
            },
            station: {
                adminOnly: true,
                icon: mdiHomeEdit,
                label: t('views.Settings.Station.title'),
                element: <Station />
            },
            service: {
                adminOnly: true,
                icon: mdiPowerPlug,
                label: t('views.Settings.Service.title'),
                element: <Service />
            },
            users: {
                adminOnly: true,
                icon: mdiAccount,
                label: t('views.Settings.Users.title'),
                element: <Users currentLocale={currentLocale} />
            },
            fdsn: {
                adminOnly: true,
                icon: mdiKey,
                label: t('views.Settings.FDSN.title'),
                element: <FDSN currentLocale={currentLocale} />
            },
            backup: {
                adminOnly: true,
                icon: mdiBackupRestore,
                label: t('views.Settings.Backup.title'),
                element: <Backup />
            },
            logs: {
                adminOnly: true,
                icon: mdiTimelineText,
                label: t('views.Settings.Logs.title'),
                element: <AppLogs />
            },
            dangerous: {
                adminOnly: true,
                icon: mdiAlertDecagram,
                label: t('views.Settings.Dangerous.title'),
                element: <Dangerous currentLocale={currentLocale} />
            }
        };
    }, [t, currentLocale]);
    const [activeTab, setActiveTab] = useState(
        (() => {
            const tabFromUrl = searchParams.get('tab');
            const availableTabs = Object.keys(tabs);
            if (tabFromUrl && availableTabs.includes(tabFromUrl)) {
                return tabFromUrl;
            }
            return availableTabs[0];
        })()
    );
    useEffect(() => {
        if (activeTab) {
            setSearchParams({ tab: activeTab });
        }
    }, [activeTab, setSearchParams]);

    const [isAdmin, setIsAdmin] = useState(false);
    const { data: isCurrentUserAdminData, loading: isCurrentUserAdminLoading } =
        useIsCurrentUserAdminQuery();
    useEffect(() => {
        if (isCurrentUserAdminData) {
            setIsAdmin(isCurrentUserAdminData.getCurrentUser.admin);
        }
    }, [isCurrentUserAdminData]);

    return (
        <div className="app-page">
            {isCurrentUserAdminLoading ? (
                <Skeleton />
            ) : (
                <>
                    <div role="tablist" className="tabs tabs-lift">
                        {Object.keys(tabs)
                            .filter(
                                (tabKey) => !tabs[tabKey as keyof typeof tabs].adminOnly || isAdmin
                            )
                            .map((tabKey, index) => (
                                <div
                                    key={`${index}_${tabKey}`}
                                    role="tab"
                                    className={`tab app-settings-tab ${activeTab === tabKey ? 'tab-active app-text-body' : 'app-text-muted'}`}
                                    onClick={() => setActiveTab(tabKey)}
                                >
                                    <Icon
                                        path={tabs[tabKey as keyof typeof tabs].icon}
                                        size={0.8}
                                    />
                                    <span>{tabs[tabKey as keyof typeof tabs].label}</span>
                                </div>
                            ))}
                    </div>

                    {tabs[activeTab as keyof typeof tabs].element}
                </>
            )}
        </div>
    );
};

export default Settings;
