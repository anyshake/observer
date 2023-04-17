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
            timePicker: new Date(Date.now() - 60000),
            waveform: {
                factors: [
                    {
                        name: "垂直分量",
                        color: "#d97706",
                        data: [],
                    },
                    {
                        name: "水平 EW",
                        color: "#128672",
                        data: [],
                    },
                    {
                        name: "水平 NS",
                        color: "#c3268a",
                        data: [],
                    },
                ],
                synthesis: [
                    {
                        name: "合成分量",
                        color: "#cf4500",
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
                            show: true,
                        },
                        zoom: {
                            enabled: true,
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
                            datetimeFormatter: {
                                hour: "HH:mm:ss",
                            },
                            style: {
                                colors: "#fff",
                            },
                        },
                    },
                    yaxis: {
                        tickAmount: 5,
                        opposite: true,
                        labels: {
                            style: {
                                colors: "#fff",
                            },
                        },
                    },
                    brush: {
                        enabled: true,
                        target: "main",
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
                timestamp: this.state.timePicker.getTime(),
            },
            timeout: 60000,
            method: AppConfig.backend.api.history.method,
        })
            .then(({ data: { data } }) => {
                successAlert({
                    title: "查询成功",
                    html: `已找到 ${data.length} 条相关数据`,
                });
                this.setState({ response: data });
                this.drawWaveform(data);
                this.analyseData(data);
            })
            .catch(() => {
                errorAlert({
                    title: "查询失败",
                    html: "未找到相关数据，请检查时间范围",
                });
            });
    };

    drawWaveform(data) {
        const verticalArr = [],
            eastWestArr = [],
            northSouthArr = [],
            synthesisArr = [];

        arrSort(data, "timestamp", "asc");
        data.forEach((item) => {
            verticalArr.push([new Date(item.timestamp), item.vertical]);
            eastWestArr.push([new Date(item.timestamp), item.east_west]);
            northSouthArr.push([new Date(item.timestamp), item.north_south]);
            synthesisArr.push([new Date(item.timestamp), item.synthesis]);
        });

        this.setState({
            waveform: {
                ...this.state.waveform,
                factors: [
                    {
                        ...this.state.waveform.factors[0],
                        data: [...verticalArr],
                    },
                    {
                        ...this.state.waveform.factors[1],
                        data: [...eastWestArr],
                    },
                    {
                        ...this.state.waveform.factors[2],
                        data: [...northSouthArr],
                    },
                ],
                synthesis: [
                    {
                        ...this.state.waveform.synthesis[3],
                        data: [...synthesisArr],
                    },
                ],
            },
        });
    }

    analyseData = (data) => {
        this.setState({
            analysis: {
                vertical: data
                    .map((obj) => Math.abs(obj.vertical))
                    .reduce((max, cur) => (cur > max ? cur : max)),
                east_west: data
                    .map((obj) => Math.abs(obj.east_west))
                    .reduce((max, cur) => (cur > max ? cur : max)),
                north_south: data
                    .map((obj) => Math.abs(obj.north_south))
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
                                            <div className="w-full px-2 mt-8">
                                                <span className="ml-2">
                                                    查询时间（本地时区）
                                                </span>
                                                <div className="relative flex flex-col min-w-0 break-words w-full shadow-lg rounded-lg">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="relative w-full max-w-full flex-grow flex-1 rounded-lg py-2">
                                                                <Datetime
                                                                    dateFormat="YYYY-MM-DD"
                                                                    timeFormat="HH:mm:ss"
                                                                    inputProps={{
                                                                        className:
                                                                            "w-full cursor-pointer focus:outline-none rounded-lg",
                                                                        readOnly: true,
                                                                        placeholder: `点击选择时间`,
                                                                    }}
                                                                    onChange={({
                                                                        _d,
                                                                    }) =>
                                                                        this.setState(
                                                                            {
                                                                                timePicker:
                                                                                    _d,
                                                                            }
                                                                        )
                                                                    }
                                                                />
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="absolute w-full px-2 text mt-20 ml-1">
                                            {`起始于 ${getTime(
                                                this.state.timePicker
                                            )}`}
                                            <br />
                                            {`系统将查询 1 分钟内的波形`}
                                        </div>

                                        <button
                                            onClick={() => {
                                                this.fetchData();
                                                timerAlert({
                                                    title: "查询中",
                                                    html: "查询速度取决于上位机性能，请耐心等待一到两分钟左右",
                                                    timer: 120000,
                                                    loading: false,
                                                    callback: () => {
                                                        errorAlert({
                                                            title: "查询失败",
                                                            text: "请求接口超时，请尝试缩小查询范围再试",
                                                        });
                                                    },
                                                });
                                            }}
                                            className="absolute w-full mt-40 text-white shadow-lg bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none font-medium rounded-lg text-sm px-5 py-2.5 text-center"
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
                        <div className="w-full xl:w-3/12 px-4">
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
                                                                    垂直分量峰值（绝对值）
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
                                                                    分量峰值（绝对值）
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
                                                                    分量峰值（绝对值）
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

                        <div className="w-full xl:w-9/12 px-4">
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

                                <div className="p-4 flex-auto shadow-lg bg-gradient-to-tr from-orange-300 to-orange-400 shadow-orange-500/40 rounded-lg">
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
                    </div>
                </div>

                <Scroller />
                <Footer />
            </>
        );
    }
}
