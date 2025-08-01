import {
    mdiAccount,
    mdiAlertDecagram,
    mdiCube,
    mdiHeart,
    mdiHomeEdit,
    mdiPowerPlug,
    mdiTimelineText
} from '@mdi/js';
import Icon from '@mdi/react';
import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Skeleton } from '../../components/Skeleton';
import { IRouterComponent } from '../../config/router';
import { useIsCurrentUserAdminQuery } from '../../graphql';
import { Dangerous } from './dangerous';
import { Logs } from './logs';
import { Metadata } from './metadata';
import { Preference } from './preference';
import { Service } from './service';
import { Station } from './station';
import { Users } from './users';

const Settings = ({ currentLocale }: IRouterComponent) => {
    const { t } = useTranslation();
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
            logs: {
                adminOnly: true,
                icon: mdiTimelineText,
                label: t('views.Settings.Logs.title'),
                element: <Logs />
            },
            dangerous: {
                adminOnly: true,
                icon: mdiAlertDecagram,
                label: t('views.Settings.Dangerous.title'),
                element: <Dangerous />
            }
        };
    }, [t, currentLocale]);
    const [activeTab, setActiveTab] = useState(Object.keys(tabs)[0]);

    const [isAdmin, setIsAdmin] = useState(false);
    const { data: isCurrentUserAdminData, loading: isCurrentUserAdminLoading } =
        useIsCurrentUserAdminQuery();
    useEffect(() => {
        if (isCurrentUserAdminData) {
            setIsAdmin(isCurrentUserAdminData.getCurrentUser.admin);
        }
    }, [isCurrentUserAdminData]);

    return (
        <div className="container mx-auto space-y-6 p-4">
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
                                    className={`tab space-x-1 ${activeTab === tabKey ? 'tab-active text-gray-700' : 'text-gray-500'}`}
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
