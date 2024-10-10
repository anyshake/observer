import { HighchartsReactRefObject } from "highcharts-react-official";
import { RefObject, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSelector } from "react-redux";
import { useSearchParams } from "react-router-dom";

import { Button } from "../../components/Button";
import { Chart, ChartProps } from "../../components/Chart";
import { Container } from "../../components/Container";
import { Form, FormProps } from "../../components/Form";
import { CollapseMode, Holder, HolderProps } from "../../components/Holder";
import { Input } from "../../components/Input";
import { Label, LabelProps } from "../../components/Label";
import { Panel } from "../../components/Panel";
import { Select, SelectProps } from "../../components/Select";
import { TimePicker } from "../../components/TimePicker";
import { apiConfig, traceCommonResponseModel1 } from "../../config/api";
import { i18nConfig } from "../../config/i18n";
import { RouterComponentProps } from "../../config/router";
import { ReduxStoreProps } from "../../config/store";
import { getFlagByCountry } from "../../helpers/i18n/getFlagByCountry";
import { sendPromiseAlert } from "../../helpers/interact/sendPromiseAlert";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { FilterPassband, getFilteredCounts } from "../../helpers/seismic/getFilteredCounts";
import { getNormalizedData } from "../../helpers/seismic/getNormalizedData";
import { asyncSleep } from "../../helpers/utils/asyncSleep";
import { CircularQueue2D } from "../../helpers/utils/CircularQueue2D";
import { getTimeString } from "../../helpers/utils/getTimeString";
import { setClipboardText } from "../../helpers/utils/setClipboardText";
import { getExportFileName } from "./getExportFileName";
import { handleSetCharts } from "./handleSetCharts";

