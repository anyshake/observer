import { HighchartsReactRefObject } from "highcharts-react-official";
import { RefObject, useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSelector } from "react-redux";

import { Banner, BannerProps } from "../../components/Banner";
import { Button } from "../../components/Button";
import { Chart, ChartProps } from "../../components/Chart";
import { Container } from "../../components/Container";
import { CollapseMode, Holder, HolderProps } from "../../components/Holder";
import { Input } from "../../components/Input";
import { Panel } from "../../components/Panel";
import { apiConfig } from "../../config/api";
import { ReduxStoreProps } from "../../config/store";
import { useSocket } from "../../helpers/hook/useSocket";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { CircularQueue2D } from "../../helpers/utils/CircularQueue2D";
import { userThrottle } from "../../helpers/utils/userThrottle";
import { getSocketUpdates, SocketUpdates } from "./getSocketUpdates";
import { handleSetBanner } from "./handleSetBanner";
import { handleSetCharts } from "./handleSetCharts";

const Realtime = () => {
	// States for the banner and charts
	const { stream } = useSelector(({ stream }: ReduxStoreProps) => stream);
	const { retention } = useSelector(({ retention }: ReduxStoreProps) => retention);
	const [banner, setBanner] = useState<BannerProps & { values?: Record<string, string> }>({
		type: "warning",
		title: "views.realtime.banner.warning.label",
		content: "views.realtime.banner.warning.text"
	});
	const [charts, setCharts] = useState<
		Record<
			string,
			{
				chart: ChartProps & {
					// rows: retention, columns: sampleRate
					// The buffer stores data in the form of [timestamp, x, y, z]
					buffer: CircularQueue2D<Float64Array>;
					ref: RefObject<HighchartsReactRefObject>;
					filter: {
						enabled: boolean;
						lowCorner?: number;
						highCorner?: number;
					};
				};
				holder: HolderProps & { values: Record<string, string> };
			}
		>
	>({
		z_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_HIDE,
				label: "views.realtime.charts.z_axis.label",
				text: "views.realtime.charts.z_axis.text",
				values: { max: "0", min: "0" }
			},
			chart: {
				buffer: new CircularQueue2D(retention, 0, Float64Array),
				backgroundColor: "#d97706",
				filter: { enabled: false },
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}Z`, type: "line", color: "#f1f5f9" }
			}
		},
		e_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_HIDE,
				label: "views.realtime.charts.e_axis.label",
				text: "views.realtime.charts.e_axis.text",
				values: { max: "0", min: "0" }
			},
			chart: {
				buffer: new CircularQueue2D(retention, 0, Float64Array),
				filter: { enabled: false },
				backgroundColor: "#10b981",
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}E`, type: "line", color: "#f1f5f9" }
			}
		},
		n_axis: {
			holder: {
				collapse: CollapseMode.COLLAPSE_HIDE,
				label: "views.realtime.charts.n_axis.label",
				text: "views.realtime.charts.n_axis.text",
				values: { max: "0", min: "0" }
			},
			chart: {
				buffer: new CircularQueue2D(retention, 0, Float64Array),
				backgroundColor: "#0ea5e9",
				filter: { enabled: false },
				ref: useRef<HighchartsReactRefObject>(null),
				series: { name: `${stream.channel}N`, type: "line", color: "#f1f5f9" }
			}
		}
	});

	// Handlers for WebSocket, passing the data to the charts
	const { t } = useTranslation();
	const { sensor } = useSelector(({ sensor }: ReduxStoreProps) => sensor);
	const handleSocketData = ({ data }: MessageEvent<SocketUpdates>) => {
		getSocketUpdates(
			data,
			(data) =>
				handleSetBanner(
					data,
					sensor.resolution,
					sensor.velocity ? "velocity" : "acceleration",
					setBanner
				),
			(data) => handleSetCharts(data, setCharts)
		);
	};
	const handleSocketError = () => {
		setBanner({
			type: "error",
			title: "views.realtime.banner.error.label",
			content: "views.realtime.banner.error.text"
		});
	};
	const handleSocketOpen = () => {
		sendUserAlert(t("views.realtime.toasts.websocket_connected"));
	};
	useSocket({
		backend: apiConfig.backend,
		endpoint: apiConfig.endpoints.socket,
		onData: handleSocketData,
		onError: handleSocketError,
		onClose: handleSocketError,
		onOpen: handleSocketOpen
	});

	// Set the height of the charts automatically according to the window size
	const [chartHeight, setChartHeight] = useState<number>(150);
	const setChartHeightToState = useCallback(() => {
		let height = Math.round((window.innerHeight * 0.6) / Object.keys(charts).length);
		if (height < 150) {
			height = 150;
		} else if (height > 500) {
			height = 500;
		}
		setChartHeight(height);
	}, [charts]);
	const handleWindowResize = userThrottle(() => {
		setChartHeightToState();
	}, 2000);
	useEffect(() => {
		setChartHeightToState();
		window.addEventListener("resize", handleWindowResize);
		return () => {
			window.removeEventListener("resize", handleWindowResize);
		};
	}, [setChartHeightToState, handleWindowResize]);

	// Handler for setting the corner frequency of the Butterworth filter
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

	// Handler for enabling/disabling the Butterworth filter
	const handleSwitchFilter = (chartKey: string) => {
		setCharts((charts) => ({
			...charts,
			[chartKey]: {
				...charts[chartKey],
				chart: {
					...charts[chartKey].chart,
					filter: {
						...charts[chartKey].chart.filter,
						enabled: !charts[chartKey].chart.filter.enabled
					}
				}
			}
		}));
	};

	return (
		<>
			<Banner
				type={banner.type}
				title={t(banner.title, { ...banner.values })}
				content={t(banner.content, { ...banner.values })}
			/>

			<Container className="pt-1">
				{Object.keys(charts).map((key) => (
					<Holder
						{...charts[key].holder}
						key={charts[key].holder.label}
						label={t(charts[key].holder.label, { channel: stream.channel })}
						text={t(charts[key].holder.text ?? "", {
							...charts[key].holder.values
						})}
						advanced={
							<Container className="max-w-96">
								<Panel
									label={t(
										`views.realtime.charts.${key}.advanced.panels.butterworth_filter.title`
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
												`views.realtime.charts.${key}.advanced.panels.butterworth_filter.low_corner_freq`
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
												`views.realtime.charts.${key}.advanced.panels.butterworth_filter.high_corner_freq`
											)}
										/>
									</Container>
									<Button
										label={t(
											`views.realtime.charts.${key}.advanced.panels.butterworth_filter.${
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
							boost={true}
							lineWidth={1}
							tooltip={true}
							zooming={true}
							animation={false}
							tickPrecision={1}
							tickInterval={100}
							height={chartHeight}
						/>
					</Holder>
				))}
			</Container>
		</>
	);
};

export default Realtime;
