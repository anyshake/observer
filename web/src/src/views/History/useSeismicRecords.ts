import { type Dispatch, type SetStateAction, useCallback, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { HistoryConstraints } from '../../config/constraints';
import { useGetSeismicRecordsLazyQuery } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sameStrings } from '../../helpers/utils/sameStrings';

export const useSeismicRecords = ({
    startTime,
    endTime,
    setIsBusy
}: {
    startTime: number;
    endTime: number;
    setIsBusy: Dispatch<SetStateAction<boolean>>;
}) => {
    const { t } = useTranslation();
    const [sampleRate, setSampleRate] = useState(0);
    const [chartData, setChartData] = useState<{ [key: string]: Array<[number, number | null]> }>(
        {}
    );
    const spectrogramSeries = useMemo(() => {
        const series: { [key: string]: Array<[number, number]> } = {};
        for (const channel in chartData) {
            const source = chartData[channel];
            let count = 0;
            for (let i = 0; i < source.length; i++) {
                if (source[i][1] !== null) {
                    count++;
                }
            }
            const points = new Array<[number, number]>(count);
            let offset = 0;
            for (let i = 0; i < source.length; i++) {
                const point = source[i];
                if (point[1] !== null) {
                    points[offset++] = [point[0], point[1]];
                }
            }
            series[channel] = points;
        }
        return series;
    }, [chartData]);
    const [getSeismicRecords] = useGetSeismicRecordsLazyQuery();
    const handleSearchRecords = useCallback(async () => {
        if (startTime > endTime) {
            sendUserAlert(t('views.History.search_records.invalid_time'), true);
            return;
        }
        const requestFn = async () => {
            setSampleRate(0);
            setChartData({});
            setIsBusy(true);
            const { data, error } = await getSeismicRecords({ variables: { startTime, endTime } });
            setIsBusy(false);
            if (error) {
                throw error;
            }
            if (!data?.getSeisRecordsByTime.length) {
                throw new Error(t('views.History.search_records.no_data'));
            }
            return data;
        };
        await sendPromiseAlert(
            requestFn(),
            t('views.History.search_records.searching'),
            ({ getSeisRecordsByTime }) => {
                const records = [...getSeisRecordsByTime].sort(
                    (a, b) => a!.timestamp - b!.timestamp
                );
                const currentChannels = new Map<string, number>();
                records.forEach((record) => {
                    record?.channelData.forEach((channel) => {
                        // channelId is 1-based; keep the physical index when earlier channels are absent.
                        currentChannels.set(channel.channelCode, channel.channelId - 1);
                    });
                });
                const channels = Array.from(currentChannels.keys());
                const channelIndices = Array.from(currentChannels.values());

                if (
                    !sameStrings(channels, prevChannelsRef.current) ||
                    channelIndices.some(
                        (index, position) => index !== prevChannelIndicesRef.current[position]
                    )
                ) {
                    setActiveChannels((prevChannels) => {
                        const newChannels = { ...prevChannels };
                        channels.forEach((channel, position) => {
                            newChannels[channel] = {
                                id: `${HistoryConstraints.id}_${channel}`,
                                index: channelIndices[position]
                            };
                        });
                        Object.keys(newChannels).forEach((channel) => {
                            if (!currentChannels.has(channel)) {
                                delete newChannels[channel];
                            }
                        });
                        prevChannelsRef.current = channels;
                        prevChannelIndicesRef.current = channelIndices;
                        return newChannels;
                    });
                }

                const sampleCounts: Record<string, number> = {};
                let gapCount = 0;
                let previousTimestamp: number | null = null;
                for (let recordIndex = 0; recordIndex < records.length; recordIndex++) {
                    const record = records[recordIndex]!;
                    if (previousTimestamp !== null && record.timestamp - previousTimestamp > 2000) {
                        gapCount++;
                    }
                    previousTimestamp = record.timestamp;
                    for (
                        let channelIndex = 0;
                        channelIndex < record.channelData.length;
                        channelIndex++
                    ) {
                        const channel = record.channelData[channelIndex];
                        sampleCounts[channel.channelCode] =
                            (sampleCounts[channel.channelCode] ?? 0) + channel.data.length;
                    }
                }

                setChartData((prevChartData) => {
                    const newChartData: typeof prevChartData = {};
                    const writeAt: Record<string, number> = {};

                    for (let channelIndex = 0; channelIndex < channels.length; channelIndex++) {
                        const channel = channels[channelIndex];
                        const previous = prevChartData[channel];
                        const base = previous?.length ?? 0;
                        const total = base + (sampleCounts[channel] ?? 0) + gapCount;
                        const dataArray = previous ?? new Array<[number, number | null]>(total);
                        if (previous) {
                            dataArray.length = total;
                        }
                        newChartData[channel] = dataArray;
                        writeAt[channel] = base;
                    }

                    let lastTimestamp: number | null = null;
                    for (let recordIndex = 0; recordIndex < records.length; recordIndex++) {
                        const { timestamp, sampleRate, channelData } = records[recordIndex]!;
                        const interval = 1000 / sampleRate;

                        if (lastTimestamp !== null && timestamp - lastTimestamp > 2000) {
                            for (
                                let channelIndex = 0;
                                channelIndex < channels.length;
                                channelIndex++
                            ) {
                                const channel = channels[channelIndex];
                                const offset = writeAt[channel];
                                const previousPoint = newChartData[channel][offset - 1];
                                const gapStart = previousPoint?.[0] ?? lastTimestamp;
                                newChartData[channel][offset] = [
                                    gapStart + (timestamp - gapStart) / 2,
                                    null
                                ];
                                writeAt[channel] = offset + 1;
                            }
                        }

                        for (
                            let channelIndex = 0;
                            channelIndex < channelData.length;
                            channelIndex++
                        ) {
                            const channel = channelData[channelIndex];
                            const dataArray = newChartData[channel.channelCode];
                            if (!dataArray) {
                                continue;
                            }
                            let offset = writeAt[channel.channelCode];
                            const samples = channel.data;
                            for (let sampleIndex = 0; sampleIndex < samples.length; sampleIndex++) {
                                dataArray[offset++] = [
                                    timestamp + sampleIndex * interval,
                                    samples[sampleIndex]
                                ];
                            }
                            writeAt[channel.channelCode] = offset;
                        }

                        lastTimestamp = timestamp;
                    }
                    setSampleRate(records.length > 0 ? records[0]!.sampleRate : 0);

                    return newChartData;
                });

                return t('views.History.search_records.success', { count: records.length });
            },
            (error) => t('views.History.search_records.error', { error })
        );
    }, [t, startTime, endTime, getSeismicRecords, setIsBusy]);

    const [activeChannels, setActiveChannels] = useState<
        Record<string, { id: string; index: number }>
    >({});
    const prevChannelsRef = useRef<string[]>([]);
    const prevChannelIndicesRef = useRef<number[]>([]);

    return { sampleRate, chartData, spectrogramSeries, activeChannels, handleSearchRecords };
};
