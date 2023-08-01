import { Component } from "react";
import Header from "../../components/Header";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Content from "../../components/Content";
import View from "../../components/View";
import Scroller from "../../components/Scroller";
import Footer from "../../components/Footer";
import TimePicker from "../../components/TimePicker";
import Card from "../../components/Card";
import Button from "../../components/Button";
import Container from "../../components/Container";
import Chart, { ChartProps } from "../../components/Chart";
import requestByTag from "../../helpers/requestByTag";
import toast, { Toaster } from "react-hot-toast";
import setChart from "./setChart";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import setGeophone from "./setGeophone";
import setADC from "./setADC";
import ModalDialog, { ModalDialogProps } from "../../components/ModalDialog";
import getTimeString from "../../helpers/getTimeString";
import Label, { LabelProps } from "../../components/Label";
import setLabels from "./setLabels";
import { IntensityScaleStandard } from "../../helpers/getIntensity";
import getLocalStorage from "../../helpers/getLocalStorage";

// 100s by default
const QUERY_TIMEOUT = 100000;
const TRACE_RANGE = 1000 * 5 * 60;

interface HistoryForm {
    readonly start: number;
    readonly end: number;
    readonly format?: "json" | "sac";
    readonly channel?: "EHZ" | "EHE" | "EHN";
    [key: string]: string | number | undefined;
}

interface TraceForm {
    readonly source: string;
    [key: string]: string | number | undefined;
}

interface HistorySelect {
    readonly from: string;
    readonly dialog: SelectDialogProps;
}

interface State {
    readonly adc: ADC;
    readonly history: HistoryForm;
    readonly trace: TraceForm;
    readonly chart: ChartProps;
    readonly geophone: Geophone;
    readonly select: HistorySelect;
    readonly modal: ModalDialogProps;
    readonly scale: IntensityScaleStandard;
    readonly labels: LabelProps[];
}

export default class History extends Component<{}, State> {
    constructor(props: {}) {
        super(props);
        this.state = {
            scale: "JMA",
            trace: {
                source: "show",
            },
            history: {
                start: Date.now() - 60000,
                end: Date.now(),
                format: "json",
                channel: "EHZ",
            },
            chart: {
                backgroundColor: "transparent",
                tickInterval: 0.1,
                tickPrecision: 0.2,
                lineWidth: 1,
                height: 400,
                tooltip: true,
                legend: true,
                zooming: true,
                series: [
                    {
                        type: "line",
                        name: "EHZ",
                        color: "#5a3eba",
                        data: [],
                    },
                    {
                        type: "line",
                        name: "EHE",
                        color: "#128672",
                        data: [],
                    },
                    {
                        type: "line",
                        name: "EHN",
                        color: "#c3268a",
                        data: [],
                    },
                ],
            },
            select: {
                from: "history",
                dialog: {
                    open: false,
                    title: "选择要导出的通道",
                    values: [
                        ["垂直通道", "EHZ"],
                        ["水平东西", "EHE"],
                        ["水平南北", "EHN"],
                    ],
                },
            },
            modal: {
                open: false,
                values: [],
                title: "选择一个事件",
            },
            geophone: {
                ehz: 0.288,
                ehe: 0.288,
                ehn: 0.288,
            },
            adc: {
                fullscale: 5,
                resolution: 24,
            },
            labels: [
                {
                    tag: "ehz-pga",
                    label: "EHZ 峰值加速度",
                    value: "0",
                    unit: "gal",
                },
                {
                    tag: "ehz-intensity",
                    label: "EHZ 瞬时最大震度",
                    value: "0",
                    unit: "",
                },
                {
                    tag: "ehe-pga",
                    label: "EHZ 峰值加速度",
                    value: "0",
                    unit: "gal",
                },
                {
                    tag: "ehe-intensity",
                    label: "EHE 瞬时最大震度",
                    value: "0",
                    unit: "",
                },
                {
                    tag: "ehn-pga",
                    label: "EHN 峰值加速度",
                    value: "0",
                    unit: "gal",
                },
                {
                    tag: "ehn-intensity",
                    label: "EHN 瞬时最大震度",
                    value: "0",
                    unit: "",
                },
            ],
        };
    }

    async componentDidMount(): Promise<void> {
        const res = await requestByTag({
            tag: "station",
        });
        if (res.data) {
            const adc = setADC(res);
            const geophone = setGeophone(res);
            const scale = getLocalStorage(
                "scale",
                "JMA"
            ) as IntensityScaleStandard;
            this.setState({ adc, geophone, scale });
        } else {
            return;
        }
    }

    promisedSetState = (newState: any) =>
        new Promise((resolve: any) => this.setState(newState, resolve));

    handleTimeChange = (type: string, value: number): void => {
        switch (type) {
            case "start":
                this.setState((state) => ({
                    history: { ...state.history, start: value },
                }));
                break;
            case "end":
                this.setState((state) => ({
                    history: { ...state.history, end: value },
                }));
                break;
        }
    };

    handleQueryHistory = async (): Promise<unknown> => {
        const { history } = this.state;
        const { start, end } = history;

        if (end - start <= 0 || !start || !end) {
            const error = "请选择正确的时间范围";
            toast.error(error);
            return Promise.reject(error);
        }

        const { error, data } = await requestByTag({
            body: history,
            tag: "history",
            timeout: QUERY_TIMEOUT,
            blob: history.format === "sac",
            filename: `${history.channel}-${history.start}-${history.end}.${history.format}`,
        });
        if (error) {
            const error = "请求失败，请检查输入后重试";
            toast.error(error);
            return Promise.reject(error);
        }

        return data;
    };