const History = ({ locale }: RouterComponentProps) => {
	// Freezing the state of the component until the metadata is initialized
	const { stream } = useSelector(({ stream }: ReduxStoreProps) => stream);
	const [isCurrentBusy, setIsCurrentBusy] = useState(!stream.initialized);
	useEffect(() => {
		setIsCurrentBusy(!stream.initialized);
	}, [stream.initialized]);

	// Attempt to retrieve start and end time from URL search params
	const { duration } = useSelector(({ duration }: ReduxStoreProps) => duration);
	const [searchParams, setSearchParams] = useSearchParams();
	const currentTimestamp = Date.now();
	const [queryDuration, setQueryDuration] = useState<{
		start_time: number;
		end_time: number;
	}>({
		start_time: searchParams.has("start")
			? Number(searchParams.get("start"))
			: currentTimestamp - 1000 * duration,
		end_time: searchParams.has("end") ? Number(searchParams.get("end")) : currentTimestamp
	});

	// Handler for timepicker changes (start and end time)
	const handleTimeChange = (value: number, is_end_time: boolean) =>
		setQueryDuration((prev) => {
			if (is_end_time) {
				return { ...prev, end_time: value };
			}
			return { ...prev, start_time: value };
		});

	// Basic states for form, other attributes are set on demand
	const [form, setForm] = useState<FormProps & { values?: Record<string, string | number> }>({
		open: false,
		inputType: "select"
	});
	const handleCloseForm = () => {
		setForm({ ...form, open: false });
	};

	// Basic states for select, other attributes are set on demand
	const [select, setSelect] = useState<SelectProps>({ open: false });
	const handleCloseSelect = () => {
		setSelect({ ...select, open: false });
	};

	// States for component data (labels, charts)
	const [labels, setLabels] = useState<
		Record<string, LabelProps & { values?: Record<string, string> }>
	>({
		z_axis: {
			label: "views.history.labels.z_axis_detail.label",
			value: "views.history.labels.z_axis_detail.value",
			values: { max: "0", min: "0" }
		},
		e_axis: {
			label: "views.history.labels.e_axis_detail.label",
			value: "views.history.labels.z_axis_detail.value",
			values: { max: "0", min: "0" }
		},
		n_axis: {
			label: "views.history.labels.n_axis_detail.label",
			value: "views.history.labels.z_axis_detail.value",
			values: { max: "0", min: "0" }
		}
	});
	const [charts, setCharts] = useState<
		Record<
			string,
			{
				chart: ChartProps & {
					// rows: retention, columns: sampleRate
					// The buffer stores data in the form of [timestamp, x, y, z]
					buffer: CircularQueue2D;
					ref: RefObject<HighchartsReactRefObject>;
					filter: { enabled: boolean; lowCorner?: number; highCorner?: number };
				};
				holder: HolderProps;
			}
		>
	>({
		z_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_HIDE,
				label: "views.history.charts.z_axis.label",
				text: "views.history.charts.z_axis.text"
			},
			chart: {
				buffer: new CircularQueue2D(0, 0),
				backgroundColor: "#d97706",
				filter: { enabled: false },
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}Z`, type: "line", color: "#f1f5f9" }
			}
		},
		e_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_SHOW,
				label: "views.history.charts.e_axis.label",
				text: "views.history.charts.e_axis.text"
			},
			chart: {
				buffer: new CircularQueue2D(0, 0),
				backgroundColor: "#10b981",
				filter: { enabled: false },
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}E`, type: "line", color: "#f1f5f9" }
			}
		},
		n_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_SHOW,
				label: "views.history.charts.n_axis.label",
				text: "views.history.charts.n_axis.text"
			},
			chart: {
				buffer: new CircularQueue2D(0, 0),
				backgroundColor: "#0ea5e9",
				filter: { enabled: false },
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}N`, type: "line", color: "#f1f5f9" }
			}
		}
	});

	// Handlers for setting corner frequencies for Butterworth filter
	const handleSetCornerFreq = (chartKey: string, lowCorner: boolean, value: number) =>
		setCharts((charts) => ({
			...charts,
			[chartKey]: {
				...charts[chartKey],
				chart: {
					...charts[chartKey].chart,
					filter: {
						...charts[chartKey].chart.filter,
						[lowCorner ? "lowCorner" : "highCorner"]: value
					}
				}
			}
		}));

	// Handler for applying Butterworth filter to chart data
	const handleSwitchFilter = (chartKey: string) => {
		setCharts((prev) => {
			const filterEnabled = !prev[chartKey].chart.filter.enabled;
			const { lowCorner, highCorner } = prev[chartKey].chart.filter;
			const { lowFreqCorner, highFreqCorner } = {
				lowFreqCorner: lowCorner ?? 0.1,
				highFreqCorner: highCorner ?? 10
			};

			// Get filtered values and apply to chart data
			const { current: chartObj } = prev[chartKey].chart.ref;
			if (chartObj) {
				const chartData = prev[chartKey].chart.buffer
					.readAll()
					.map((item) => {
						const timestamp = item[0];
						const channelData = item.slice(1);
						const normalizedData = Float32Array.from(
							getNormalizedData(Array.from(channelData), 0)
						);
						if (filterEnabled) {
							const filteredData = getFilteredCounts(normalizedData, {
								poles: 4,
								lowFreqCorner,
								highFreqCorner,
								sampleRate: normalizedData.length,
								passbandType: FilterPassband.BAND_PASS
							});
							return Array.from(filteredData).map((value, index) => [
								timestamp + (1000 / filteredData.length) * index,
								value
							]);
						}
						return Array.from(normalizedData).map((value, index) => [
							timestamp + (1000 / normalizedData.length) * index,
							value
						]);
					})
					.flat();
				chartObj.chart.series[0].setData(chartData, true, false, false);
			}

			// Update chart state
			const currentChart = {
				...prev[chartKey],
				chart: {
					...prev[chartKey].chart,
					filter: { ...prev[chartKey].chart.filter, enabled: filterEnabled },
					title: filterEnabled ? `Band pass [${lowFreqCorner}-${highFreqCorner} Hz]` : ""
				}
			};
			return { ...prev, [chartKey]: currentChart };
		});
	};

	// Handler for querying waveform data, server returns waveform data in JSON format
	const { t } = useTranslation();
	const handleQueryWaveform = async () => {
		const { start_time, end_time } = queryDuration;
		if (!start_time || !end_time || start_time >= end_time) {
			sendUserAlert(t("views.history.toasts.duration_error"), true);
			return;
		}

		const { backend } = apiConfig;
		const payload = { start_time, end_time, channel: "", format: "json" };

		const res = await sendPromiseAlert(
			requestRestApi<
				typeof apiConfig.endpoints.history.model.request,
				typeof apiConfig.endpoints.history.model.response.common,
				typeof apiConfig.endpoints.history.model.response.error
			>({
				backend,
				payload,
				timeout: 180,
				throwError: true,
				endpoint: apiConfig.endpoints.history
			}),
			t("views.history.toasts.is_fetching_waveform"),
			t("views.history.toasts.fetch_waveform_success"),
			t("views.history.toasts.fetch_waveform_error")
		);

		handleSetCharts(res, setCharts, setLabels);
	};

	// Handler for exporting waveform data as file, server returns waveform data in SAC / MiniSEED binary format
	const handleExportAsFile = () => {
		const { start_time, end_time } = queryDuration;
		if (!start_time || !end_time || start_time >= end_time) {
			sendUserAlert(t("views.history.toasts.duration_error"), true);
			return;
		}

		const handleSubmitForm = async (format: string, channelCode: string) => {
			setForm((prev) => ({ ...prev, open: false }));

			const { backend } = apiConfig;
			const payload = { start_time, end_time, channel: channelCode, format };

			const baseFileName = getExportFileName(
				start_time,
				`${stream.channel}${channelCode}`,
				stream
			);
			const miniSeedFileName = `${baseFileName}.mseed`;
			const sacFileName = `${baseFileName}.sac`;

			await sendPromiseAlert(
				requestRestApi<
					typeof apiConfig.endpoints.history.model.request,
					typeof apiConfig.endpoints.history.model.response.common,
					typeof apiConfig.endpoints.history.model.response.error
				>({
					backend,
					payload,
					timeout: 180,
					throwError: true,
					endpoint: apiConfig.endpoints.history,
					blobOptions: {
						fileName: format === "miniseed" ? miniSeedFileName : sacFileName
					}
				}),
				t("views.history.toasts.is_exporting_file"),
				t("views.history.toasts.export_file_success"),
				t("views.history.toasts.export_file_error")
			);
		};

		const handleChooseFormat = async (channel: string) => {
			setForm((prev) => ({ ...prev, open: false }));
			await asyncSleep(300);

			setForm((prev) => ({
				...prev,
				open: true,
				selectOptions: [
					{ label: "MiniSEED", value: "miniseed" },
					{ label: "SAC", value: "sac" }
				],
				onSubmit: (format) => handleSubmitForm(format, channel),
				title: "views.history.forms.choose_format.title",
				cancelText: "views.history.forms.choose_format.cancel",
				submitText: "views.history.forms.choose_format.submit",
				placeholder: "views.history.forms.choose_format.placeholder"
			}));
		};

		setForm((prev) => ({
			...prev,
			open: true,
			selectOptions: [
				{ label: "Z Axis", value: "Z" },
				{ label: "E Axis", value: "E" },
				{ label: "N Axis", value: "N" }
			],
			onSubmit: handleChooseFormat,
			title: "views.history.forms.choose_channel.title",
			cancelText: "views.history.forms.choose_channel.cancel",
			submitText: "views.history.forms.choose_channel.submit",
			placeholder: "views.history.forms.choose_channel.placeholder"
		}));
	};

	// Handler for querying seismic event data and selecting its corresponding start and end time
	const { fallback: fallbackLocale } = i18nConfig;
	const handleQueryEvent = async () => {
		const { backend } = apiConfig;
		const payload = { source: "list" };

		const res = await sendPromiseAlert(
			requestRestApi<
				typeof apiConfig.endpoints.trace.model.request,
				typeof apiConfig.endpoints.trace.model.response.common,
				typeof apiConfig.endpoints.trace.model.response.error
			>({
				backend,
				payload,
				throwError: true,
				endpoint: apiConfig.endpoints.trace
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
					throwError: true,
					payload: { source },
					endpoint: apiConfig.endpoints.trace
				}),
				t("views.history.toasts.is_fetching_events"),
				t("views.history.toasts.fetch_events_success"),
				t("views.history.toasts.fetch_events_error")
			)) as unknown as typeof traceCommonResponseModel1;
			if (!res?.data || res?.data?.length === 0) {
				sendUserAlert(t("views.history.toasts.no_events_found"), true);
				return;
			}
			const eventList = res.data.map(
				({ distance, magnitude, region, timestamp, depth, estimation }, index) => [
					`[${magnitude.map((m) => `${m.type} ${m.value.toFixed(1)}`).join(", ")}] ${region}`,
					JSON.stringify({
						id: index,
						start_time: Math.round(
							timestamp + estimation.p_wave * 1000 - duration * 500
						),
						end_time: Math.round(timestamp + estimation.s_wave * 1000 + duration * 500)
					}),
					t("views.history.selects.choose_event.template", {
						time: getTimeString(timestamp),
						distance: distance.toFixed(1),
						p_wave: estimation.p_wave.toFixed(1),
						s_wave: estimation.s_wave.toFixed(1),
						depth: depth !== -1 ? depth.toFixed(1) : "?"
					})
				]
			);

			const handleSelectEvent = (value: string) => {
				setSelect((prev) => ({ ...prev, open: false }));
				const { start_time, end_time } = JSON.parse(value);
				setQueryDuration({ start_time, end_time });
				sendUserAlert(t("views.history.toasts.event_select_success"));
			};
			setSelect((prev) => ({
				...prev,
				open: true,
				options: eventList,
				onClose: handleCloseSelect,
				onSelect: handleSelectEvent,
				title: "views.history.selects.choose_event.title"
			}));
		};

		setForm((prev) => {
			const currentLocale = locale ?? fallbackLocale;
			return {
				...prev,
				open: true,
				selectOptions: res.data
					.sort((a, b) => {
						if ("country" in a && "country" in b) {
							if (a.country === b.country) {
								const aLocale =
									a.locales[currentLocale as keyof typeof a.locales] ??
									a.locales[a.default as keyof typeof a.locales];
								const bLocale =
									b.locales[currentLocale as keyof typeof b.locales] ??
									b.locales[b.default as keyof typeof b.locales];
								return aLocale.localeCompare(bLocale);
							}
							return a.country.localeCompare(b.country);
						}
						return 0;
					})
					.map((seisSource) => {
						if ("locales" in seisSource) {
							const source =
								seisSource.locales[
									currentLocale as keyof typeof seisSource.locales
								] ??
								seisSource.locales[
									seisSource.default as keyof typeof seisSource.locales
								];
							return {
								label: `${getFlagByCountry(seisSource.country)} ${source}`,
								value: seisSource.id
							};
						}
						return { label: "", value: "" };
					}),
				onSubmit: handleSubmitForm,
				title: "views.history.forms.choose_source.title",
				cancelText: "views.history.forms.choose_source.cancel",
				submitText: "views.history.forms.choose_source.submit",
				placeholder: "views.history.forms.choose_source.placeholder"
			};
		});
	};

	// Handler for getting share link with current query duration
	const handleGetShareLink = async () => {
		const { start_time, end_time } = queryDuration;
		if (!start_time || !end_time || start_time >= end_time) {
			sendUserAlert(t("views.history.toasts.duration_error"), true);
			return;
		}

		const newSearchParams = new URLSearchParams();
		newSearchParams.set("start", String(start_time));
		newSearchParams.set("end", String(end_time));
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

	return (
		<>
			<Container
				className={`my-6 gap-4 grid lg:grid-cols-2 ${
					isCurrentBusy ? "cursor-progress" : ""
				}`}
			>
				<Panel label={t("views.history.panels.query_history")}>
					<TimePicker
						value={queryDuration.start_time}
						currentLocale={locale ?? fallbackLocale}
						label={t("views.history.time_pickers.start_time")}
						onChange={(value) => handleTimeChange(value, false)}
					/>
					<TimePicker
						value={queryDuration.end_time}
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
						className={`bg-green-700 hover:bg-green-800 ${
							isCurrentBusy ? "cursor-wait" : ""
						}`}
						onClick={() => {
							if (!isCurrentBusy) {
								setIsCurrentBusy(true);
								handleExportAsFile();
								setIsCurrentBusy(false);
							}
						}}
						label={t("views.history.buttons.exoprt_as_file")}
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

				<Panel className="" label={t("views.history.panels.analyze_history")}>
					{Object.values(labels).map(({ label, value, values, ...rest }) => (
						<Label
							{...rest}
							key={label}
							value={t(value, values)}
							label={t(label, { channel: stream.channel })}
						/>
					))}
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
				<Select {...select} onClose={handleCloseSelect} title={t(select.title ?? "")} />
			</Container>

			{Object.keys(charts).map((key) => (
				<Holder
					key={charts[key].holder.label}
					text={t(charts[key].holder.text ?? "")}
					label={t(charts[key].holder.label ?? "", { channel: stream.channel })}
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
											handleSetCornerFreq(key, true, Number(value))
										}
										defaultValue={0.1}
										type="number"
										disabled={charts[key].chart.filter.enabled}
										numberLimit={{ max: 100, min: 0.1 }}
										label={t(
											`views.history.charts.${key}.advanced.panels.butterworth_filter.low_corner_freq`
										)}
									/>
									<Input
										onValueChange={(value) =>
											handleSetCornerFreq(key, false, Number(value))
										}
										defaultValue={10}
										type="number"
										disabled={charts[key].chart.filter.enabled}
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
