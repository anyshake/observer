import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Footer from "../../components/Footer";
import Scroller from "../../components/Scroller";

export default class historyWaveform extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "history",
        };
    }

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4 ">
                    <Navbar navigation={"历史数据"} />
                </div>

                <Scroller />
                <Footer />
            </>
        );
    }
}
