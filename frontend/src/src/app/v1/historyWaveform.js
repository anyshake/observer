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
    selectAlert,
    successAlert,
    timerAlert,
} from "../../helpers/alerts/sweetAlert";
import createRequest from "../../helpers/requests/createRequest";
import AppConfig from "../../config";
import getApiUrl from "../../helpers/utilities/getApiUrl";
import arrSort from "../../helpers/utilities/arrSort";
import arrMaximum from "../../helpers/utilities/arrMaximum";
import Modal from "react-responsive-modal";
import "react-responsive-modal/styles.css";

const FETCH_TIMEOUT = 30 * 1000;

export default class historyWaveform extends Component {
    constructor(props) {
        super(props);
        this.state = {
            sidebarMark: "history",
            timePicker: new Date(Date.now() - 120000),
            showModal: false,
            cardLimit: 3,
            waveform: {
                factors: [
                    {
                        name: "垂直分量",
                        color: "#5a3eba",
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
                        animations: {
                            enabled: false,
                        },
                    },
                    dataLabels: {
                        enabled: false,
                    },
                    legend: {
                        show: true,
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
            trace: {
                source: [],
                list: [],
            },
        };
    }

    queryHistoryData = () => {
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
            timeout: FETCH_TIMEOUT,
            method: AppConfig.backend.api.history.method,
        })
            .then(({ data: { data } }) => {
                successAlert({
                    title: "查询成功",
                    html: `已找到 ${data.length} 条相关数据`,
                });
                this.setState({
                    response: data,
                });
                this.drawWaveform(data);
                this.analyseData(data);
            })
            .catch(() => {
                errorAlert({
                    title: "查询失败",
                    html: "未找到相关数据，请检查时间范围",
                });
            });

        timerAlert({
            title: "查询中",
            html: "正在查询加速度数据，这可能需要一些时间",
            timer: 30000,
            loading: true,
            callback: () => {
                errorAlert({
                    title: "查询失败",
                    text: "请求接口超时，请尝试缩小查询范围再试",
                });
            },
        });
    };

    fillData(dataArr, timestampArr) {
        const arr = dataArr.map((item, index) => {
            const timeDiff =
                index !== dataArr.length - 1
                    ? timestampArr[index + 1] - timestampArr[index]
                    : 1000;
            const timeSpan = timeDiff / item.length;

            return item.map((obj, _index) => [
                new Date(timestampArr[index] + _index * timeSpan),
                obj,
            ]);
        });

        return arr.flat();
    }

    filterData = (arr, key) => {
        return arrMaximum(arr.map((obj) => obj[key]).flat());
    };

    drawWaveform(data) {
        arrSort(data, "timestamp", "asc");

        this.setState({
            waveform: {
                ...this.state.waveform,
                factors: [
                    {
                        ...this.state.waveform.factors[0],
                        data: [
                            ...this.fillData(
                                data.map((obj) => obj.vertical),
                                data.map((obj) => obj.timestamp)
                            ),
                        ],
                    },
                    {
                        ...this.state.waveform.factors[1],
                        data: [
                            ...this.fillData(
                                data.map((obj) => obj.east_west),
                                data.map((obj) => obj.timestamp)
                            ),
                        ],
                    },
                    {
                        ...this.state.waveform.factors[2],
                        data: [
                            ...this.fillData(
                                data.map((obj) => obj.north_south),
                                data.map((obj) => obj.timestamp)
                            ),
                        ],
                    },
                ],
                synthesis: [
                    {
                        ...this.state.waveform.synthesis[3],
                        data: [
                            ...this.fillData(
                                data.map((obj) => obj.synthesis),
                                data.map((obj) => obj.timestamp)
                            ),
                        ],
                    },
                ],
            },
        });
    }

    analyseData = (data) => {
        this.setState({
            analysis: {
                vertical: this.filterData(data, "vertical"),
                east_west: this.filterData(data, "east_west"),
                north_south: this.filterData(data, "north_south"),
                synthesis: this.filterData(data, "synthesis"),
            },
        });
    };

    render() {
        return (
            <>
                <Sidebar sidebarMark={this.state.sidebarMark} />
                <div className="content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4">
                    <Navbar navigation={"历史数据"} />

                    {this.state.showModal && (
                        <Modal
                            open={this.state.showModal}
                            onClose={() =>
                                this.setState({
                                    showModal: false,
                                })
                            }
                            classNames={{
                                modal: "rounded-lg w-[80%] md:w-[60%] lg:w-[50%]",
                            }}
                        >
                            <h2 className="text-xl font-bold">选择一个地震</h2>
                            <h4 className="text-md mt-2">
                                点击卡片可以查看对应时刻测站波形
                            </h4>
                            <section className="overflow-hidden text-gray-700 mt-4">
                                <div className="container mx-auto px-5 py-2">
                                    <div className="flex flex-wrap items-center justify-center gap-2">
                                        {this.state.trace.list.map(
                                            (item, index) =>
                                                this.state.cardLimit >
                                                    index && (
                                                    <div
                                                        key={index}
                                                        className="cursor-pointer w-full border-2 border-b-4 border-gray-200 rounded-xl hover:bg-gray-50"
                                                        onClick={() => {
                                                            this.setState({
                                                                showModal: false,
                                                                timePicker:
                                                                    new Date(
                                                                        item.timestamp +
                                                                            item.estimated *
                                                                                1000
                                                                    ),
                                                            });

                                                            setTimeout(
                                                                () =>
                                                                    this.queryHistoryData(),
                                                                100
                                                            );
                                                        }}
                                                    >
                                                        <div className="p-5 flex flex-col">
                                                            <p className="text-sky-500 font-bold text-xs">
                                                                {`${item.region}`}
                                                            </p>
                                                            <p className="text-sky-500 font-bold text-xs">
                                                                {`${getTime(
                                                                    new Date(
                                                                        item.timestamp
                                                                    )
                                                                )}`}
                                                            </p>
                                                            <p className="text-gray-600 font-bold">
                                                                {`${item.magnitude} 级 / ${item.event}`}
                                                            </p>
                                                            <p className="text-gray-400 text-sm">
                                                                {`震源深度：${
                                                                    item.depth !==
                                                                    -1
                                                                        ? `${item.depth} km`
                                                                        : `数据源未提供`
                                                                }`}
                                                            </p>
                                                            <p className="text-gray-400 text-sm">
                                                                {`距离测站：${item.distance.toFixed(
                                                                    2
                                                                )} km`}
                                                            </p>
                                                            <p className="text-gray-400 text-sm">
                                                                {`传播时长：${item.estimated.toFixed(
                                                                    2
                                                                )} s`}
                                                            </p>
                                                        </div>
                                                    </div>
                                                )
                                        )}
                                        {this.state.cardLimit <
                                        this.state.trace.list.length ? (
                                            <button
                                                className="mt-6 bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-lg"
                                                onClick={() =>
                                                    this.setState({
                                                        cardLimit:
                                                            this.state
                                                                .cardLimit + 6,
                                                    })
                                                }
                                            >
                                                加载更多
                                            </button>
                                        ) : (
                                            <span className="mt-6 py-2 px-4">
                                                没有更多了
                                            </span>
                                        )}
                                    </div>
                                </div>
                            </section>
                        </Modal>
                    )}

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
                                    <div className="h-[350px] flex flex-col space-y-10">
                                        <div className="flex flex-wrap -mx-2">
                                            <div className="w-full px-2">
                                                <span className="ml-2">
                                                    {`查询时间（时区 ${
                                                        Intl.DateTimeFormat().resolvedOptions()
                                                            .timeZone
                                                    }）`}
                                                </span>
                                                <div className="flex flex-col min-w-0 break-words w-full shadow-lg rounded-lg mt-4">
                                                    <div className="px-4 py-3 bg-transparent">
                                                        <div className="flex flex-wrap items-center">
                                                            <div className="w-full max-w-full flex-grow flex-1 rounded-lg py-2">
                                                                <Datetime
                                                                    dateFormat="YYYY-MM-DD"
                                                                    timeFormat="HH:mm:ss"
                                                                    inputProps={{
                                                                        readOnly: true,
                                                                        className:
                                                                            "w-full cursor-pointer focus:outline-none rounded-lg",
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

                                        <div className="w-full px-2 text ml-1">
                                            {`查询时间 ${getTime(
                                                this.state.timePicker
                                            )}`}
                                            <br />
                                            系统将查询 5 分钟内震动波形
                                        </div>

                                        <div className="flex flex-col justify-center items-center gap-4 font-medium text-sm">
                                            <button
                                                onClick={() => {
                                                    this.queryHistoryData();
                                                }}
                                                className="w-full text-white shadow-lg bg-purple-700 hover:bg-purple-800 focus:ring-4 focus:outline-none rounded-lg text-sm px-5 py-2.5 text-center"
                                            >
                                                送出查询
                                            </button>
                                            <button
                                                onClick={() => {
                                                    createRequest({
                                                        url: getApiUrl({
                                                            tls: AppConfig
                                                                .backend.tls,
                                                            host: AppConfig
                                                                .backend.host,
                                                            port: AppConfig
                                                                .backend.port,
                                                            version:
                                                                AppConfig
                                                                    .backend
                                                                    .version,
                                                            api: AppConfig
                                                                .backend.api
                                                                .trace.uri,
                                                            type: AppConfig
                                                                .backend.api
                                                                .trace.type,
                                                        }),
                                                        data: {
                                                            source: "show",
                                                        },
                                                        timeout: FETCH_TIMEOUT,
                                                        method: AppConfig
                                                            .backend.api.trace
                                                            .method,
                                                    })
                                                        .then(
                                                            ({
                                                                data: { data },
                                                            }) => {
                                                                this.setState({
                                                                    trace: {
                                                                        source: data,
                                                                    },
                                                                });
                                                                selectAlert({
                                                                    title: "请选择",
                                                                    html: "请选择一个地震数据源",
                                                                    input: "select",
                                                                    inputOptions:
                                                                        data
                                                                            .map(
                                                                                (
                                                                                    item
                                                                                ) => ({
                                                                                    [item.value]:
                                                                                        item.name,
                                                                                })
                                                                            )
                                                                            .reduce(
                                                                                (
                                                                                    obj,
                                                                                    curr
                                                                                ) =>
                                                                                    Object.assign(
                                                                                        obj,
                                                                                        curr
                                                                                    ),
                                                                                {}
                                                                            ),
                                                                    callback: (
                                                                        e
                                                                    ) => {
                                                                        createRequest(
                                                                            {
                                                                                url: getApiUrl(
                                                                                    {
                                                                                        tls: AppConfig
                                                                                            .backend
                                                                                            .tls,
                                                                                        host: AppConfig
                                                                                            .backend
                                                                                            .host,
                                                                                        port: AppConfig
                                                                                            .backend
                                                                                            .port,
                                                                                        version:
                                                                                            AppConfig
                                                                                                .backend
                                                                                                .version,
                                                                                        api: AppConfig
                                                                                            .backend
                                                                                            .api
                                                                                            .trace
                                                                                            .uri,
                                                                                        type: AppConfig
                                                                                            .backend
                                                                                            .api
                                                                                            .trace
                                                                                            .type,
                                                                                    }
                                                                                ),
                                                                                data: {
                                                                                    source: e,
                                                                                },
                                                                                timeout:
                                                                                    FETCH_TIMEOUT,
                                                                                method: AppConfig
                                                                                    .backend
                                                                                    .api
                                                                                    .trace
                                                                                    .method,
                                                                            }
                                                                        )
                                                                            .then(
                                                                                ({
                                                                                    data: {
                                                                                        data,
                                                                                    },
                                                                                }) => {
                                                                                    successAlert(
                                                                                        {
                                                                                            title: "查询成功",
                                                                                            html: `已找到 ${data.length} 条相关数据，按下确认继续`,
                                                                                        }
                                                                                    ).then(
                                                                                        (
                                                                                            result
                                                                                        ) => {
                                                                                            if (
                                                                                                result.isConfirmed
                                                                                            ) {
                                                                                                arrSort(
                                                                                                    data,
                                                                                                    "timestamp",
                                                                                                    "desc"
                                                                                                );
                                                                                                this.setState(
                                                                                                    {
                                                                                                        showModal: true,
                                                                                                        trace: {
                                                                                                            list: data,
                                                                                                        },
                                                                                                    }
                                                                                                );
                                                                                            }
                                                                                        }
                                                                                    );
                                                                                }
                                                                            )
                                                                            .catch(
                                                                                () => {
                                                                                    errorAlert(
                                                                                        {
                                                                                            title: "查询失败",
                                                                                            html: "未找到相关数据，请稍后重试",
                                                                                        }
                                                                                    );
                                                                                }
                                                                            );
                                                                        timerAlert(
                                                                            {
                                                                                title: "查询中",
                                                                                html: "正在请求地震列表数据，这可能需要一些时间",
                                                                                timer: 30000,
                                                                                loading: true,
                                                                                callback:
                                                                                    () => {
                                                                                        errorAlert(
                                                                                            {
                                                                                                title: "查询失败",
                                                                                                text: "请求接口超时，请稍候重试",
                                                                                            }
                                                                                        );
                                                                                    },
                                                                            }
                                                                        );
                                                                    },
                                                                });
                                                            }
                                                        )
                                                        .catch(() => {
                                                            errorAlert({
                                                                title: "查询失败",
                                                                html: "未找到相关数据，请稍后重试",
                                                            });
                                                        });
                                                    timerAlert({
                                                        title: "查询中",
                                                        html: "正在获取地震数据源，这可能需要一些时间",
                                                        timer: 30000,
                                                        loading: true,
                                                        callback: () => {
                                                            errorAlert({
                                                                title: "查询失败",
                                                                text: "请求接口超时，请稍后再试",
                                                            });
                                                        },
                                                    });
                                                }}
                                                className="w-full text-white shadow-lg bg-indigo-700 hover:bg-indigo-800 focus:ring-4 focus:outline-none rounded-lg px-5 py-2.5 text-center"
                                            >
                                                地震反查
                                            </button>
                                        </div>
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
