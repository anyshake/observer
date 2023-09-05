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
import restfulApiByTag from "../../helpers/request/restfulApiByTag";
import toast, { Toaster } from "react-hot-toast";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import setGeophone from "./setGeophone";
import setADC from "./setADC";
import ModalDialog, { ModalDialogProps } from "../../components/ModalDialog";
import getTimeString from "../../helpers/utils/getTimeString";
import Label, { LabelProps } from "../../components/Label";
import setLabels from "./setLabels";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import { fallbackScale } from "../../config/global";
import { ReduxStoreProps } from "../../config/store";
import { update as updateADC } from "../../store/adc";
import { update as updateGeophone } from "../../store/geophone";
import { connect } from "react-redux";
import mapStateToProps from "../../helpers/utils/mapStateToProps";
import Area, { AreaProps } from "../../components/Area";
import setAreas from "./setAreas";
import { WithTranslation, withTranslation } from "react-i18next";

// Query duration is 100s by default
const TRACE_RANGE = 1000 * 5 * 60;
// Query timeout is 100s by default
const QUERY_TIMEOUT = 100 * 1000;

export interface HistoryArea {
    readonly tag: string;
    readonly area: AreaProps;
    readonly chart: ChartProps;
}

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

interface HistoryState {
    readonly adc: ADC;
    readonly labels: LabelProps[];
    readonly history: HistoryForm;
    readonly trace: TraceForm;
    readonly areas: HistoryArea[];
    readonly geophone: Geophone;
    readonly select: HistorySelect;
    readonly modal: ModalDialogProps;
    readonly scale: IntensityStandardProperty;
}

class History extends Component<
    ReduxStoreProps & WithTranslation,
    HistoryState
