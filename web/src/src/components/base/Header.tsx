import { mdiEarth, mdiExitToApp } from '@mdi/js';
import Icon from '@mdi/react';

import { localeConfig } from '../../config/locale';

interface IHeader {
    readonly title: string;
    readonly onLogout: () => void;
    readonly locales: Record<string, string>;
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly onSwitchLocale: (newLocale: string) => void;
}

export const Header = ({ title, onLogout, locales, currentLocale, onSwitchLocale }: IHeader) => {
    return (
        <div className="navbar bg-base-200 fixed z-10 px-4 shadow-md">
            <div className="flex-1 pl-16">
                <span className="text-xl font-bold text-gray-800 select-none">{title}</span>
            </div>
            <div className="flex space-x-1">
                <div className="dropdown dropdown-end">
                    <button tabIndex={0} className="btn btn-sm btn-ghost text-gray-500">
                        <Icon className="flex-shrink-0" path={mdiEarth} size={0.8} />
                    </button>
                    <ul
                        tabIndex={0}
                        className="dropdown-content menu bg-base-100 rounded-box w-36 p-2 shadow"
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
                <button className="btn btn-sm btn-ghost text-gray-500" onClick={onLogout}>
                    <Icon className="flex-shrink-0" path={mdiExitToApp} size={0.8} />
                </button>
            </div>
        </div>
    );
};
