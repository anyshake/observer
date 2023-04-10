import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Card from "../../components/Card";

export default class stationInfo extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "index",
            cardList: [{}, {}, {}, {}, {}, {}],
        };
    }

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4 ">
                    <div
                        class="p-4 mb-4 text-sm text-blue-700 bg-blue-100 rounded-lg dark:bg-blue-200 dark:text-blue-800"
                        role="alert"
                    >
                        <span class="font-medium">Info alert!</span> Change a
                        few things up and try submitting again.
                    </div>
                    <div
                        class="p-4 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-red-200 dark:text-red-800"
                        role="alert"
                    >
                        <span class="font-medium">Danger alert!</span> Change a
                        few things up and try submitting again.
                    </div>
                    <div
                        class="p-4 mb-4 text-sm text-green-700 bg-green-100 rounded-lg dark:bg-green-200 dark:text-green-800"
                        role="alert"
                    >
                        <span class="font-medium">Success alert!</span> Change a
                        few things up and try submitting again.
                    </div>
                    <Card cardList={this.state.cardList} />
                </div>
            </>
        );
    }
}
