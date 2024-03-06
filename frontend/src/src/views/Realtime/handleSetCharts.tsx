import { Dispatch, RefObject, SetStateAction } from "react";
import { SocketUpdates } from "./getSocketUpdates";
import { ChartProps } from "../../components/Chart";
import { HolderProps } from "../../components/Holder";
import { HighchartsReactRefObject } from "highcharts-react-official";
import store from "../../config/store";
import { getVoltageArr } from "../../helpers/seismic/getVoltageArr";
import { getVelocityArr } from "../../helpers/seismic/getVelocityArr";
import { getAccelerationArr } from "../../helpers/seismic/getAccelerationArr";
import { fallbackScale, globalConfig } from "../../config/global";
import {
    FilterPassband,
    getFilteredCounts,
} from "../../helpers/seismic/getFilteredCounts";

export const handleSetCharts = (
    res: SocketUpdates,
    stateFn: Dispatch<
        SetStateAction<
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
                    holder: HolderProps & { values: Record<string, string> };
                }
            >
        >
    >
) => {
    if (!res.ts) {
        return;
    }
    stateFn((prev) => {
        const { ts, ...channels } = res;
        const { retention } = store.getState().retention;
        Object.keys(prev).forEach((key) => {
            if (!(key in res)) {
                return;
            }

            // Append new data to buffer and remove expired data
            const channelData = channels[key as keyof typeof channels];
            const { buffer } = prev[key].chart;
            buffer.push({ ts, data: channelData });
            const timeoutThreshold = ts - retention * 1000;
            while (buffer[0].ts < timeoutThreshold) {
                buffer.shift();
            }

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
                    ? `Band pass [${lowFreqCorner.toFixed(
                          1
                      )}-${highFreqCorner.toFixed(1)} Hz]`
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

            // Get seismic values and apply to holder fields
            const { adc } = store.getState().adc;
            const voltageArr = getVoltageArr(
                channelData,
                adc.resolution,
                adc.fullscale
            );
            const { geophone } = store.getState().geophone;
            const sensitivity = geophone.sensitivity / 100;
            const velocityArr = getVelocityArr(voltageArr, sensitivity);
            const channelDataSpanMS = 1000 / channelData.length;
            const accelerationArr = getAccelerationArr(
                velocityArr,
                channelDataSpanMS
            );
            const { scale } = store.getState().scale;
            const intensityStandard =
                globalConfig.scales.find(
                    (scaleObj) => scaleObj.property().value === scale
                ) ?? fallbackScale;
            const { holder } = prev[key];
            holder.values = {
                pga: accelerationArr
                    .reduce((a, b) => Math.max(Math.abs(a), Math.abs(b)), 0)
                    .toFixed(5),
                pgv: velocityArr
                    .reduce((a, b) => Math.max(Math.abs(a), Math.abs(b)), 0)
                    .toFixed(5),
                intensity: `${scale} ${intensityStandard.getIntensity(
                    {
                        rawData: channelData,
                        currentPGV: Math.max(...velocityArr),
                        currentPGA: Math.max(...accelerationArr),
                    },
                    { adc, geophone }
                )}`,
            };
        });

        return prev;
    });
};
