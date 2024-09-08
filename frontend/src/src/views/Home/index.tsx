import { mdiBug, mdiClockTimeFour, mdiWebCheck } from "@mdi/js";
import { HighchartsReactRefObject } from "highcharts-react-official";
import { RefObject, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { Banner, BannerProps } from "../../components/Banner";
import { Chart, ChartProps } from "../../components/Chart";
import { Container } from "../../components/Container";
import { Holder, HolderProps } from "../../components/Holder";
import { Label, LabelProps } from "../../components/Label";
import { MapBox, MapBoxProps } from "../../components/MapBox";
import { useInterval } from "../../helpers/hook/useInterval";
import { getStationUpdates } from "./getStationUpdates";
import { handleSetBanner } from "./handleSetBanner";
import { handleSetCharts } from "./handleSetCharts";
import { handleSetLabels } from "./handleSetLabels";
import { handleSetMap } from "./handleSetMap";

const Home = () => {
	const [banner, setBanner] = useState<BannerProps & { values?: Record<string, string> }>({
		type: "warning",
		title: "views.home.banner.warning.title",
		content: "views.home.banner.warning.content"
	});
	const [labels, setLabels] = useState<Record<string, LabelProps>>({
		messages: {
			color: true,
			value: "0",
			icon: mdiWebCheck,
			unit: "views.home.labels.messages.unit",
			label: "views.home.labels.messages.label"
		},
		errors: {
			color: true,
			value: "0",
			icon: mdiBug,
			unit: "views.home.labels.errors.unit",
			label: "views.home.labels.errors.label"
		},
		elapsed: {
			color: true,
			value: "0",
			icon: mdiClockTimeFour,
			unit: "views.home.labels.elapsed.unit",
			label: "views.home.labels.elapsed.label"
		}
	});
	const [charts, setCharts] = useState<
		Record<
			string,
			{
				chart: ChartProps & {
					ref: RefObject<HighchartsReactRefObject>;
				};
				holder: HolderProps & { values: Record<string, string> };
			}
		>
	>({
		cpu: {
			chart: {
				height: 250,
				lineWidth: 5,
				backgroundColor: "#22c55e",
				ref: useRef<HighchartsReactRefObject>(null),
				series: { type: "line", color: "#fff" }
			},
			holder: {
				label: "views.home.charts.cpu.label",
				text: "views.home.charts.cpu.text",
				values: { usage: "0.00" }
			}
		},
		memory: {
			chart: {
				height: 250,
				lineWidth: 5,
				backgroundColor: "#06b6d4",
				ref: useRef<HighchartsReactRefObject>(null),
				series: { type: "line", color: "#fff" }
			},
			holder: {
				label: "views.home.charts.memory.label",
				text: "views.home.charts.memory.text",
				values: { usage: "0.00" }
			}
		}
	});
	const [map, setMap] = useState<{
		mapbox: MapBoxProps;
		holder: HolderProps & { values: Record<string, string> };
	}>({
		mapbox: {
			zoom: 6,
			minZoom: 3,
			maxZoom: 7,
			flyTo: true,
			center: [0, 0],
			dragging: false,
			tile: "/tiles/{z}/{x}/{y}.webp"
		},
		holder: {
			label: "views.home.map.label",
			text: "views.home.map.text",
			values: { longitude: "0.00", latitude: "0.00", elevation: "0.00" }
		}
	});

	useInterval(
		() =>
			getStationUpdates(
				(res) => {
					handleSetBanner(res, setBanner);
				},
				(res) => {
					handleSetLabels(res, setLabels);
				},
				(res) => {
					handleSetCharts(res, setCharts);
				},
				(res) => {
					handleSetMap(res, setMap);
				}
			),
		2000,
		true
	);

	const { t } = useTranslation();

	return (
		<>
			<Banner
				type={banner.type}
				title={t(banner.title, { ...banner.values })}
				content={t(banner.content, { ...banner.values })}
			/>

			<Container className="mt-5 grid lg:grid-cols-3">
				{Object.values(labels).map(({ label, unit, ...item }) => (
					<Label {...item} key={label} label={t(label)} unit={t(unit ?? "")} />
				))}
			</Container>

			<Container className="mt-5 gap-4 grid grid-cols-1 md:grid-cols-2">
				{Object.values(charts).map(({ holder, chart }) => (
					<Holder
						key={holder.label}
						label={t(holder.label)}
						text={t(holder.text ?? "", { ...holder.values })}
					>
						<Chart {...chart} />
					</Holder>
				))}
			</Container>

			<Container>
				<Holder
					label={t(map.holder.label)}
					text={t(map.holder.text ?? "", { ...map.holder.values })}
				>
					<MapBox className="h-[400px]" {...map.mapbox} />
				</Holder>
			</Container>
		</>
	);
};

export default Home;