> {
    constructor(props: ReduxStoreProps & WithTranslation) {
        super(props);
        this.state = {
            trace: {
                source: "show",
            },
            history: {
                start: Date.now() - 60000,
                end: Date.now(),
                format: "json",
                channel: "EHZ",
            },
            areas: [
                {
                    tag: "ehz",
                    area: {
                        label: { id: "views.history.areas.ehz.label" },
                    },
                    chart: {
                        backgroundColor: "#d97706",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHZ",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
                {
                    tag: "ehe",
                    area: {
                        label: { id: "views.history.areas.ehe.label" },
                    },
                    chart: {
                        backgroundColor: "#10b981",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHE",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
                {
                    tag: "ehn",
                    area: {
                        label: { id: "views.history.areas.ehn.label" },
                    },
                    chart: {
                        backgroundColor: "#0ea5e9",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHN",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
            ],
            select: {
                from: "history",
                dialog: {
                    open: false,
                    title: { id: "views.history.selects.choose_channel.title" },
                    values: [
                        ["Vertical", "EHZ"],
                        ["East-West", "EHE"],
                        ["North-South", "EHN"],
                    ],
                },
            },
            modal: {
                open: false,
                values: [],
                title: { id: "views.history.modals.choose_event.title" },
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
                    label: { id: "views.history.labels.ehz_pga.label" },
                    unit: { id: "views.history.labels.ehz_pga.unit" },
                    value: "0",
                },
                {
                    tag: "ehz-intensity",
                    label: { id: "views.history.labels.ehz_intensity.label" },
                    unit: { id: "views.history.labels.ehz_intensity.unit" },
                    value: "Unknown",
                },
                {
                    tag: "ehe-pga",
                    label: { id: "views.history.labels.ehe_pga.label" },
                    unit: { id: "views.history.labels.ehe_pga.unit" },
                    value: "0",
                },
                {
                    tag: "ehe-intensity",
                    label: { id: "views.history.labels.ehe_intensity.label" },
                    unit: { id: "views.history.labels.ehe_intensity.unit" },
                    value: "Unknown",
                },
                {
                    tag: "ehn-pga",
                    label: { id: "views.history.labels.ehn_pga.label" },
                    unit: { id: "views.history.labels.ehn_pga.unit" },
                    value: "0",
                },
                {
                    tag: "ehn-intensity",
                    label: { id: "views.history.labels.ehn_intensity.label" },
                    unit: { id: "views.history.labels.ehn_intensity.unit" },
                    value: "Unknown",
                },
            ],
            scale: fallbackScale.property(),
        };
    }

    async componentDidMount(): Promise<void> {
        // Get ADC, Geophone, scale standard from redux
        let { adc } = this.props.adc;
        const { scale } = this.props.scale;
        let { geophone } = this.props.geophone;
        const { ehz, ehe, ehn } = geophone;
        const { resolution } = adc;

        // Query again from server if value is not set
        if (resolution === -1 || ehz * ehe * ehn === 0) {
            const res = await restfulApiByTag({
                tag: "station",
            });
            if (res.data) {
                // Get new formatted state
                geophone = setGeophone(res);
                adc = setADC(res);
                // Apply to Redux store
                const { updateADC, updateGeophone } = this.props;
                updateGeophone && updateGeophone(geophone);
                updateADC && updateADC(adc);
            } else {
                // Show error and return if failed
                const { t } = this.props;
                toast.error(t("views.history.toasts.metadata_error"));
                return;
            }
        }

        // Apply to component state
        this.setState({ adc, geophone, scale });
    }

    // Promise version of this.setState used for async/await
    promisedSetState = (newState: Partial<HistoryState>) => {
        return new Promise<void>((resolve) =>
            this.setState(newState as HistoryState, resolve)
        );
    };

    // Time picker change handler for start and end time
    handleTimeChange = (typ: string, value: number): void => {
        switch (typ) {
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

    // Fetch history waveform with specified format (JSON or SAC)
    handleQueryHistory = async (): Promise<unknown> => {
        const { history } = this.state;
        const { start, end } = history;

        // Check if start time is earlier than end time
        if (end - start <= 0 || !start || !end) {
            const { t } = this.props;
            const error = t("views.history.toasts.time_error");
            toast.error(error);
            return Promise.reject(error);
        }

        // Auto detect format and filename
        const { error, data } = await restfulApiByTag({
            body: history,
            tag: "history",
            timeout: QUERY_TIMEOUT,
            blob: history.format === "sac",
            filename: `${history.channel}-${history.start}-${history.end}.${history.format}`,
        });
        if (error) {
            const { t } = this.props;
            const error = t("views.history.toasts.export_sac_error");
            toast.error(error);
            return Promise.reject(error);
        }

        return data;
    };

    // Fetch events list from specified source and open modal dialog
    handleQueryEvents = async (): Promise<unknown> => {
        // Get events list from server
        const { trace } = this.state;
        const { error, data } = await restfulApiByTag({
            body: trace,
            tag: "trace",
            timeout: QUERY_TIMEOUT,
        });
        if (error) {
            return Promise.reject(error);
        }

        // Open modal dialog if success
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
                    const desc = `[M${magnitude.toFixed(
                        1
                    )}] ${event} / 时刻 ${getTimeString(
                        timestamp
                    )} / 深度 ${depth.toFixed(1)} km / 传播 ${estimated.toFixed(
                        1
                    )} s`;

                    return [region, timestamp + estimated * 1000, desc];
                }),
            },
        }));
    };

    // Choose event handler for event list modal dialog
    handleChooseEvent = async (value: string): Promise<void> => {
        // Get time range from event timestamp
        const span = TRACE_RANGE / 2;
        const time = new Date(value).getTime();
        const start = time - span;
        const end = time + span;

        // Update state and close modal dialog
        await this.promisedSetState({
            history: { start, end, format: "json" },
            modal: { ...this.state.modal, open: false },
        });
    };

    // Select dialog handler for data source & SAC file export option
    handleSelect = async (from: string, value: string): Promise<void> => {
        // Close select dialog
        const { t } = this.props;
        const select = {
            from: "history",
            dialog: {
                ...this.state.select.dialog,
                open: false,
            },
        };

        // Determine method based on `from` field
        switch (from) {
            case "history":
                await this.promisedSetState({
                    select,
                    history: {
                        ...this.state.history,
                        channel: value as any,
                        format: "sac",
                    },
                });
                await toast.promise(this.handleQueryHistory(), {
                    loading: t("views.history.toasts.is_exporting_sac"),
                    success: t("views.history.toasts.export_sac_success"),
                    error: t("views.history.toasts.export_sac_error"),
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
                    loading: t("views.history.toasts.is_fetching_events"),
                    success: t("views.history.toasts.fetch_events_success"),
                    error: t("views.history.toasts.fetch_events_error"),
                });
                break;
        }
    };

    // Fetch JSON format waveform and update labels
    handleQueryWaveform = async (): Promise<void> => {
        // Update format to JSON
        await this.promisedSetState({
            history: { ...this.state.history, format: "json" },
        });

        // Fetch waveform and update labels
        const { t } = this.props;
        const res = await toast.promise(this.handleQueryHistory(), {
            loading: t("views.history.toasts.is_fetching_waveform"),
            success: t("views.history.toasts.fetch_waveform_success"),
            error: t("views.history.toasts.fetch_waveform_error"),
        });
        if (res) {
            const { adc, geophone, scale } = this.state;
            const labels = setLabels(
                this.state.labels,
                res,
                adc,
                geophone,
                scale
            );
            const areas = setAreas(this.state.areas, res);
            this.setState({ areas, labels });
        }
    };

    // Open select dialog for SAC file export option
    handleQuerySACFile = async (): Promise<void> => {
        // Reset dialog content and open select dialog
        this.setState((state) => ({
            select: {
                ...state.select,
                from: "history",
                dialog: {
                    open: true,
                    values: [
                        ["Vertical", "EHZ"],
                        ["East-West", "EHE"],
                        ["North-South", "EHN"],
                    ],
                    title: { id: "views.history.selects.choose_channel.title" },
                },
            },
        }));
    };

    // Fetch event source list and open select dialog
    handleQuerySource = async (): Promise<unknown> => {
        // Set payload and fetch source list
        const { t } = this.props;
        const trace = {
            source: "show",
        };

        // Avoiding use toast.promise due to restfulApiByTag never reject
        const loader = toast.loading(
            t("views.history.toasts.is_fetching_source")
        );
        const { data, error } = await restfulApiByTag({
            body: trace,
            tag: "trace",
        });
        toast.remove(loader);

        // Show error and return if failed
        if (error || !data) {
            const error = t("views.history.toasts.fetch_source_error");
            toast.error(error);
            return Promise.reject(error);
        } else {
            toast.success(t("views.history.toasts.fetch_source_success"));
        }

        // Open data source select dialog if success
        this.setState((state) => ({
            select: {
                from: "trace",
                dialog: {
                    ...state.select.dialog,
                    open: true,
                    values: data.map((item: any) => [item.name, item.value]),
                    title: { id: "views.history.selects.choose_source.title" },
                },
            },
        }));
    };

    render() {
        const { areas, select, modal, history, labels } = this.state;
        const { from, dialog } = select;
        const { start, end } = history;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container className="mb-6" layout="grid">
                        <Card
                            className="h-[360px]"
                            label={{
                                id: "views.history.cards.query_history",
                            }}
                        >
                            <TimePicker
                                value={start}
                                label={{
                                    id: "views.history.time_pickers.start_time",
                                }}
                                onChange={(value) =>
                                    this.handleTimeChange("start", value)
                                }
                            />
                            <TimePicker
                                value={end}
                                label={{
                                    id: "views.history.time_pickers.end_time",
                                }}
                                onChange={(value) =>
                                    this.handleTimeChange("end", value)
                                }
                            />

                            <Button
                                className="mt-6 bg-indigo-700 hover:bg-indigo-800"
                                onClick={this.handleQueryWaveform}
                                label={{
                                    id: "views.history.buttons.query_waveform",
                                }}
                            />
                            <Button
                                className="bg-green-700 hover:bg-green-800"
                                onClick={this.handleQuerySACFile}
                                label={{
                                    id: "views.history.buttons.query_sac_file",
                                }}
                            />
                            <Button
                                className="bg-yellow-700 hover:bg-yellow-800"
                                onClick={this.handleQuerySource}
                                label={{
                                    id: "views.history.buttons.query_source",
                                }}
                            />
                        </Card>

                        <Card
                            label={{
                                id: "views.history.cards.analyse_history",
                            }}
                        >
                            <Container layout="grid">
                                {labels.map((label, index) => (
                                    <Label key={index} {...label} />
                                ))}
                            </Container>
                        </Card>
                    </Container>

                    {areas.map(({ area, chart }, index) => (
                        <Area key={index} {...area}>
                            <Chart
                                {...chart}
                                tooltip={true}
                                zooming={true}
                                animation={false}
                                tickPrecision={1}
                                tickInterval={10}
                            />
                        </Area>
                    ))}
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

export default connect(mapStateToProps, {
    updateGeophone,
    updateADC,
})(withTranslation()(History));
