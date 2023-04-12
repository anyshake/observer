import React, { Component } from "react";
import AppConfig from "../config";

export default class Footer extends Component {
    constructor(props) {
        super(props);
        this.state = {
            copyright: AppConfig.frontend.copyright,
        };
    }

    render() {
        return (
            <footer className="fixed bottom-0 w-full bg-gray-200 text-gray-500">
                <div className="container mx-auto flex flex-wrap flex-col sm:flex-row">
                    <div className="container mx-auto py-2 px-4 flex flex-wrap flex-col sm:flex-row">
                        <span className="text-xs text-center mt-1 ml-4 md:ml-12 lg:ml-16 md:text-left">
                            {this.props.extra ||
                                `Constructing Real-time Seismic Network Ambitiously.`}
                        </span>
                        <span className="text-sm inline-flex sm:ml-auto sm:mt-0 mt-2 justify-center sm:justify-start">
                            {this.state.copyright}
                        </span>
                    </div>
                </div>
            </footer>
        );
    }
}
