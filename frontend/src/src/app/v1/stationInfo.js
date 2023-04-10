import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Card from "../../components/Card";
import createRequest from "../../helpers/requests/createRequest";
import AppConfig from "../../config";
import { timerAlert } from "../../helpers/alerts/sweetAlert";
import getApiUrl from "../../helpers/utilities/getApiUrl";
import ReactPolling from "react-polling";

export default class stationInfo extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "index",
            cards: [
                {
                    title: `解码讯息量`,
                    key: `this.state.response.status.messages`,
                    unit: `条`,
                },
                {
                    title: `错误讯息量`,
                    key: `this.state.response.status.errors`,
                    unit: `条`,
                },
                {
                    title: `已推送讯息`,
                    key: `this.state.response.status.pushed`,
                    unit: `条`,
                },
                {
                    title: `推送失败讯息`,
                    key: `this.state.response.status.failures`,
                    unit: `条`,
                },
                {
                    title: `待推送讯息`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                },
                {
                    title: `时间偏移量`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                },
                {
                    title: `错误讯息量`,
                    key: `this.state.response.status.errors`,
                    unit: `条`,
                },
                {
                    title: `已推送讯息`,
                    key: `this.state.response.status.pushed`,
                    unit: `条`,
                },
                {
                    title: `推送失败讯息`,
                    key: `this.state.response.status.failures`,
                    unit: `条`,
                },
            ],
            response: {
                uuid: ``,
                station: ``,
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
                render={({ isPolling }) => (
                    <>
                        <Sidebar sidebarMark={this.state.sidebarMark} />
                        <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4 ">
                            <div
                                className={
                                    this.state.response.station.length > 0
                                        ? `p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-cyan-500 to-yellow-500`
                                        : `p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-blue-500 to-pink-500`
                                }
                            >
                                <div className="flex flex-col gap-y-2">
                                    <div className="flex flex-row space-y-1 gap-2 font-bold">
                                        {this.state.response.uptime ? (
                                            <>
                                                <svg
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    viewBox="0 0 448 512"
                                                    className="w-6 h-6 ml-3"
                                                    fill="currentColor"
                                                >
                                                    <path d="M0 64C0 46.3 14.3 32 32 32c229.8 0 416 186.2 416 416c0 17.7-14.3 32-32 32s-32-14.3-32-32C384 253.6 226.4 96 32 96C14.3 96 0 81.7 0 64zM0 416a64 64 0 1 1 128 0A64 64 0 1 1 0 416zM32 160c159.1 0 288 128.9 288 288c0 17.7-14.3 32-32 32s-32-14.3-32-32c0-123.7-100.3-224-224-224c-17.7 0-32-14.3-32-32s14.3-32 32-32z" />
                                                </svg>
                                                <span>服务器连接已建立</span>
                                            </>
                                        ) : (
                                            <>
                                                <svg
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    viewBox="0 0 640 512"
                                                    className="w-6 h-6 ml-3"
                                                    fill="currentColor"
                                                >
                                                    <path d="M579.8 267.7c56.5-56.5 56.5-148 0-204.5c-50-50-128.8-56.5-186.3-15.4l-1.6 1.1c-14.4 10.3-17.7 30.3-7.4 44.6s30.3 17.7 44.6 7.4l1.6-1.1c32.1-22.9 76-19.3 103.8 8.6c31.5 31.5 31.5 82.5 0 114L422.3 334.8c-31.5 31.5-82.5 31.5-114 0c-27.9-27.9-31.5-71.8-8.6-103.8l1.1-1.6c10.3-14.4 6.9-34.4-7.4-44.6s-34.4-6.9-44.6 7.4l-1.1 1.6C206.5 251.2 213 330 263 380c56.5 56.5 148 56.5 204.5 0L579.8 267.7zM60.2 244.3c-56.5 56.5-56.5 148 0 204.5c50 50 128.8 56.5 186.3 15.4l1.6-1.1c14.4-10.3 17.7-30.3 7.4-44.6s-30.3-17.7-44.6-7.4l-1.6 1.1c-32.1 22.9-76 19.3-103.8-8.6C74 372 74 321 105.5 289.5L217.7 177.2c31.5-31.5 82.5-31.5 114 0c27.9 27.9 31.5 71.8 8.6 103.9l-1.1 1.6c-10.3 14.4-6.9 34.4 7.4 44.6s34.4 6.9 44.6-7.4l1.1-1.6C433.5 260.8 427 182 377 132c-56.5-56.5-148-56.5-204.5 0L60.2 244.3z" />
                                                </svg>
                                                <span>正在建立服务器连接</span>
                                            </>
                                        )}
                                    </div>
                                    {this.state.response.uptime ? (
                                        <h3 className="pl-3 text-md font-medium">
                                            {`测站名称：${this.state.response.station}`}
                                            <br />
                                            {`测站标识符：${this.state.response.uuid}`}
                                            <br />
                                            {`服务器已上线 ${this.state.response.uptime} 秒`}
                                        </h3>
                                    ) : (
                                        <h3 className="pl-3 text-md font-medium">
                                            {`测站名称：未知`}
                                            <br />
                                            {`测站标识符：未知`}
                                            <br />
                                            {`服务器已上线 0 秒`}
                                        </h3>
                                    )}
                                </div>
                            </div>

                            <div className="flex flex-wrap my-2 -mx-2">
                                {this.state.cards.map((item, index) => (
                                    <Card
                                        key={index}
                                        icon={item.icon}
                                        title={item.title}
                                        unit={item.unit}
                                        value={
                                            // eslint-disable-next-line no-eval
                                            eval(item.key)
                                        }
                                    />
                                ))}
                            </div>
                        </div>
                    </>
                )}
            />
        );
    }
}
