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
                    title: `已解码讯息量`,
                    key: `this.state.response.status.messages`,
                    unit: `条`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 512 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M256 512A256 256 0 1 0 256 0a256 256 0 1 0 0 512zM369 209L241 337c-9.4 9.4-24.6 9.4-33.9 0l-64-64c-9.4-9.4-9.4-24.6 0-33.9s24.6-9.4 33.9 0l47 47L335 175c9.4-9.4 24.6-9.4 33.9 0s9.4 24.6 0 33.9z" />
                        </svg>
                    ),
                },
                {
                    title: `帧错误讯息量`,
                    key: `this.state.response.status.errors`,
                    unit: `条`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 512 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M256 0c53 0 96 43 96 96v3.6c0 15.7-12.7 28.4-28.4 28.4H188.4c-15.7 0-28.4-12.7-28.4-28.4V96c0-53 43-96 96-96zM41.4 105.4c12.5-12.5 32.8-12.5 45.3 0l64 64c.7 .7 1.3 1.4 1.9 2.1c14.2-7.3 30.4-11.4 47.5-11.4H312c17.1 0 33.2 4.1 47.5 11.4c.6-.7 1.2-1.4 1.9-2.1l64-64c12.5-12.5 32.8-12.5 45.3 0s12.5 32.8 0 45.3l-64 64c-.7 .7-1.4 1.3-2.1 1.9c6.2 12 10.1 25.3 11.1 39.5H480c17.7 0 32 14.3 32 32s-14.3 32-32 32H416c0 24.6-5.5 47.8-15.4 68.6c2.2 1.3 4.2 2.9 6 4.8l64 64c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0l-63.1-63.1c-24.5 21.8-55.8 36.2-90.3 39.6V240c0-8.8-7.2-16-16-16s-16 7.2-16 16V479.2c-34.5-3.4-65.8-17.8-90.3-39.6L86.6 502.6c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3l64-64c1.9-1.9 3.9-3.4 6-4.8C101.5 367.8 96 344.6 96 320H32c-17.7 0-32-14.3-32-32s14.3-32 32-32H96.3c1.1-14.1 5-27.5 11.1-39.5c-.7-.6-1.4-1.2-2.1-1.9l-64-64c-12.5-12.5-12.5-32.8 0-45.3z" />
                        </svg>
                    ),
                },
                {
                    title: `已推送讯息量`,
                    key: `this.state.response.status.pushed`,
                    unit: `条`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 512 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M498.1 5.6c10.1 7 15.4 19.1 13.5 31.2l-64 416c-1.5 9.7-7.4 18.2-16 23s-18.9 5.4-28 1.6L284 427.7l-68.5 74.1c-8.9 9.7-22.9 12.9-35.2 8.1S160 493.2 160 480V396.4c0-4 1.5-7.8 4.2-10.7L331.8 202.8c5.8-6.3 5.6-16-.4-22s-15.7-6.4-22-.7L106 360.8 17.7 316.6C7.1 311.3 .3 300.7 0 288.9s5.9-22.8 16.1-28.7l448-256c10.7-6.1 23.9-5.5 34 1.4z" />
                        </svg>
                    ),
                },
                {
                    title: `推送失败讯息量`,
                    key: `this.state.response.status.failures`,
                    unit: `条`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 512 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M256 512A256 256 0 1 0 256 0a256 256 0 1 0 0 512zm0-384c13.3 0 24 10.7 24 24V264c0 13.3-10.7 24-24 24s-24-10.7-24-24V152c0-13.3 10.7-24 24-24zM224 352a32 32 0 1 1 64 0 32 32 0 1 1 -64 0z" />
                        </svg>
                    ),
                },
                {
                    title: `等待推送讯息量`,
                    key: `this.state.response.status.queued`,
                    unit: `条`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 384 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M32 0C14.3 0 0 14.3 0 32S14.3 64 32 64V75c0 42.4 16.9 83.1 46.9 113.1L146.7 256 78.9 323.9C48.9 353.9 32 394.6 32 437v11c-17.7 0-32 14.3-32 32s14.3 32 32 32H64 320h32c17.7 0 32-14.3 32-32s-14.3-32-32-32V437c0-42.4-16.9-83.1-46.9-113.1L237.3 256l67.9-67.9c30-30 46.9-70.7 46.9-113.1V64c17.7 0 32-14.3 32-32s-14.3-32-32-32H320 64 32zM96 75V64H288V75c0 25.5-10.1 49.9-28.1 67.9L192 210.7l-67.9-67.9C106.1 124.9 96 100.4 96 75z" />
                        </svg>
                    ),
                },
                {
                    title: `系统时间偏移量`,
                    key: `this.state.response.status.offset`,
                    unit: `秒`,
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 512 512"
                            fill="currentColor"
                            className="w-8 h-8"
                        >
                            <path d="M256 0a256 256 0 1 1 0 512A256 256 0 1 1 256 0zM232 120V256c0 8 4 15.5 10.7 20l96 64c11 7.4 25.9 4.4 33.3-6.7s4.4-25.9-6.7-33.3L280 243.2V120c0-13.3-10.7-24-24-24s-24 10.7-24 24z" />
                        </svg>
                    ),
                },
            ],
            response: {
                uuid: ``,
                station: ``,
                uptime: 0,
                location: {
                    latitude: -1,
                    longitude: -1,
                    altitude: -1,
                },
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
                promise={this.fetchData}
                render={({ isPolling }) => (
                    <>
                        <Sidebar sidebarMark={this.state.sidebarMark} />
                        <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4 ">
                            <div
                                className={
                                    this.state.response.station.length > 0
                                        ? `p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-cyan-500 to-yellow-500`
                                        : `p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-blue-500 to-orange-500`
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
                                            {`测站位置：经纬 ${this.state.response.location.latitude},
                                            ${this.state.response.location.longitude} / 海拔
                                            ${this.state.response.location.altitude}m`}
                                            <br />
                                            {`测站标识符：${this.state.response.uuid}`}
                                            <br />
                                            {`服务器已上线 ${this.state.response.uptime} 秒`}
                                        </h3>
                                    ) : (
                                        <h3 className="pl-3 text-md font-medium">
                                            {`测站名称：正在获取`}
                                            <br />
                                            {`测站位置：正在获取`}
                                            <br />
                                            {`测站标识符：正在获取`}
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
