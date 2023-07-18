import { Component } from "react";
import { Link } from "react-router-dom";
import HomeIcon from "../assets/icons/house-solid.svg";
import ArrowIcon from "../assets/icons/angle-right-solid.svg";
import getRouterUri from "../helpers/getRouterUri";
import getRouterTitle from "../helpers/getRouterTitle";

export default class Navbar extends Component<{}> {
    render() {
        const uri = getRouterUri();
        const title = getRouterTitle();

        return (
            <nav className="px-5 py-3 rounded-lg bg-gray-100">
                <ol className="text-sm font-medium text-gray-700 flex space-x-2">
                    <li className="cursor-pointer hover:text-gray-900">
                        <Link className="flex" to={"/"}>
                            <img
                                className="self-center w-4 h-4"
                                src={HomeIcon}
                                alt=""
                            />
                            <span className="ml-2">主页</span>
                        </Link>
                    </li>

                    {uri !== "/" && (
                        <li className="flex">
                            <img
                                className="self-center w-4 h-4"
                                src={ArrowIcon}
                                alt=""
                            />
                            <Link
                                className="ml-2 cursor-pointer hover:text-gray-900"
                                to={uri}
                            >
                                {title}
                            </Link>
                        </li>
                    )}
                </ol>
            </nav>
        );
    }
}
