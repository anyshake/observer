import React, { Component } from "react";
import redirectRouter from "../helpers/router/redirectRouter";

export default class Navbar extends Component {
    constructor(props) {
        super(props);
        this.state = {
            navigation: this.props.navigation || "",
        };
    }

    render() {
        return (
            <nav className="flex px-5 py-3 text-gray-700  rounded-lg bg-gray-50">
                <ol className="inline-flex items-center space-x-1 md:space-x-3">
                    <li className="inline-flex items-center">
                        <div
                            onClick={() =>
                                redirectRouter({
                                    dest: "/",
                                    replace: false,
                                })
                            }
                            className="cursor-pointer inline-flex items-center text-sm font-medium text-gray-700 hover:text-gray-900"
                        >
                            <svg
                                className="w-4 h-4 mr-2"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"></path>
                            </svg>
                            主页
                        </div>
                    </li>
                    {this.state.navigation.length > 0 && (
                        <li>
                            <div className="flex items-center">
                                <svg
                                    className="w-6 h-6 text-gray-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        fillRule="evenodd"
                                        d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                                        clipRule="evenodd"
                                    />
                                </svg>
                                <div className="ml-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ml-2">
                                    {this.state.navigation}
                                </div>
                            </div>
                        </li>
                    )}
                </ol>
            </nav>
        );
    }
}
