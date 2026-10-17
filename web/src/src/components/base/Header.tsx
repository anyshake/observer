import { mdiEarth, mdiExitToApp } from '@mdi/js';
import Icon from '@mdi/react';

import { localeConfig } from '../../config/locale';
import { ThemeToggle } from '../ui/ThemeToggle';

interface IHeader {
    readonly title: string;
    readonly onLogout: () => void;
    readonly locales: Record<string, string>;
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly onSwitchLocale: (newLocale: string) => void;
}

export const Header = ({ title, onLogout, locales, currentLocale, onSwitchLocale }: IHeader) => {
    return (
        <div className="navbar app-header">
            <div className="app-header-brand">
                <span className="app-header-title">{title}</span>
            </div>
            <div className="app-header-actions">
                <ThemeToggle />
                <div className="dropdown dropdown-end">
                    <button tabIndex={0} className="btn btn-sm btn-ghost app-header-button">
                        <Icon className="app-icon" path={mdiEarth} size={0.8} />
                    </button>
                    <ul
                        tabIndex={0}
                        className="dropdown-content menu rounded-box app-header-locales"
                    >
                        {Object.entries(locales).map(([key, value]) => (
                            <li
                                className={`app-text-body ${key === currentLocale ? 'app-text-bold' : ''}`}
                                key={key}
                            >
                                <a onClick={() => onSwitchLocale(key)}>{value}</a>
                            </li>
                        ))}
                    </ul>
                </div>
                <button className="btn btn-sm btn-ghost app-header-button" onClick={onLogout}>
                    <Icon className="app-icon" path={mdiExitToApp} size={0.8} />
                </button>
            </div>
        </div>
    );
};
