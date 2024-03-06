import { Dispatch, RefObject, SetStateAction } from "react";
import { apiConfig } from "../../config/api";
import { ChartProps } from "../../components/Chart";
import { HighchartsReactRefObject } from "highcharts-react-official";
import { HolderProps } from "../../components/Holder";
import {
    FilterPassband,
    getFilteredCounts,
} from "../../helpers/seismic/getFilteredCounts";

export const handleSetCharts = (
    res:
        | typeof apiConfig.endpoints.history.model.response.common
        | typeof apiConfig.endpoints.history.model.response.error,
    stateFn: Dispatch<
        SetStateAction<
            Record<
                string,
                {
                    chart: ChartProps & {
                        buffer: {
                            ts: number;
                            data: number[];
                        }[];
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
        >
    >
) => {
    if (!res?.data) {
        return;
    }
    stateFn((prev) => {
        Object.keys(prev).forEach((key) => {
            if (!res.data.every((obj) => key in obj)) {
                return;
            }

            // Set channel buffer from API response
            const buffer = res.data.map(({ ts, ...channels }) => ({
                data: channels[key as keyof typeof channels],
                ts,
            }));
            prev[key].chart.buffer = buffer;

            // Get filter settings and apply to chart
            const {
                enabled: filterEnabled,
                lowCorner,
                highCorner,
            } = prev[key].chart.filter;
            const { lowFreqCorner, highFreqCorner } = {
                lowFreqCorner: lowCorner ?? 0.1,
                highFreqCorner: highCorner ?? 10,
            };
            prev[key].chart = {
                ...prev[key].chart,
                title: filterEnabled
                    ? `Band pass [${lowFreqCorner}-${highFreqCorner} Hz]`
                    : "",
            };

            // Get filtered values and apply to chart data
            const chartData = buffer
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
            const { current: chartObj } = prev[key].chart.ref;
            if (!!chartObj) {
                const { series } = chartObj.chart;
                series[0].setData(chartData, true, false, false);
            }
        });

        return prev;
    });
};
