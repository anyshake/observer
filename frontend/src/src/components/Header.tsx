import { Component } from "react";
// import BranchIcon from "../assets/icons/code-branch-solid.svg";
import GLOBAL_CONFIG from "../config/global";

export default class Header extends Component {
    render() {
        // const { name, version } = GLOBAL_CONFIG.app_settings;
        const { name } = GLOBAL_CONFIG.app_settings;
        return (
            <header className="fixed w-full z-10 flex justify-between bg-gray-200 items-center h-16 px-5">
                <h1 className="ml-14 text-gray-800 text-xl font-bold">
                    {name}
                </h1>

                <div className="flex text-gray-500">
                    {/* <img className="w-4 h-4" src={BranchIcon} alt="" />
                    <span className="ml-1 text-sm">{version}</span> */}
                </div>
            </header>
        );
    }
}
