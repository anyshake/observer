import React, { Component } from "react";
import Station from "./components/Station";
import Sidebar from "./components/Sidebar";

export default class App extends Component {
    render() {
        return (
            <div className="bg-gradient-to-br bg-indigo-800 from-indigo-600 via-indigo-800 to-indigo-900">
                <Sidebar />
                <div className="flex flex-col min-h-screen">
                    <div className="container lg:w-3/4 w-8/9 mx-auto mt-[32px]">
                        <div className="w-full px-4 md:px-0 md:mt-8 mb-16 text-gray-800">
                            <Station />
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
