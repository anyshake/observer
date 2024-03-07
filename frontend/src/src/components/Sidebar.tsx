import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import MenuIcon from "../assets/icons/maximize-solid.svg";
import { MenuItem } from "../config/menu";

interface SidebarProps {
    readonly title: string;
    readonly links: MenuItem[];
    readonly currentLocale: string;
}

export const Sidebar = (props: SidebarProps) => {
    const { title, links, currentLocale } = props;
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    const location = useLocation();
    const { hash, pathname } = location;

    return (
        <aside
            className={`z-20 w-[235px] fixed flex min-h-screen duration-700 bg-gray-800 text-white ${
                isSidebarOpen ? "translate-x-none" : "-translate-x-48"
            }`}
        >
            <div
                className={`-right-6 ease-in duration-300 border-4 border-white absolute top-2 rounded-full ${
                    isSidebarOpen ? "translate-x-0" : "translate-x-24 scale-x-0"
                }`}
            >
                <div className="bg-gradient-to-r from-indigo-500 via-purple-500 to-purple-500 pl-16 pr-6 py-2 rounded-full">
                    <div className="duration-100 mr-16 font-bold text-center min-w-16">
                        {title}
                    </div>
                </div>
            </div>
            <div
                className="-right-6 cursor-pointer duration-500 border-4 border-white bg-gray-800 hover:bg-purple-500 absolute top-2 p-3 rounded-full hover:rotate-45"
                onClick={() => {
                    setIsSidebarOpen(!isSidebarOpen);
                }}
            >
                <img className="size-4" src={MenuIcon} alt="" />
            </div>
            <div
                className={`mt-20 flex flex-col space-y-2 w-full h-full ${
                    isSidebarOpen ? "" : "hidden"
                }`}
            >
                {links.map(({ url, icon, label }) => (
                    <Link
                        className={`cursor-pointer w-full bg-gray-800 p-3 pl-8 rounded-full duration-300 flex items-center ${
                            url === hash || url === pathname
                                ? "font-bold ml-2"
                                : "hover:font-bold hover:ml-2"
                        }`}
                        to={url}
                        key={url}
                    >
                        <img src={icon} className="size-4" alt="" />
                        <span className="ml-4">{label[currentLocale]}</span>
                    </Link>
                ))}
            </div>
            <div
                className={`mt-20 flex flex-col space-y-2 w-full h-full ${
                    isSidebarOpen ? "hidden" : ""
                }`}
            >
                {links.map(({ url, icon }) => (
                    <Link
                        key={url}
                        to={url}
                        className={`cursor-pointer justify-end w-full bg-gray-800 p-4 rounded-full duration-300 flex ${
                            url === hash || url === pathname
                                ? "ml-2"
                                : "hover:ml-2"
                        }`}
                    >
                        <img src={icon} className="size-4" alt="" />
                    </Link>
                ))}
            </div>
        </aside>
    );
};
