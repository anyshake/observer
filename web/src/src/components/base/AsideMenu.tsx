import { mdiArrowExpandAll } from '@mdi/js';
import Icon from '@mdi/react';
import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

import { localeConfig } from '../../config/locale';
import { IMenuItem } from '../../config/menu';

interface IAsideMenu {
    readonly title: string;
    readonly menu: IMenuItem[];
    readonly currentLocale: keyof typeof localeConfig.resources;
}

export const AsideMenu = ({ title, menu, currentLocale }: IAsideMenu) => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);
    const { hash, pathname } = useLocation();

    return (
        <aside
            className={`fixed z-20 flex min-h-screen w-[235px] bg-gray-800 text-white transition-all duration-700 ${
                isSidebarOpen ? 'left-0' : '-left-48'
            }`}
        >
            <div
                className={`absolute top-2 -right-6 rounded-full border-4 border-white duration-300 ease-in ${
                    isSidebarOpen ? 'translate-x-0' : 'translate-x-24 scale-x-0'
                }`}
            >
                <div className="rounded-full bg-gradient-to-tl from-purple-500 to-indigo-500 py-3 pr-14">
                    <div className="text-md ml-6 text-center font-bold duration-100 select-none">
                        {title}
                    </div>
                </div>
            </div>
            <div
                className={`absolute top-2.5 -right-6 cursor-pointer rounded-full border-4 border-white p-3 duration-500 hover:rotate-45 hover:bg-purple-500 ${isSidebarOpen ? 'rotate-45 bg-purple-500' : 'bg-gray-800'}`}
                onClick={() => {
                    setIsSidebarOpen(!isSidebarOpen);
                }}
            >
                <Icon className="flex-shrink-0" path={mdiArrowExpandAll} size={0.8} />
            </div>
            <div
                className={`mt-20 flex h-full w-full flex-col space-y-2 ${
                    isSidebarOpen ? '' : 'hidden'
                }`}
            >
                {menu.map(({ url, icon, label }) => (
                    <Link
                        className={`flex w-full cursor-pointer items-center rounded-full bg-gray-800 p-3 pl-6 duration-300 ${
                            url === hash || url === pathname
                                ? 'ml-2 font-bold'
                                : 'hover:ml-2 hover:font-bold'
                        }`}
                        to={url}
                        key={url}
                    >
                        <Icon className="flex-shrink-0" path={icon} size={0.8} />
                        <span className="ml-4">{label[currentLocale]}</span>
                    </Link>
                ))}
            </div>
            <div
                className={`mt-20 flex h-full w-full flex-col space-y-2 ${
                    isSidebarOpen ? 'hidden' : ''
                }`}
            >
                {menu.map(({ url, icon }) => (
                    <Link
                        key={url}
                        to={url}
                        className={`flex w-full cursor-pointer justify-end rounded-full bg-gray-800 p-3 duration-300 ${
                            url === hash || url === pathname ? 'ml-2' : 'hover:ml-2'
                        }`}
                    >
                        <Icon className="flex-shrink-0" path={icon} size={0.8} />
                    </Link>
                ))}
            </div>
        </aside>
    );
};
