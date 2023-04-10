import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";

export default class realtimeWaveform extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "waveform",
        };
    }

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4">
                    <Navbar />
                </div>
            </>
        );
    }
}
