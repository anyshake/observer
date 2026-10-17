import { mdiArrowExpandAll } from '@mdi/js';
import Icon from '@mdi/react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation } from 'react-router-dom';

import { IMenuItem } from '../../config/menu';

interface IAsideMenu {
    readonly title: string;
    readonly menu: IMenuItem[];
}

export const AsideMenu = ({ title, menu }: IAsideMenu) => {
    const { pathname } = useLocation();
    const [t] = useTranslation();
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    return (
        <aside
            className={`app-sidebar ${
                isSidebarOpen ? 'app-sidebar-open' : 'app-sidebar-collapsed'
            }`}
        >
            <div
                className={`app-sidebar-brand ${
                    isSidebarOpen ? 'app-sidebar-brand-visible' : 'app-sidebar-brand-hidden'
                }`}
            >
                <div className="app-sidebar-brand-surface">
                    <div className="app-sidebar-title">{title}</div>
                </div>
            </div>
            <div
                className={`app-sidebar-toggle ${isSidebarOpen ? 'app-sidebar-toggle-open' : 'app-sidebar-toggle-collapsed'}`}
                onClick={() => {
                    setIsSidebarOpen(!isSidebarOpen);
                }}
            >
                <Icon className="app-icon" path={mdiArrowExpandAll} size={0.8} />
            </div>
            <div className={`app-sidebar-menu ${isSidebarOpen ? '' : 'app-hidden'}`}>
                {menu.map(({ url, icon, label }) => (
                    <Link
                        className={`app-sidebar-link ${
                            url === pathname ? 'app-sidebar-link-active' : 'app-sidebar-link-idle'
                        }`}
                        to={url}
                        key={url}
                    >
                        <Icon className="app-icon" path={icon} size={0.8} />
                        <span className="app-sidebar-label">{t(label)}</span>
                    </Link>
                ))}
            </div>
            <div className={`app-sidebar-menu ${isSidebarOpen ? 'app-hidden' : ''}`}>
                {menu.map(({ url, icon }) => (
                    <Link
                        key={url}
                        to={url}
                        className={`app-sidebar-icon-link ${
                            url === pathname ? 'app-sidebar-icon-active' : 'app-sidebar-icon-idle'
                        }`}
                    >
                        <Icon className="app-icon" path={icon} size={0.8} />
                    </Link>
                ))}
            </div>
        </aside>
    );
};
