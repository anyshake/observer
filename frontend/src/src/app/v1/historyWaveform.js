import React, { Component } from "react";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Footer from "../../components/Footer";
import Scroller from "../../components/Scroller";
import ReactApexChart from "react-apexcharts";
import "react-datetime/css/react-datetime.css";
import Datetime from "react-datetime";
import getTime from "../../helpers/utilities/getTime";
import {
    errorAlert,
    successAlert,
    timerAlert,
} from "../../helpers/alerts/sweetAlert";
import createRequest from "../../helpers/requests/createRequest";
import AppConfig from "../../config";
import getApiUrl from "../../helpers/utilities/getApiUrl";
import arrSort from "../../helpers/utilities/arrSort";

export default class historyWaveform extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "history",
            timePicker: {
                start: new Date(Date.now() - 60000),
                end: new Date(),
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
                            format: "yy/MM/dd HH:mm:ss",
                        },
                    },
                    xaxis: {
                        type: "datetime",
                        labels: {
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
                data: [
                    {
                        timestamp: -1,
                        latitude: -1,
                        longitude: -1,
                        altitude: -1,
                        vertical: -1,
                        east_west: -1,
                        north_south: -1,
                        synthesis: -1,
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

    fetchData = () => {
        createRequest({
            url: getApiUrl({
                tls: AppConfig.backend.tls,
                host: AppConfig.backend.host,
                port: AppConfig.backend.port,
                version: AppConfig.backend.version,
                api: AppConfig.backend.api.history.uri,
                type: AppConfig.backend.api.history.type,
            }),
            data: {
                start: this.state.timePicker.start.getTime(),
                end: this.state.timePicker.end.getTime(),
            },
            timeout: 300000,
            method: AppConfig.backend.api.history.method,
        })
            .catch(() => {
                errorAlert({
                    title: "查询失败",
                    html: "未找到相关数据，请检查时间范围",
                });
            })
            .then(({ data: { data } }) => {
                successAlert({
                    title: "查询成功",
                    html: `已找到 ${data.length} 条相关数据`,
                });
                this.setState({ response: data });
                this.drawWaveform(data);
                this.analyseData(data);
            });
    };

    drawWaveform(data) {
        const verticalArr = [],
            eastWestArr = [],
            northSouthArr = [],
            synthesisArr = [];

        arrSort(data, "timestamp", "asc");
        data.forEach((item) => {
            verticalArr.push([item.vertical, new Date(item.timestamp)]);
            eastWestArr.push([item.east_west, new Date(item.timestamp)]);
            northSouthArr.push([item.north_south, new Date(item.timestamp)]);
            synthesisArr.push([item.synthesis, new Date(item.timestamp)]);
        });

        // this.setState({
        //     waveform: {
        //         ...this.state.waveform,
        //         factors: [
        //             {
        //                 ...this.state.waveform.factors[0],
        //                 data: verticalArr,
        //             },
        //             {
        //                 ...this.state.waveform.factors[1],
        //                 data: eastWestArr,
        //             },
        //             {
        //                 ...this.state.waveform.factors[2],
        //                 data: northSouthArr,
        //             },
        //         ],
        //         synthesis: [
        //             {
        //                 ...this.state.waveform.synthesis[0],
        //                 data: synthesisArr,
        //             },
        //         ],
        //     },
        // });
    }

    analyseData = (data) => {
        this.setState({
            analysis: {
                vertical: data
                    .map((obj) => obj.vertical)
                    .reduce((max, cur) => (cur > max ? cur : max)),
                east_west: data
                    .map((obj) => obj.east_west)
                    .reduce((max, cur) => (cur > max ? cur : max)),
                north_south: data
                    .map((obj) => obj.north_south)
                    .reduce((max, cur) => (cur > max ? cur : max)),
                synthesis: data
                    .map((obj) => obj.synthesis)
                    .reduce((max, cur) => (cur > max ? cur : max)),
            },
        });
    };

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4">
                    <Navbar navigation={"历史数据"} />

                    <div className="flex flex-wrap mt-6">
                        <div className="w-full xl:w-3/12 px-4">
                            <div className="relative flex flex-col bg-white w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3 bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                回溯
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                历史数据
                                            </h2>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-4 shadow-lg flex-auto text-gray-600 ">
                                    <div className="relative h-[350px]">
                                        <div className="flex flex-wrap -mx-2">
                                            <div className="w-full px-2">
                                                <span className="ml-2">
                                                    起始时间（本地时区）
                                                </span>
                                                <div className="relative flex flex-col min-w-0 break-words w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1 rounded-lg py-2">
                                                                <Datetime
                                                                    inputProps={{
                                                                        className:
                                                                            "w-full",
                                                                        readOnly: true,
                                                                    }}
                                                                    onChange={({
                                                                        _d,
                                                                    }) =>
                                                                        this.setState(
                                                                            {
                                                                                timePicker:
                                                                                    {
                                                                                        ...this
                                                                                            .state
                                                                                            .timePicker,
                                                                                        start: _d,
                                                                                    },
                                                                            }
                                                                        )
                                                                    }
                                                                />
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>

                                            <div className="w-full px-2">
                                                <span className="ml-2">
                                                    结束时间（本地时区）
                                                </span>
                                                <div className="relative flex flex-col min-w-0 break-words w-full mb-4 shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1 rounded-lg py-2">
                                                                <Datetime
                                                                    inputProps={{
                                                                        className:
                                                                            "w-full",
                                                                        readOnly: true,
                                                                    }}
                                                                    onChange={({
                                                                        _d,
                                                                    }) =>
                                                                        this.setState(
                                                                            {
                                                                                timePicker:
                                                                                    {
                                                                                        ...this
                                                                                            .state
                                                                                            .timePicker,
                                                                                        end: _d,
                                                                                    },
                                                                            }
                                                                        )
                                                                    }
                                                                />
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>

                                            <div className="w-full px-2 text mt-4 ml-2 text-sm">
                                                {`起始于 ${getTime(
                                                    this.state.timePicker.start
                                                )}`}
                                                <br />
                                                {`截止于 ${getTime(
                                                    this.state.timePicker.end
                                                )}`}
                                            </div>
                                        </div>

                                        <button
                                            onClick={() => {
                                                this.fetchData();
                                                timerAlert({
                                                    title: "查询中",
                                                    html: "查询速度取决于数据量，请耐心等待一到五分钟左右",
                                                    timer: 300000,
                                                    loading: false,
                                                    callback: () => {
                                                        errorAlert({
                                                            title: "查询失败",
                                                            text: "请求接口超时，请尝试缩小查询范围再试",
                                                        });
                                                    },
                                                });
                                            }}
                                            className="absolute w-full mt-9 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none font-medium rounded-lg text-sm px-5 py-2.5 text-center"
                                        >
                                            送出查询
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="w-full xl:w-9/12 xl:mb-0 px-4">
                            <div className="relative flex flex-col w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3  bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                历史
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                三分量加速度
                                            </h2>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-4 flex-auto shadow-lg bg-gradient-to-tr from-pink-300 to-pink-400 shadow-pink-500/40 rounded-lg">
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
                    </div>

                    <div className="flex flex-wrap mt-6">
                        <div className="w-full mb-12 xl:w-10/12 xl:mb-0 px-4">
                            <div className="relative flex flex-col w-full mb-6 shadow-lg rounded-lg">
                                <div className="px-4 py-3  bg-transparent">
                                    <div className="flex flex-wrap items-center">
                                        <div className="relative w-full max-w-full flex-grow flex-1">
                                            <h6 className="text-gray-500 mb-1 text-xs font-semibold">
                                                历史
                                            </h6>
                                            <h2 className="text-gray-700 text-xl font-semibold">
                                                合成加速度
                                            </h2>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-4 flex-auto shadow-lg bg-gradient-to-tr from-teal-300 to-teal-400 shadow-teal-500/40 rounded-lg">
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

                        <div className="w-full xl:w-2/12 px-4">
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
                                                                    垂直分量峰值
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
                                                                    EW 分量峰值
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
                                                                    NS 分量峰值
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
                                                                    合成分量峰值
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
