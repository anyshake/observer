import { Options, SeriesOptionsType } from "highcharts";
import HighCharts from "highcharts/highcharts";
import HighchartsBoost from "highcharts/modules/boost";
import { HighchartsReact, HighchartsReactRefObject } from "highcharts-react-official";
import { ForwardedRef, forwardRef, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

const hasWebGLSupport = () => {
	if (window.WebGLRenderingContext) {
		const canvas = document.createElement("canvas");
		const names = ["webgl", "experimental-webgl", "webgl2", "moz-webgl", "webkit-3d"];
		return names.some((name) => {
			try {
				return !!canvas.getContext(name);
			} catch (e) {
				return false;
			}
		});
	}

	return false;
};

export interface ChartProps {
	readonly boost?: boolean;
	readonly title?: string;
	readonly height?: number;
	readonly legend?: boolean;
	readonly tooltip?: boolean;
	readonly zooming?: boolean;
	readonly lineWidth?: number;
	readonly lineColor?: string;
	readonly animation?: boolean;
	readonly tickInterval?: number;
	readonly tickPrecision?: number;
	readonly backgroundColor?: string;
	readonly series: SeriesOptionsType;
}

export const Chart = forwardRef(
	(props: ChartProps, ref: ForwardedRef<HighchartsReactRefObject>) => {
		const pixelRatio = window.devicePixelRatio * 0.65;
		const {
			boost,
			title,
			series,
			height,
			legend,
			tooltip,
			zooming,
			animation,
			lineWidth,
			tickInterval,
			tickPrecision,
			lineColor,
			backgroundColor
		} = props;

		const [supportWebGL, setSupportWebGL] = useState(false);

		useEffect(() => {
			setSupportWebGL(hasWebGLSupport());
		}, []);

		const [chartOptions, setChartOptions] = useState<Options>({
			chart: {
				zooming: zooming ? { type: "x" } : {},
				marginTop: 20,
				height,
				animation,
				backgroundColor
			},
			xAxis: {
				labels: {
					style: { color: "#fff" },
					format: "{value:%H:%M:%S}"
				},
				type: "datetime",
				tickColor: "#fff",
				lineColor
			},
			yAxis: {
				labels: {
					style: { color: "#fff" },
					format: tickPrecision ? `{value:${tickPrecision}f}` : `{value:0.2f}`
				},
				title: { text: "" },
				opposite: true,
				lineColor,
				tickInterval
			},
			tooltip: {
				enabled: tooltip,
				followPointer: true,
				followTouchMove: true,
				xDateFormat: "%Y-%m-%d %H:%M:%S",
				padding: 12
			},
			legend: { enabled: legend, itemStyle: { color: "#fff" } },
			plotOptions: {
				series: {
					lineWidth,
					turboThreshold: boost ? 10 : 0,
					boostThreshold: boost ? 1 : 0,
					states: { hover: { enabled: false } }
				}
			},
			title: {
				text: title,
				style: {
					color: "#fff",
					fontSize: "10px",
					fontWeight: "normal"
				}
			},
			boost: { enabled: supportWebGL, pixelRatio },
			accessibility: { enabled: false },
			credits: { enabled: false },
			time: { useUTC: false },
			series: [series]
		});

		const { t } = useTranslation();

		useEffect(() => {
			HighchartsBoost(HighCharts);
		}, []);

		useEffect(() => {
			HighCharts.setOptions({
				lang: {
					resetZoom: t("components.chart.reset_zoom"),
					resetZoomTitle: t("components.chart.reset_zoom_title")
				}
			});
		}, [t]);

		useEffect(() => {
			setChartOptions((prev) => ({
				...prev,
				chart: { ...prev.chart, height },
				title: { ...prev.title, text: title },
				boost: { ...prev.boost, enabled: supportWebGL }
			}));
		}, [height, title, supportWebGL]);

		return <HighchartsReact ref={ref} options={chartOptions} highcharts={HighCharts} />;
	}
);
