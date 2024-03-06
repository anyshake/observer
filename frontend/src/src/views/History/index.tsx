import { useTranslation } from "react-i18next";
import { Container } from "../../components/Container";
import { Panel } from "../../components/Panel";
import { TimePicker } from "../../components/TimePicker";
import { Button } from "../../components/Button";
import { CollapseMode, Holder, HolderProps } from "../../components/Holder";
import { Chart, ChartProps } from "../../components/Chart";
import { Label, LabelProps } from "../../components/Label";
import { RefObject, useEffect, useRef, useState } from "react";
import { i18nConfig } from "../../config/i18n";
import { RouterComponentProps } from "../../config/router";
import { HighchartsReactRefObject } from "highcharts-react-official";
import { Select, SelectProps } from "../../components/Select";
import { ReduxStoreProps } from "../../config/store";
import { useSelector } from "react-redux";
import { Input } from "../../components/Input";
import { useSearchParams } from "react-router-dom";
import { setClipboardText } from "../../helpers/utils/setClipboardText";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { apiConfig, traceCommonResponseModel1 } from "../../config/api";
import { sendPromiseAlert } from "../../helpers/interact/sendPromiseAlert";
import { handleSetCharts } from "./handleSetCharts";
import { handleSetLabels } from "./handleSetLabels";
import { Form, FormProps } from "../../components/Form";
import { getSACFileName } from "./getSACFileName";
import { getTimeString } from "../../helpers/utils/getTimeString";
import {
    FilterPassband,
    getFilteredCounts,
} from "../../helpers/seismic/getFilteredCounts";

