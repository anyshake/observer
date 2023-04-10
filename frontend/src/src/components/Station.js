import React, { Component } from "react";
import getApiUrl from "../helpers/utilities/getApiUrl";
import AppConfig from "../config";
import createRequest from "../helpers/requests/createRequest";
import ReactPolling from "react-polling";
import { timerAlert } from "../helpers/alerts/sweetAlert";

export default class Station extends Component {
    constructor(props) {
        super(props);
        this.state = {
            cards: [
                {
                    name: `解码讯息量`,
                    color: `bg-green-600`,
                    key: `this.state.response.status.messages`,
                    unit: `条`,
                },
                {
                    name: `错误讯息量`,
                    color: `bg-red-600`,
                    key: `this.state.response.status.errors`,
                    unit: `条`,
                },
                {
                    name: `已推送讯息`,
                    color: `bg-indigo-600`,
                    key: `this.state.response.status.pushed`,
                    unit: `条`,
                },
                {
                    name: `推送失败讯息`,
                    color: `bg-pink-600`,
                    key: `this.state.response.status.failures`,
                    unit: `条`,
                },
                {
                    name: `待推送讯息`,
                    color: `bg-yellow-600`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `错误讯息量`,
                    color: `bg-red-600`,
                    key: `this.state.response.status.errors`,
                    unit: `条`,
                },
                {
                    name: `已推送讯息`,
                    color: `bg-indigo-600`,
                    key: `this.state.response.status.pushed`,
                    unit: `条`,
                },
                {
                    name: `推送失败讯息`,
                    color: `bg-pink-600`,
                    key: `this.state.response.status.failures`,
                    unit: `条`,
                },
                {
                    name: `待推送讯息`,
                    color: `bg-yellow-600`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `待推送讯息`,
                    color: `bg-yellow-600`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `待推送讯息`,
                    color: `bg-yellow-600`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    name: `时间偏移量`,
                    color: `bg-blue-600`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
            ],
            response: {
                uptime: 0,
                memory: {
                    total: 0,
                    free: 0,
                    used: 0,
                    percent: 0,
                },
                disk: {
                    total: 0,
                    free: 0,
                    used: 0,
                    percent: 0,
                },
                status: {
                    messages: 0,
                    pushed: 0,
                    errors: 0,
                    failures: 0,
                    queued: 0,
                    offset: 0,
                },
                os: {
                    os: ``,
                    arch: ``,
                    distro: ``,
                    hostname: ``,
                },
                cpu: {
                    model: ``,
                    percent: 0,
                },
            },
        };
    }

    fetchData = (url) => {
        return createRequest({
            url: url,
            method: AppConfig.backend.api.station.method,
        });
    };

    render() {
        return (
            <ReactPolling
                interval={1000}
                url={getApiUrl({
                    tls: AppConfig.backend.tls,
                    host: AppConfig.backend.host,
                    port: AppConfig.backend.port,
                    version: AppConfig.backend.version,
                    api: AppConfig.backend.api.station.uri,
                    type: AppConfig.backend.api.station.type,
                })}
                onSuccess={(res) => {
                    const {
                        data: { data },
                    } = res;
                    this.setState({ response: data });
                    return res;
                }}
                onFailure={() =>
                    timerAlert({
                        title: "连接失败",
                        html: "请检查网络连接，页面即将刷新",
                        loading: false,
                        timer: 3000,
                        callback: () => {
                            window.location.reload();
                        },
                    })
                }
                retryCount={3}
                promise={this.fetchData}
                render={({ isPolling }) => {
                    if (!isPolling) {
                        timerAlert({
                            title: "正在连接",
                            html: "正在连接服务器",
                            loading: true,
                            timer: 1000,
                        });
                    }

                    return (
                        <div className="flex flex-wrap">
                            {this.state.cards.map((item, index) => (
                                <div
                                    className="w-full md:w-1/2 xl:w-1/3 p-3"
                                    key={index}
                                >
                                    <div className="bg-gray-200 border rounded-lg shadow p-2">
                                        <div className="flex flex-row items-center">
                                            <div className="flex-shrink pr-4">
                                                <div
                                                    className={`rounded-sm p-3 ${item.color}`}
                                                />
                                            </div>
                                            <div className="flex-1 text-right md:text-center">
                                                <h5 className="font-bold uppercase text-gray-500">
                                                    {item.name}
                                                </h5>
                                                <h3 className="font-bold text-3xl">
                                                    {
                                                        // eslint-disable-next-line no-eval
                                                        `${eval(item.key)} ${
                                                            item.unit
                                                        }`
                                                    }
                                                </h3>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    );
                }}
            />
        );
    }
}
