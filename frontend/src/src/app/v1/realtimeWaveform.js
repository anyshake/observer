import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Footer from "../../components/Footer";
import Scroller from "../../components/Scroller";
import createSocket from "../../helpers/requests/createSocket";
import getApiUrl from "../../helpers/utilities/getApiUrl";
import AppConfig from "../../config";
import ReactApexChart from "react-apexcharts";
import getTime from "../../helpers/utilities/getTime";
import Notification from "../../components/Notification";
import Navbar from "../../components/Navbar";

export default class realtimeWaveform extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "waveform",
            webSocket: null,
            lastDataArr: {
                timestamp: -1,
                vertical: [],
                east_west: [],
                north_south: [],
            },
            waveform: {
                factors: [
                    {
                        name: "垂直分量",
                        color: "#d97706",
                        data: [],
                    },
                    {
                        name: "水平 EW",
                        color: "#0d9488",
                        data: [],
                    },
                    {
                        name: "水平 NS",
                        color: "#4f46e5",
                        data: [],
                    },
                ],
                synthesis: [
                    {
                        name: "合成分量",
                        color: "#be185d",
                        data: [],
                    },
                ],
                options: {
                    stroke: {
                        width: 2,
                        curve: "smooth",
                    },
                    hollow: {
                        margin: 15,
                        size: "40%",
                    },
                    chart: {
                        height: 350,
                        toolbar: {
                            show: false,
                        },
                        zoom: {
                            enabled: false,
                        },
                        animations: {
                            enabled: false,
                            // easing: "linear",
                            // dynamicAnimation: {
                            //     speed: 1000 * 1.1,
                            // },
                        },
                    },
                    dataLabels: {
                        enabled: false,
                    },
                    legend: {
                        show: false,
                        labels: {
                            useSeriesColors: true,
                        },
                    },
                    tooltip: {
                        enabled: true,
                        theme: "dark",
                        fillSeriesColor: false,
                        x: {
                            format: "20yy/MM/dd HH:mm:ss",
                        },
                    },
                    xaxis: {
                        type: "datetime",
                        labels: {
                            datetimeUTC: false,
                            datetimeFormatter: {
                                hour: "HH:mm:ss",
                            },
                            style: {
                                colors: "#fff",
                            },
                        },
                    },
                    yaxis: {
                        opposite: true,
                        labels: {
                            style: {
                                colors: "#fff",
                            },
                        },
                    },
                },
            },
            response: {
                uuid: ``,
                station: ``,
                acceleration: [
                    {
                        timestamp: -1,
                        altitude: -1,
                        latitude: -1,
                        longitude: -1,
                        vertical: [],
                        east_west: [],
                        north_south: [],
                    },
                ],
            },
            analysis: {
                vertical: 0,
                east_west: 0,
                north_south: 0,
                synthesis: 0,
            },
        };
    }

    componentDidMount() {
        const socketUrl = getApiUrl({
            tls: AppConfig.backend.tls,
            host: AppConfig.backend.host,
            port: AppConfig.backend.port,
            version: AppConfig.backend.version,
            api: AppConfig.backend.api.socket.uri,
            type: AppConfig.backend.api.socket.type,
        });

        this.setState({
            webSocket: createSocket({
                url: socketUrl,
                type: AppConfig.backend.api.socket.method,
                onMessageCallback: ({ data: result }) => {
                    const data = JSON.parse(result);
                    this.setState({
                        response: data,
                    });
                    this.drawWaveform(data);
                    this.analyseData(data);
                },
            }),
        });
    }

    componentWillUnmount() {
        if (this.state.webSocket) {
            this.state.webSocket.close();
            this.setState({
                webSocket: null,
            });
        }
    }

    fillData(data, key) {
        if (
            this.state.lastDataArr.timestamp < 1 ||
            this.state.lastDataArr.timestamp === data.timestamp
        ) {
            return [];
        }

        const timeDiff = data.timestamp - this.state.lastDataArr.timestamp;
        const timeSpan = timeDiff / data[key].length;

        return data[key].map((item, index) => {
            return [
                new Date(
                    data.timestamp - (data[key].length - index) * timeSpan
                ),
                item,
            ];
        });
    }

    drawWaveform({ acceleration }) {
        this.state.waveform.synthesis[0].data.length > 35 * 300 &&
            this.state.waveform.synthesis[0].data.splice(0, 35);
        this.state.waveform.factors.forEach((_, index) => {
            if (this.state.waveform.factors[index].data.length > 35 * 300) {
                this.state.waveform.factors[index].data.splice(0, 35);
            }
        });

        this.setState({
            waveform: {
                ...this.state.waveform,
                factors: [
                    {
                        ...this.state.waveform.factors[0],
                        data: [
                            ...this.state.waveform.factors[0].data,
                            ...this.fillData(acceleration, "vertical"),
                        ],
                    },
                    {
                        ...this.state.waveform.factors[1],
                        data: [
                            ...this.state.waveform.factors[1].data,
                            ...this.fillData(acceleration, "east_west"),
                        ],
                    },
                    {
                        ...this.state.waveform.factors[2],
                        data: [
                            ...this.state.waveform.factors[2].data,
                            ...this.fillData(acceleration, "north_south"),
                        ],
                    },
                ],
                synthesis: [
                    {
                        ...this.state.waveform.synthesis[0],
                        data: [
                            ...this.state.waveform.synthesis[0].data,
                            ...this.fillData(acceleration, "synthesis"),
                        ],
                    },
                ],
            },
        });
    }

    analyseData = ({ acceleration }) => {
        this.setState({
            analysis: {
                vertical:
                    acceleration.vertical[acceleration.vertical.length - 1],
                east_west:
                    acceleration.east_west[acceleration.east_west.length - 1],
                north_south:
                    acceleration.north_south[
                        acceleration.north_south.length - 1
                    ],
                synthesis:
                    acceleration.synthesis[acceleration.synthesis.length - 1],
            },
            lastDataArr: acceleration,
        });
    };

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4">
                    <Navbar navigation={"实时波形"} />

                    <Notification
                        className={
                            this.state.response.uuid.length !== 0
                                ? `shadow-xl p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-cyan-500 to-yellow-500`
                                : `shadow-xl p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-blue-500 to-orange-500`
                        }
                        icon={
                            this.state.response.uuid.length !== 0 ? (
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 448 512"
                                    className="w-6 h-6 ml-3"
                                    fill="currentColor"
                                >
                                    <path d="M0 64C0 46.3 14.3 32 32 32c229.8 0 416 186.2 416 416c0 17.7-14.3 32-32 32s-32-14.3-32-32C384 253.6 226.4 96 32 96C14.3 96 0 81.7 0 64zM0 416a64 64 0 1 1 128 0A64 64 0 1 1 0 416zM32 160c159.1 0 288 128.9 288 288c0 17.7-14.3 32-32 32s-32-14.3-32-32c0-123.7-100.3-224-224-224c-17.7 0-32-14.3-32-32s14.3-32 32-32z" />
                                </svg>
                            ) : (
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 640 512"
                                    className="w-6 h-6 ml-3"
                                    fill="currentColor"
                                >
                                    <path d="M579.8 267.7c56.5-56.5 56.5-148 0-204.5c-50-50-128.8-56.5-186.3-15.4l-1.6 1.1c-14.4 10.3-17.7 30.3-7.4 44.6s30.3 17.7 44.6 7.4l1.6-1.1c32.1-22.9 76-19.3 103.8 8.6c31.5 31.5 31.5 82.5 0 114L422.3 334.8c-31.5 31.5-82.5 31.5-114 0c-27.9-27.9-31.5-71.8-8.6-103.8l1.1-1.6c10.3-14.4 6.9-34.4-7.4-44.6s-34.4-6.9-44.6 7.4l-1.1 1.6C206.5 251.2 213 330 263 380c56.5 56.5 148 56.5 204.5 0L579.8 267.7zM60.2 244.3c-56.5 56.5-56.5 148 0 204.5c50 50 128.8 56.5 186.3 15.4l1.6-1.1c14.4-10.3 17.7-30.3 7.4-44.6s-30.3-17.7-44.6-7.4l-1.6 1.1c-32.1 22.9-76 19.3-103.8-8.6C74 372 74 321 105.5 289.5L217.7 177.2c31.5-31.5 82.5-31.5 114 0c27.9 27.9 31.5 71.8 8.6 103.9l-1.1 1.6c-10.3 14.4-6.9 34.4 7.4 44.6s34.4 6.9 44.6-7.4l1.1-1.6C433.5 260.8 427 182 377 132c-56.5-56.5-148-56.5-204.5 0L60.2 244.3z" />
                                </svg>
                            )
                        }
                        title={
                            this.state.response.uuid.length !== 0
                                ? `最后更新于 ${getTime(
                                      new Date(
                                          this.state.response.acceleration.timestamp
                                      )
                                  )}`
                                : `暂未收到数据`
                        }
                        text={
                            this.state.response.uuid.length !== 0
                                ? `${this.state.response.station} - ${this.state.response.uuid}`
                                : `正在等待服务器数据...`
                        }
                    />

                    <div className="flex flex-wrap mt-6">
                        <div className="w-full mb-12 xl:mb-0 px-4">
                            <div className="relative flex flex-col w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3  bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                即时
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                实时分量加速度
                                            </h2>
                                        </div>
                                    </div>
                                </div>
                                <div className="p-4 flex-auto shadow-lg bg-gradient-to-tr from-purple-300 to-purple-400 shadow-purple-500/40 rounded-lg">
                                    <div className="relative h-[350px]">
                                        <ReactApexChart
                                            height="350px"
                                            series={this.state.waveform.factors}
                                            options={
                                                this.state.waveform.options
                                            }
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="w-full mb-12 xl:mb-0 px-4">
                            <div className="relative flex flex-col w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3 bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                即时
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                实时合成加速度
                                            </h2>
                                        </div>
                                    </div>
                                </div>
                                <div className="p-4 flex-auto shadow-lg bg-gradient-to-tr from-indigo-300 to-indigo-400 shadow-indigo-500/40 rounded-lg">
                                    <div className="relative h-[350px]">
                                        <ReactApexChart
                                            type="area"
                                            height="350px"
                                            series={
                                                this.state.waveform.synthesis
                                            }
                                            options={
                                                this.state.waveform.options
                                            }
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="w-full px-4">
                            <div className="relative flex flex-col bg-white w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3 bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                数据
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                数据分析
                                            </h2>
                                        </div>
                                    </div>
                                </div>
                                <div className="p-4 shadow-lg flex-auto">
                                    <div className="relative h-[350px]">
                                        <div className="flex flex-wrap -mx-2">
                                            <div className="w-full px-2">
                                                <div className="relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1">
                                                                <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                                    垂直分量当前值
                                                                </h6>
                                                                <h2 className="text-gray-700 text-xl font-semibold">
                                                                    {
                                                                        this
                                                                            .state
                                                                            .analysis
                                                                            .vertical
                                                                    }
                                                                </h2>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>

                                                <div className="relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1">
                                                                <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                                    EW
                                                                    分量当前值
                                                                </h6>
                                                                <h2 className="text-gray-700 text-xl font-semibold">
                                                                    {
                                                                        this
                                                                            .state
                                                                            .analysis
                                                                            .east_west
                                                                    }
                                                                </h2>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>

                                                <div className="relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1">
                                                                <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                                    NS
                                                                    分量当前值
                                                                </h6>
                                                                <h2 className="text-gray-700 text-xl font-semibold">
                                                                    {
                                                                        this
                                                                            .state
                                                                            .analysis
                                                                            .north_south
                                                                    }
                                                                </h2>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>

                                                <div className="relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1">
                                                                <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                                    合成分量当前值
                                                                </h6>
                                                                <h2 className="text-gray-700 text-xl font-semibold">
                                                                    {
                                                                        this
                                                                            .state
                                                                            .analysis
                                                                            .synthesis
                                                                    }
                                                                </h2>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <Scroller />
                <Footer />
            </>
        );
    }
}
