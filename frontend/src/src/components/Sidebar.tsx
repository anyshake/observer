import { Component } from "react";
import { Link } from "react-router-dom";
import MENU_CONFIG from "../config/menu";
import MenuIcon from "../assets/icons/maximize-solid.svg";
import getRouterUri from "../helpers/router/getRouterUri";
import { WithTranslation, withTranslation } from "react-i18next";

export interface SidebarState {
    readonly isOpen: boolean;
}

class Sidebar extends Component<WithTranslation, SidebarState> {
    constructor(props: WithTranslation) {
        super(props);
        this.state = {
            isOpen: false,
        };
    }

    render() {
        const { t } = this.props;
        const { isOpen } = this.state;
        const currentUri = getRouterUri();
        const { title, list } = MENU_CONFIG;

        return (
            <aside
                className={`z-20 w-[235px] fixed flex min-h-screen duration-700 bg-gray-800 text-white ${
                    isOpen ? `translate-x-none` : `-translate-x-48`
                }`}
            >
                <div
                    className={`-right-6 ease-in duration-300 border-4 border-white absolute top-2 rounded-full ${
                        isOpen ? `translate-x-0` : `translate-x-24 scale-x-0`
                    }`}
                >
                    <div className="bg-gradient-to-r from-indigo-500 via-purple-500 to-purple-500 pl-16 pr-6 py-2 rounded-full">
                        <div className="duration-100 mr-16 font-bold">
                            {t(title)}
                        </div>
                    </div>
                </div>

                <div
                    className="-right-6 cursor-pointer duration-500 border-4 border-white bg-gray-800 hover:bg-purple-500 absolute top-2 p-3 rounded-full hover:rotate-45"
                    onClick={() =>
                        this.setState({
                            isOpen: !isOpen,
                        })
                    }
                >
                    <img className="w-4 h-4" src={MenuIcon} alt="" />
                </div>

                <div
                    className={`mt-20 flex flex-col space-y-2 w-full h-full ${
                        isOpen || `hidden`
                    }`}
                >
                    {list.map(({ uri, icon, label }, index) => (
                        <Link
                            className={`cursor-pointer w-full bg-gray-800 p-3 pl-8 rounded-full duration-300 flex items-center ${
                                uri === currentUri
                                    ? `font-bold ml-2`
                                    : `hover:font-bold hover:ml-2`
                            }`}
                            to={uri}
                            key={index}
                        >
                            <img src={icon} className="w-4 h-4" alt="" />
                            <span className="ml-4">{t(label)}</span>
                        </Link>
                    ))}
                </div>

                <div
                    className={`mt-20 flex flex-col space-y-2 w-full h-full ${
                        isOpen ? `hidden` : ``
                    }`}
                >
                    {list.map(({ uri, icon }, index) => (
                        <Link
                            key={index}
                            to={uri}
                            className={`cursor-pointer justify-end w-full bg-gray-800 p-4 rounded-full duration-300 flex ${
                                uri === currentUri ? `ml-2` : `hover:ml-2`
                            }`}
                        >
                            <img src={icon} className="w-4 h-4" alt="" />
                        </Link>
                    ))}
                </div>
            </aside>
        );
    }
}

export default withTranslation()(Sidebar);
