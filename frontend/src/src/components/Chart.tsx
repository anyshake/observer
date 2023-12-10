import { Component } from "react";
import * as Highcharts from "highcharts";
import { HighchartsReact } from "highcharts-react-official";
import HighchartsBoost from "highcharts/modules/boost";
import { Options } from "highcharts";
import { WithTranslation, withTranslation } from "react-i18next";
HighchartsBoost(Highcharts);

export interface ChartProps {
    readonly height: number;
    readonly legend?: boolean;
    readonly zooming?: boolean;
    readonly tooltip?: boolean;
    readonly lineWidth?: number;
    readonly lineColor?: string;
    readonly animation?: boolean;
    readonly tickInterval?: number;
    readonly tickPrecision?: number;
    readonly backgroundColor?: string;
    // See https://github.com/highcharts/highcharts/issues/12242
    readonly series: Record<string, any>;
    readonly [key: string]: any;
}

export interface ChartState extends Options {}

class Chart extends Component<ChartProps & WithTranslation, ChartState> {
    constructor(props: ChartProps & WithTranslation) {
        super(props);
        const {
            height,
            legend,
            tooltip,
            zooming,
            animation,
            lineWidth,
            tickInterval,
            tickPrecision,
            lineColor,
            backgroundColor,
        } = this.props;
        this.state = {
            accessibility: { enabled: false },
            boost: {
                enabled: true,
                seriesThreshold: 5,
            },
            chart: {
                zooming: zooming ? { type: "x" } : {},
                marginTop: 20,
                height: height,
                animation: animation,
                backgroundColor: backgroundColor,
            },
            legend: {
                enabled: legend,
                itemStyle: { color: "#fff" },
            },
            plotOptions: {
                series: {
                    lineWidth: lineWidth,
                    states: { hover: { enabled: false } },
                },
            },
            xAxis: {
                labels: {
                    style: { color: "#fff" },
                    format: "{value:%H:%M:%S}",
                },
                type: "datetime",
                tickColor: "#fff",
                lineColor: lineColor,
            },
            yAxis: {
                labels: {
                    style: { color: "#fff" },
                    format: tickPrecision
                        ? `{value:${tickPrecision}f}`
                        : `{value:0.2f}`,
                },
                title: { text: "" },
                opposite: true,
                lineColor: lineColor,
                tickInterval: tickInterval,
            },
            tooltip: {
                enabled: tooltip,
                followPointer: true,
                followTouchMove: true,
                xDateFormat: "%Y-%m-%d %H:%M:%S",
                padding: 12,
            },
            credits: { enabled: false },
            time: { useUTC: false },
            title: { text: "" },
        };
    }

    componentDidUpdate(): void {
        const { t } = this.props;
        Highcharts.setOptions({
            lang: {
                resetZoom: t("components.chart.reset_zoom"),
                resetZoomTitle: t("components.chart.reset_zoom_title"),
            },
        });
    }

    render() {
        const { series } = this.props;
        const { state: options } = this;

        // sort series data again
        if (series.data) {
            series.data.sort((a: any, b: any) => {
                return a[0] - b[0];
            });
        } else if (series.length) {
            for (let i of series as any[]) {
                i.data.sort((a: any, b: any) => {
                    return a[0] - b[0];
                });
            }
        }

        return (
            <HighchartsReact
                options={{ ...options, series }}
                highcharts={Highcharts}
            />
        );
    }
}

export default withTranslation()(Chart);
