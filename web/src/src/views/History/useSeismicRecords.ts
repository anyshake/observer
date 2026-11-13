import { type Dispatch, type SetStateAction, useCallback, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getDisplayUnit, isGroundMotionChannel } from '../../components/chart/channelUnits';
import type { IGroundMotionRecord } from '../../components/chart/groundMotion';
import { HistoryConstraints } from '../../config/constraints';
import { useGetSeismicRecordsLazyQuery } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { useRequestSignal } from '../../helpers/request/useRequestSignal';
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
    const getRequestSignal = useRequestSignal();
    const [sampleRate, setSampleRate] = useState(0);
    const [channelUnits, setChannelUnits] = useState<Record<string, string>>({});
    const [chartData, setChartData] = useState<{ [key: string]: Array<[number, number | null]> }>(
        {}
    );
    // Per-second peaks sorted by record start; looked up while hovering the waveform.
    const [groundMotion, setGroundMotion] = useState<{ [key: string]: IGroundMotionRecord[] }>({});
    const [spectrogramSeries, setSpectrogramSeries] = useState<Record<string, [number, number][]>>(
        {}
    );
    const [getSeismicRecords] = useGetSeismicRecordsLazyQuery();
    const handleSearchRecords = useCallback(async () => {
        if (startTime > endTime) {
            sendUserAlert(t('views.History.search_records.invalid_time'), true);
            return;
        }
        const requestFn = async () => {
            setSampleRate(0);
            setChartData({});
            setGroundMotion({});
            setChannelUnits({});
            setSpectrogramSeries({});
            setIsBusy(true);
            const signal = getRequestSignal();
            const { data, error } = await getSeismicRecords({
                variables: { startTime, endTime },
                context: { signal }
            }).finally(() => {
                if (!signal.aborted) {
                    setIsBusy(false);
                }
            });
            signal.throwIfAborted();
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
                const units: Record<string, string> = {};
                records.forEach((record) => {
                    record?.channelData.forEach((channel) => {
                        // channelId is 1-based; keep the physical index when earlier channels are absent.
                        currentChannels.set(channel.channelCode, channel.channelId - 1);
                        const unit = getDisplayUnit(channel.inputUnit, channel.sensitivity);
                        const previous = units[channel.channelCode];
                        units[channel.channelCode] =
                            previous === undefined || previous === unit ? unit : '';
                    });
                });
                setChannelUnits(units);
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

                const motion: { [key: string]: IGroundMotionRecord[] } = {};
                for (const record of records) {
                    const duration = Math.max(
                        1,
                        ...record!.channelData.map(
                            (channel) => (channel.data.length * 1000) / record!.sampleRate
                        )
                    );
                    for (const { channelCode, pga, pgv, pgd, inputUnit } of record!.channelData) {
                        if (!isGroundMotionChannel(channelCode, inputUnit, { pga, pgv, pgd })) {
                            continue;
                        }
                        (motion[channelCode] ??= []).push({
                            start: record!.timestamp,
                            end: record!.timestamp + duration,
                            pga: pga ?? null,
                            pgv: pgv ?? null,
                            pgd: pgd ?? null
                        });
                    }
                }
                setGroundMotion(motion);

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
                    const spectra: Record<string, [number, number][]> = {};
                    const spectraWriteAt: Record<string, number> = {};

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
                        spectra[channel] = new Array(sampleCounts[channel] ?? 0);
                        spectraWriteAt[channel] = 0;
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
                            const spectrum = spectra[channel.channelCode];
                            let spectrumOffset = spectraWriteAt[channel.channelCode];
                            const divisor = units[channel.channelCode] ? channel.sensitivity : 1;
                            for (let sampleIndex = 0; sampleIndex < samples.length; sampleIndex++) {
                                dataArray[offset++] = [
                                    timestamp + sampleIndex * interval,
                                    samples[sampleIndex] / divisor
                                ];
                                spectrum[spectrumOffset++] = [
                                    timestamp + sampleIndex * interval,
                                    samples[sampleIndex]
                                ];
                            }
                            writeAt[channel.channelCode] = offset;
                            spectraWriteAt[channel.channelCode] = spectrumOffset;
                        }

                        lastTimestamp = timestamp;
                    }
                    setSpectrogramSeries(spectra);
                    setSampleRate(records.length > 0 ? records[0]!.sampleRate : 0);

                    return newChartData;
                });

                return t('views.History.search_records.success', { count: records.length });
            },
            (error) => t('views.History.search_records.error', { error })
        );
    }, [t, startTime, endTime, getSeismicRecords, getRequestSignal, setIsBusy]);

    const [activeChannels, setActiveChannels] = useState<
        Record<string, { id: string; index: number }>
    >({});
    const prevChannelsRef = useRef<string[]>([]);
    const prevChannelIndicesRef = useRef<number[]>([]);

    return {
        sampleRate,
        chartData,
        groundMotion,
        channelUnits,
        spectrogramSeries,
        activeChannels,
        handleSearchRecords
    };
};
