import React, { Component } from "react";
import createSocket from "./helpers/requests/socket";
import AppConfig from "./config";

export default class App extends Component {
    componentDidMount() {
        const wsUrl = `${AppConfig.backend.tls ? "wss" : "ws"}://${
            AppConfig.backend.host
        }:${AppConfig.backend.port}/api/v1/socket`;

        createSocket({
            url: wsUrl,
            type: "arraybuffer",
            onMessageCallback: ({ data }) => {
                console.log(JSON.parse(data));
            },
            onCloseCallback: (e) => {
                console.log(e);
            },
        });
    }

    render() {
        return (
            <div className="min-h-screen bg-cyan-600 text-amber-500 flex flex-col justify-center">
                <h1 className="text-center text-5xl font-bold tracking-wide">
                    地震项目测试前端
                </h1>
            </div>
        );
    }
}