const History = (props: RouterComponentProps) => {
    const { t } = useTranslation();

    const { station } = useSelector(({ station }: ReduxStoreProps) => station);
    const { duration } = useSelector(
        ({ duration }: ReduxStoreProps) => duration
    );

    const [isCurrentBusy, setIsCurrentBusy] = useState(!station.initialized);

    useEffect(() => {
        setIsCurrentBusy(!station.initialized);
    }, [station.initialized]);

    const currentTimestamp = Date.now();
    const [searchParams, setSearchParams] = useSearchParams();

    const [queryDuration, setQueryDuration] = useState<{
        start: number;
        end: number;
    }>({
        start: searchParams.has("start")
            ? Number(searchParams.get("start"))
            : currentTimestamp - 1000 * duration,
        end: searchParams.has("end")
            ? Number(searchParams.get("end"))
            : currentTimestamp,
    });

    const handleTimeChange = (value: number, end: boolean) =>
        setQueryDuration((prev) => {
            if (end) {
                return { ...prev, end: value };
            }
            return { ...prev, start: value };
        });

    const [form, setForm] = useState<
        FormProps & { values?: Record<string, string | number> }
    >({ open: false, inputType: "select" });

    const handleCloseForm = () => setForm({ ...form, open: false });

    const [select, setSelect] = useState<SelectProps>({ open: false });

    const handleCloseSelect = () => setSelect({ ...select, open: false });

    const [labels, setLabels] = useState<
        Record<string, LabelProps & { values?: Record<string, string> }>
    >({
        ehz: { label: "views.history.labels.ehz_detail.label", value: "-" },
        ehe: { label: "views.history.labels.ehe_detail.label", value: "-" },
        ehn: { label: "views.history.labels.ehn_detail.label", value: "-" },
    });

    const [charts, setCharts] = useState<
        Record<
            string,
            {
                chart: ChartProps & {
                    buffer: { ts: number; data: number[] }[];
                    ref: RefObject<HighchartsReactRefObject>;
                    filter: {
                        enabled: boolean;
                        lowCorner?: number;
                        highCorner?: number;
                    };
                };
                holder: HolderProps;
            }
        >
    >({
        ehz: {
            holder: {
                collapse: CollapseMode.COLLAPSE_HIDE,
                label: "views.history.charts.ehz.label",
                text: "views.history.charts.ehz.text",
            },
            chart: {
                buffer: [],
                backgroundColor: "#d97706",
                filter: { enabled: false },
                ref: useRef<HighchartsReactRefObject>(null),
                series: { name: "EHZ", type: "line", color: "#f1f5f9" },
            },
        },
        ehe: {
            holder: {
                collapse: CollapseMode.COLLAPSE_SHOW,
                label: "views.history.charts.ehe.label",
                text: "views.history.charts.ehe.text",
            },
            chart: {
                buffer: [],
                backgroundColor: "#10b981",
                filter: { enabled: false },
                ref: useRef<HighchartsReactRefObject>(null),
                series: { name: "EHE", type: "line", color: "#f1f5f9" },
            },
        },
        ehn: {
            holder: {
                collapse: CollapseMode.COLLAPSE_SHOW,
                label: "views.history.charts.ehn.label",
                text: "views.history.charts.ehn.text",
            },
            chart: {
                buffer: [],
                backgroundColor: "#0ea5e9",
                filter: { enabled: false },
                ref: useRef<HighchartsReactRefObject>(null),
                series: { name: "EHN", type: "line", color: "#f1f5f9" },
            },
        },
    });

    const handleSetCornerFreq = (
        chartKey: string,
        lowCorner: boolean,
        value: number
    ) =>
        setCharts((charts) => ({
            ...charts,
            [chartKey]: {
                ...charts[chartKey],
                chart: {
                    ...charts[chartKey].chart,
                    filter: {
                        ...charts[chartKey].chart.filter,
                        [lowCorner ? "lowCorner" : "highCorner"]: value,
                    },
                },
            },
        }));

    const handleSwitchFilter = (chartKey: string) => {
        setCharts((prev) => {
            const filterEnabled = !prev[chartKey].chart.filter.enabled;
            const { lowCorner, highCorner } = prev[chartKey].chart.filter;
            const { lowFreqCorner, highFreqCorner } = {
                lowFreqCorner: lowCorner ?? 0.1,
                highFreqCorner: highCorner ?? 10,
            };

            // Get filtered values and apply to chart data
            const chartData = prev[chartKey].chart.buffer
                .map(({ ts, data }) => {
                    const filteredData = filterEnabled
                        ? getFilteredCounts(data, {
                              poles: 4,
                              lowFreqCorner,
                              highFreqCorner,
                              sampleRate: data.length,
                              passbandType: FilterPassband.BAND_PASS,
                          })
                        : data;
                    const dataSpanMS = 1000 / filteredData.length;
                    return filteredData.map((value, index) => [
                        ts + dataSpanMS * index,
                        value,
                    ]);
                })
                .reduce((acc, curArr) => acc.concat(curArr), []);
            const { current: chartObj } = prev[chartKey].chart.ref;
            if (!!chartObj) {
                const { series } = chartObj.chart;
                series[0].setData(chartData, true, false, false);
            }

            const currentChart = {
                ...prev[chartKey],
                chart: {
                    ...prev[chartKey].chart,
                    filter: {
                        ...prev[chartKey].chart.filter,
                        enabled: filterEnabled,
                    },
                    title: filterEnabled
                        ? `Band pass [${lowFreqCorner}-${highFreqCorner} Hz]`
                        : "",
                },
            };
            return { ...prev, [chartKey]: currentChart };
        });
    };

    const handleQueryWaveform = async () => {
        const { start, end } = queryDuration;
        if (!start || !end || start >= end) {
            sendUserAlert(t("views.history.toasts.duration_error"), true);
            return;
        }

        const { backend } = apiConfig;
        const payload = { start, end, channel: "", format: "json" };

        const res = await sendPromiseAlert(
            requestRestApi<
                typeof apiConfig.endpoints.history.model.request,
                typeof apiConfig.endpoints.history.model.response.common,
                typeof apiConfig.endpoints.history.model.response.error
            >({
                backend,
                payload,
                timeout: 120,
                throwError: true,
                endpoint: apiConfig.endpoints.history,
            }),
            t("views.history.toasts.is_fetching_waveform"),
            t("views.history.toasts.fetch_waveform_success"),
            t("views.history.toasts.fetch_waveform_error")
        );

        handleSetLabels(res, setLabels);
        handleSetCharts(res, setCharts);
    };

    const handleExportSACFile = () => {
        const { start, end } = queryDuration;
        if (!start || !end || start >= end) {
            sendUserAlert(t("views.history.toasts.duration_error"), true);
            return;
        }

        const handleSubmitForm = async (channel: string) => {
            setForm((prev) => ({ ...prev, open: false }));

            const { backend } = apiConfig;
            const payload = { start, end, channel, format: "sac" };
            const sacFileName = getSACFileName(start, channel, station);

            await sendPromiseAlert(
                requestRestApi<
                    typeof apiConfig.endpoints.history.model.request,
                    typeof apiConfig.endpoints.history.model.response.common,
                    typeof apiConfig.endpoints.history.model.response.error
                >({
                    backend,
                    payload,
                    timeout: 120,
                    throwError: true,
                    endpoint: apiConfig.endpoints.history,
                    blobOptions: { filename: sacFileName },
                }),
                t("views.history.toasts.is_exporting_sac"),
                t("views.history.toasts.export_sac_success"),
                t("views.history.toasts.export_sac_error")
            );
        };

        setForm((prev) => ({
            ...prev,
            open: true,
            selectOptions: [
                { label: "EHZ", value: "EHZ" },
                { label: "EHE", value: "EHE" },
                { label: "EHN", value: "EHN" },
            ],
            onSubmit: handleSubmitForm,
            title: "views.history.forms.choose_channel.title",
            cancelText: "views.history.forms.choose_channel.cancel",
            submitText: "views.history.forms.choose_channel.submit",
            placeholder: "views.history.forms.choose_channel.placeholder",
        }));
    };

    const handleQueryEvent = async () => {
        const { backend } = apiConfig;
        const payload = { source: "show" };

        const res = await sendPromiseAlert(
            requestRestApi<
                typeof apiConfig.endpoints.trace.model.request,
                typeof apiConfig.endpoints.trace.model.response.common,
                typeof apiConfig.endpoints.trace.model.response.error
            >({
                backend,
                payload,
                timeout: 30,
                throwError: true,
                endpoint: apiConfig.endpoints.trace,
            }),
            t("views.history.toasts.is_fetching_source"),
            t("views.history.toasts.fetch_source_success"),
            t("views.history.toasts.fetch_source_error")
        );
        if (!res?.data) {
            return;
        }

        const handleSubmitForm = async (source: string) => {
            setForm((prev) => ({ ...prev, open: false }));

            const res = (await sendPromiseAlert(
                requestRestApi<
                    typeof apiConfig.endpoints.trace.model.request,
                    typeof apiConfig.endpoints.trace.model.response.common,
                    typeof apiConfig.endpoints.trace.model.response.error
                >({
                    backend,
                    timeout: 60,
                    throwError: true,
                    payload: { source },
                    endpoint: apiConfig.endpoints.trace,
                }),
                t("views.history.toasts.is_fetching_events"),
                t("views.history.toasts.fetch_events_success"),
                t("views.history.toasts.fetch_events_error")
            )) as unknown as typeof traceCommonResponseModel1;
            if (!res?.data) {
                return;
            }

            const handleSelectEvent = (value: string) => {
                setSelect((prev) => ({ ...prev, open: false }));
                const [start, end] = value.split("|").map(Number);
                setQueryDuration({ start, end });
                sendUserAlert(t("views.history.toasts.event_select_success"));
            };

            const eventList = res.data.map(
                ({
                    distance,
                    magnitude,
                    region,
                    event,
                    timestamp,
                    depth,
                    estimation,
                }) => [
                    region,
                    `${timestamp + estimation.p * 1000}|${
                        timestamp + estimation.s * 1000
                    }`,
                    t("views.history.selects.choose_event.template", {
                        event,
                        time: getTimeString(timestamp),
                        magnitude: magnitude.toFixed(1),
                        distance: distance.toFixed(1),
                        p_wave: estimation.p.toFixed(1),
                        s_wave: estimation.s.toFixed(1),
                        depth: depth !== -1 ? depth.toFixed(1) : "-",
                    }),
                ]
            );
            setSelect((prev) => ({
                ...prev,
                open: true,
                options: eventList,
                onClose: handleCloseSelect,
                onSelect: handleSelectEvent,
                title: "views.history.selects.choose_event.title",
            }));
        };

        setForm((prev) => ({
            ...prev,
            open: true,
            selectOptions: res.data.map((source) => {
                if ("name" in source && "value" in source) {
                    return { label: source.name, value: source.value };
                }
                return { label: "", value: "" };
            }),
            onSubmit: handleSubmitForm,
            title: "views.history.forms.choose_source.title",
            cancelText: "views.history.forms.choose_source.cancel",
            submitText: "views.history.forms.choose_source.submit",
            placeholder: "views.history.forms.choose_source.placeholder",
        }));
    };

    const handleGetShareLink = async () => {
        const { start, end } = queryDuration;
        if (!start || !end || start >= end) {
            sendUserAlert(t("views.history.toasts.duration_error"), true);
            return;
        }

        const newSearchParams = new URLSearchParams();
        newSearchParams.set("start", String(start));
        newSearchParams.set("end", String(end));
        setSearchParams(newSearchParams);
        const newFullUrl = window.location.href;
        const success = await setClipboardText(newFullUrl);
        sendUserAlert(
            success
                ? t("views.history.toasts.copy_link_success")
                : t("views.history.toasts.copy_link_error"),
            !success
        );
    };

    const { locale } = props;
    const { fallback: fallbackLocale } = i18nConfig;

    return (
        <>
            <Container
                className={`my-6 gap-4 grid lg:grid-cols-2 ${
                    isCurrentBusy ? "cursor-progress" : ""
                }`}
            >
                <Panel label={t("views.history.panels.query_history")}>
                    <TimePicker
                        value={queryDuration.start}
                        currentLocale={locale ?? fallbackLocale}
                        label={t("views.history.time_pickers.start_time")}
                        onChange={(value) => handleTimeChange(value, false)}
                    />
                    <TimePicker
                        value={queryDuration.end}
                        currentLocale={locale ?? fallbackLocale}
                        label={t("views.history.time_pickers.end_time")}
                        onChange={(value) => handleTimeChange(value, true)}
                    />

                    <Button
                        className={`bg-indigo-700 hover:bg-indigo-800 ${
                            isCurrentBusy ? "cursor-wait" : ""
                        }`}
                        onClick={async () => {
                            if (!isCurrentBusy) {
                                setIsCurrentBusy(true);
                                await handleQueryWaveform();
                                setIsCurrentBusy(false);
                            }
                        }}
                        label={t("views.history.buttons.query_waveform")}
                    />
                    <Button
                        className="bg-green-700 hover:bg-green-800"
                        onClick={handleExportSACFile}
                        label={t("views.history.buttons.query_sac_file")}
                    />
                    <Button
                        className={`bg-yellow-700 hover:bg-yellow-800 ${
                            isCurrentBusy ? "cursor-wait" : ""
                        }`}
                        onClick={async () => {
                            if (!isCurrentBusy) {
                                setIsCurrentBusy(true);
                                await handleQueryEvent();
                                setIsCurrentBusy(false);
                            }
                        }}
                        label={t("views.history.buttons.query_source")}
                    />
                    <Button
                        className="bg-cyan-700 hover:bg-cyan-800"
                        onClick={handleGetShareLink}
                        label={t("views.history.buttons.get_share_link")}
                    />
                </Panel>

                <Panel
                    className=""
                    label={t("views.history.panels.analyze_history")}
                >
                    {Object.values(labels).map(
                        ({ label, value, values, ...rest }) => (
                            <Label
                                {...rest}
                                key={label}
                                label={t(label)}
                                value={t(value, values)}
                            />
                        )
                    )}
                </Panel>

                <Form
                    {...form}
                    onClose={handleCloseForm}
                    title={t(form.title ?? "")}
                    cancelText={t(form.cancelText ?? "")}
                    submitText={t(form.submitText ?? "")}
                    placeholder={t(form.placeholder ?? "")}
                    content={t(form.content ?? "", { ...form.values })}
                />
                <Select
                    {...select}
                    onClose={handleCloseSelect}
                    title={t(select.title ?? "")}
                />
            </Container>

            {Object.keys(charts).map((key) => (
                <Holder
                    key={charts[key].holder.label}
                    text={t(charts[key].holder.text ?? "")}
                    label={t(charts[key].holder.label ?? "")}
                    advanced={
                        <Container className="max-w-96">
                            <Panel
                                label={t(
                                    `views.history.charts.${key}.advanced.panels.butterworth_filter.title`
                                )}
                                embedded={true}
                            >
                                <Container className="flex flex-col md:flex-row gap-4">
                                    <Input
                                        onValueChange={(value) =>
                                            handleSetCornerFreq(
                                                key,
                                                true,
                                                Number(value)
                                            )
                                        }
                                        defaultValue={0.1}
                                        type="number"
                                        disabled={
                                            charts[key].chart.filter.enabled
                                        }
                                        numberLimit={{ max: 100, min: 0.1 }}
                                        label={t(
                                            `views.history.charts.${key}.advanced.panels.butterworth_filter.low_corner_freq`
                                        )}
                                    />
                                    <Input
                                        onValueChange={(value) =>
                                            handleSetCornerFreq(
                                                key,
                                                false,
                                                Number(value)
                                            )
                                        }
                                        defaultValue={10}
                                        type="number"
                                        disabled={
                                            charts[key].chart.filter.enabled
                                        }
                                        numberLimit={{ max: 100, min: 0.1 }}
                                        label={t(
                                            `views.history.charts.${key}.advanced.panels.butterworth_filter.high_corner_freq`
                                        )}
                                    />
                                </Container>
                                <Button
                                    label={t(
                                        `views.history.charts.${key}.advanced.panels.butterworth_filter.${
                                            charts[key].chart.filter.enabled
                                                ? "disable_filter"
                                                : "enable_filter"
                                        }`
                                    )}
                                    className="bg-indigo-600 hover:bg-indigo-700"
                                    onClick={() => handleSwitchFilter(key)}
                                />
                            </Panel>
                        </Container>
                    }
                >
                    <Chart
                        {...charts[key].chart}
                        height={300}
                        boost={true}
                        lineWidth={1}
                        tooltip={true}
                        zooming={true}
                        animation={true}
                        tickPrecision={1}
                        tickInterval={100}
                    />
                </Holder>
            ))}
        </>
    );
};

export default History;