    handleQueryEvents = async (): Promise<unknown> => {
        const { trace } = this.state;
        const { error, data } = await requestByTag({
            body: trace,
            tag: "trace",
            timeout: QUERY_TIMEOUT,
        });
        if (error) {
            const error = "请求失败，请检查输入后重试";
            toast.error(error);
            return Promise.reject(error);
        }

        this.setState((state) => ({
            modal: {
                ...state.modal,
                open: true,
                values: data.map((item: any) => {
                    const {
                        magnitude,
                        region,
                        event,
                        timestamp,
                        depth,
                        estimated,
                    } = item;
                    const desc = `[M ${magnitude.toFixed(
                        1
                    )}] ${region} / 时刻 ${getTimeString(
                        timestamp
                    )} / 深度 ${depth} km / 传播 ${estimated.toFixed(1)} s`;

                    return [event, timestamp + estimated * 1000, desc];
                }),
            },
        }));
    };

    handleChooseEvent = async (value: string): Promise<void> => {
        const span = TRACE_RANGE / 2;
        const start = new Date(value).getTime() - span;
        const end = new Date(value).getTime() + span;

        await this.promisedSetState({
            history: { start, end, format: "json" },
            modal: { ...this.state.modal, open: false },
        });
    };

    handleSelect = async (from: string, value: string): Promise<void> => {
        const select = {
            from: "history",
            dialog: {
                ...this.state.select.dialog,
                open: false,
            },
        };

        switch (from) {
            case "history":
                await this.promisedSetState({
                    select,
                    history: {
                        ...this.state.history,
                        channel: value,
                        format: "sac",
                    },
                });
                await toast.promise(this.handleQueryHistory(), {
                    loading: "正在查询...",
                    success: "历史波形数据导出成功",
                    error: "历史波形数据导出失败",
                });
                break;
            case "trace":
                await this.promisedSetState({
                    select,
                    trace: {
                        ...this.state.trace,
                        source: value,
                    },
                });
                await toast.promise(this.handleQueryEvents(), {
                    loading: "正在查询...",
                    success: "地震事件查询成功",
                    error: "地震事件查询失败",
                });
                break;
        }
    };

    handleQueryWaveform = async (): Promise<void> => {
        await this.promisedSetState({
            history: { ...this.state.history, format: "json" },
        });

        const res = await toast.promise(this.handleQueryHistory(), {
            loading: "正在查询...",
            success: "历史波形数据查询成功",
            error: "历史波形数据查询失败",
        });
        if (res) {
            const { adc, geophone, scale } = this.state;
            const chart = setChart(this.state.chart, res, adc, geophone);
            const labels = setLabels(
                this.state.labels,
                res,
                adc,
                geophone,
                scale
            );
            this.setState({ chart, labels });
        }
    };

    handleQuerySACFile = async (): Promise<void> => {
        this.setState((state) => ({
            select: {
                from: "history",
                dialog: {
                    ...state.select,
                    title: "选择要导出的通道",
                    values: [
                        ["垂直通道", "EHZ"],
                        ["水平东西", "EHE"],
                        ["水平南北", "EHN"],
                    ],
                    open: true,
                },
            },
        }));
    };

    handleQuerySource = async (): Promise<unknown> => {
        const trace = {
            source: "show",
        };

        const { data, error } = await toast.promise(
            requestByTag({
                body: trace,
                tag: "trace",
            }),
            {
                loading: "正在获取数据源...",
                success: "成功取得数据源",
                error: "数据源获取失败",
            }
        );

        if (error || !data) {
            const error = "请求失败，请检查输入后重试";
            toast.error(error);
            return Promise.reject(error);
        }

        this.setState((state) => ({
            select: {
                from: "trace",
                dialog: {
                    ...state.select.dialog,
                    open: true,
                    title: "选择要地震数据来源",
                    values: data.map((item: any) => [item.name, item.value]),
                },
            },
        }));
    };

    render() {
        const { chart, select, modal, history, labels } = this.state;
        const { from, dialog } = select;
        const { start, end } = history;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />
                    <Container layout="grid">
                        <Card className="h-[430px]" label="历史查询">
                            <TimePicker
                                value={start}
                                label="选择起始时间"
                                onChange={(value) =>
                                    this.handleTimeChange("start", value)
                                }
                            />
                            <TimePicker
                                value={end}
                                label="选择结束时间"
                                onChange={(value) =>
                                    this.handleTimeChange("end", value)
                                }
                            />

                            <Button
                                className="mt-6 bg-indigo-700 hover:bg-indigo-800"
                                onClick={this.handleQueryWaveform}
                                label="调阅波形"
                            />
                            <Button
                                className="bg-green-700 hover:bg-green-800"
                                onClick={this.handleQuerySACFile}
                                label="数据下载"
                            />
                            <Button
                                className="bg-yellow-700 hover:bg-yellow-800"
                                onClick={this.handleQuerySource}
                                label="事件反查"
                            />
                        </Card>

                        <Card
                            className="h-[430px] rounded-lg bg-pink-300"
                            label="加速度波形图"
                        >
                            <Chart {...chart} />
                        </Card>
                    </Container>

                    <Card className="rounded-lg" label="数据分析">
                        <Container layout="grid">
                            {labels.map((label, index) => (
                                <Label key={index} {...label} />
                            ))}
                        </Container>
                    </Card>
                </Content>

                <Scroller />
                <Footer />

                <SelectDialog
                    {...dialog}
                    onSelect={(value: string) => this.handleSelect(from, value)}
                />
                <ModalDialog
                    {...modal}
                    onSelect={this.handleChooseEvent}
                    onClose={() =>
                        this.setState({ modal: { ...modal, open: false } })
                    }
                />
                <Toaster position="top-center" />
            </View>
        );
    }
}
