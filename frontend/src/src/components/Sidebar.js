import React, { Component } from "react";
import AppConfig from "../config";
import { Link } from "react-router-dom";

export default class Sidebar extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isSidebarOpen: false,
            sidebarVersion: AppConfig.frontend.version,
            sidebarTitle: AppConfig.frontend.title,
            sidebarList: AppConfig.sidebar,
            sidebarMark: this.props.sidebarMark,
        };
    }

    render() {
        return (
            <>
                <div className="fixed w-full z-30 flex bg-white p-2 items-center justify-center h-16 px-10">
                    <div
                        className={`${
                            this.state.isSidebarOpen || `ml-10`
                        } text-gray-800 transform ease-in-out duration-500 flex-none h-full flex items-center justify-center text-lg font-bold`}
                    >
                        {this.state.sidebarTitle}
                    </div>

                    <div className="grow h-full flex items-center justify-center " />
                    <div className="flex-none h-full text-center flex items-center justify-center text-gray-500">
                        <div className="flex space-x-2 items-center lg:px-10">
                            <div className="flex-none flex justify-center">
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 512 512"
                                    className="w-4 h-4"
                                    fill="currentColor"
                                >
                                    <path d="M320 488c0 9.5-5.6 18.1-14.2 21.9s-18.8 2.3-25.8-4.1l-80-72c-5.1-4.6-7.9-11-7.9-17.8s2.9-13.3 7.9-17.8l80-72c7-6.3 17.2-7.9 25.8-4.1s14.2 12.4 14.2 21.9v40h16c35.3 0 64-28.7 64-64V153.3C371.7 141 352 112.8 352 80c0-44.2 35.8-80 80-80s80 35.8 80 80c0 32.8-19.7 61-48 73.3V320c0 70.7-57.3 128-128 128H320v40zM456 80a24 24 0 1 0 -48 0 24 24 0 1 0 48 0zM192 24c0-9.5 5.6-18.1 14.2-21.9s18.8-2.3 25.8 4.1l80 72c5.1 4.6 7.9 11 7.9 17.8s-2.9 13.3-7.9 17.8l-80 72c-7 6.3-17.2 7.9-25.8 4.1s-14.2-12.4-14.2-21.9V128H176c-35.3 0-64 28.7-64 64V358.7c28.3 12.3 48 40.5 48 73.3c0 44.2-35.8 80-80 80s-80-35.8-80-80c0-32.8 19.7-61 48-73.3V192c0-70.7 57.3-128 128-128h16V24zM56 432a24 24 0 1 0 48 0 24 24 0 1 0 -48 0z" />
                                </svg>
                            </div>

                            <div className="md:block text-sm md:text-md">
                                {this.state.sidebarVersion}
                            </div>
                        </div>
                    </div>
                </div>

                <aside
                    className={`${
                        this.state.isSidebarOpen
                            ? `translate-x-none`
                            : `-translate-x-48`
                    } w-60 fixed transition transform ease-in-out duration-1000 z-50 flex h-screen bg-gray-800`}
                >
                    <div
                        className={`${
                            this.state.isSidebarOpen
                                ? `translate-x-0`
                                : `translate-x-24 scale-x-0`
                        } w-full -right-6 transition transform ease-in duration-300 flex items-center justify-between border-4 border-white absolute top-2 rounded-full h-12`}
                    >
                        <div className="flex items-center space-x-3 group bg-gradient-to-r from-indigo-500 via-purple-500 to-purple-500 pl-16 pr-6 py-2 rounded-full text-white">
                            <div className="transform ease-in-out duration-300 mr-16 font-bold">
                                面板菜单
                            </div>
                        </div>
                    </div>

                    <div
                        onClick={() =>
                            this.setState({
                                isSidebarOpen: !this.state.isSidebarOpen,
                            })
                        }
                        className="-right-6 cursor-pointer transition transform ease-in-out duration-500 flex border-4 border-white bg-[#1E293B] hover:bg-purple-500 absolute top-2 p-3 rounded-full text-white hover:rotate-45"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={3}
                            stroke="currentColor"
                            className="w-4 h-4"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"
                            />
                        </svg>
                    </div>

                    <div
                        className={`${
                            this.state.isSidebarOpen ? `flex` : `hidden`
                        } text-white mt-20 flex-col space-y-2 w-full h-[calc(100vh)]`}
                    >
                        {this.state.sidebarList.map((item, index) => (
                            <Link
                                key={index}
                                to={item.link}
                                className={`${
                                    this.state.sidebarMark === item.tag
                                        ? `text-purple-500`
                                        : `hover:text-purple-500`
                                } cursor-pointer hover:ml-4 w-full text-white bg-[#1E293B] p-2 pl-8 rounded-full transform ease-in-out duration-300 flex flex-row items-center space-x-3`}
                            >
                                {item.icon}
                                <div>{item.title}</div>
                            </Link>
                        ))}
                    </div>

                    <div
                        className={`${
                            this.state.isSidebarOpen ? `hidden` : `flex`
                        } mt-20 flex-col space-y-2 w-full h-[calc(100vh)]`}
                    >
                        {this.state.sidebarList.map((item, index) => (
                            <Link
                                key={index}
                                to={item.link}
                                className={`${
                                    this.state.sidebarMark === item.tag
                                        ? `text-purple-500`
                                        : `hover:text-purple-500`
                                } cursor-pointer justify-end pr-5 text-white w-full bg-[#1E293B] p-3 rounded-full transform ease-in-out duration-300 flex`}
                            >
                                {item.icon}
                            </Link>
                        ))}
                    </div>
                </aside>
            </>
        );
    }
}
