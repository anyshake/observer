import { createRef, RefObject, useCallback, useRef, useState } from 'react';

import { getDisplayUnit, isGroundMotionChannel } from '../../components/chart/channelUnits';
import { DequeChartHandle } from '../../components/chart/DequeChart';
import { DequeSpectrogramHandle } from '../../components/chart/DequeSpectrogram';
import type { IGroundMotion, IGroundMotionRecord } from '../../components/chart/groundMotion';
import { RealTimeConstraints, RetentionConstraints } from '../../config/constraints';
import { getSocketApiUrl } from '../../helpers/app/getSocketApiUrl';
import { useSocket } from '../../helpers/request/useSocket';
import { sameStrings } from '../../helpers/utils/sameStrings';
import { useCredentialStore } from '../../stores/credential';

export const useRealtimeStream = () => {
    const { getCredential } = useCredentialStore();
    const [sampleRate, setSampleRate] = useState(0);
    const [updatedAt, setUpdatedAt] = useState(0);
    const [channelUnits, setChannelUnits] = useState<Record<string, string>>({});
    const [latestValues, setLatestValues] = useState<Record<string, number>>({});
    const [activeChannels, setActiveChannels] = useState<
        Record<string, { id: string; index: number }>
    >({});

    // Latest per-second peaks of each channel; null when the station has no metadata.
    const [groundMotion, setGroundMotion] = useState<Record<string, IGroundMotion>>({});
    // Per-second peaks kept for the longest retention, looked up while hovering the waveform.
    const motionHistoryRef = useRef<Record<string, IGroundMotionRecord[]>>({});
    const getGroundMotionRecords = useCallback(
        (channel: string) => motionHistoryRef.current[channel],
        []
    );

    const waveformRefs = useRef<{ [key: string]: RefObject<DequeChartHandle> }>({});
    const spectrogramRefs = useRef<{ [key: string]: RefObject<DequeSpectrogramHandle> }>({});
    const prevChannelsRef = useRef<string[]>([]);
    const prevChannelIndicesRef = useRef<number[]>([]);
    const updateChannels = useCallback((channelData: Record<string, { channel_id: number }>) => {
        const currentChannels = Object.keys(channelData);
        // Backend channel_id starts at 1; object key order does not reflect physical channel order.
        const channelIndices = currentChannels.map(
            (channel) => channelData[channel].channel_id - 1
        );

        if (
            sameStrings(currentChannels, prevChannelsRef.current) &&
            channelIndices.every(
                (index, position) => index === prevChannelIndicesRef.current[position]
            )
        ) {
            return;
        }
        setActiveChannels((prevChannels) => {
            const newChannels = { ...prevChannels };
            currentChannels.forEach((channel, position) => {
                if (!newChannels[channel]) {
                    waveformRefs.current[channel] =
                        createRef<DequeChartHandle>() as RefObject<DequeChartHandle>;
                    spectrogramRefs.current[channel] =
                        createRef<DequeSpectrogramHandle>() as RefObject<DequeSpectrogramHandle>;
                }
                newChannels[channel] = {
                    id: `${RealTimeConstraints.id}_${channel}`,
                    index: channelIndices[position]
                };
            });
            Object.keys(newChannels).forEach((channel) => {
                if (!currentChannels.includes(channel)) {
                    delete newChannels[channel];
                    delete waveformRefs.current[channel];
                    delete spectrogramRefs.current[channel];
                }
            });
            prevChannelsRef.current = currentChannels;
            prevChannelIndicesRef.current = channelIndices;
            return newChannels;
        });
    }, []);

    const updatedAtRef = useRef(0);
    const { readyState, sendMessage } = useSocket(
        {
            url: getSocketApiUrl(),
            heartbeat: { message: 'ping', response: 'pong' },
            onOpen: () => {
                sendMessage(getCredential().token);
                sendMessage('client hello');
            },
            onData: ({ data }) => {
                const { channel_data, sample_rate, record_time, current_time } = data;

                const units: Record<string, string> = {};
                const latest: Record<string, number> = {};
                for (const channel in channel_data) {
                    const { data_array, input_unit, sensitivity } = channel_data[channel];
                    const unit = getDisplayUnit(input_unit, sensitivity);
                    units[channel] = unit;
                    const samples = unit
                        ? data_array.map((value: number) => value / sensitivity)
                        : data_array;
                    latest[channel] = samples[samples.length - 1];
                    waveformRefs.current[channel]?.current?.addData(
                        samples,
                        record_time,
                        current_time,
                        sample_rate
                    );
                    spectrogramRefs.current[channel]?.current?.addData(
                        data_array,
                        record_time,
                        current_time,
                        sample_rate
                    );
                }
                setChannelUnits((prev) =>
                    sameStrings(Object.keys(prev), Object.keys(units)) &&
                    Object.keys(units).every((channel) => prev[channel] === units[channel])
                        ? prev
                        : units
                );
                setLatestValues(latest);

                const motion: Record<string, IGroundMotion> = {};
                const history: Record<string, IGroundMotionRecord[]> = {};
                const expired = record_time - RetentionConstraints.maxValue * 1000;
                for (const channel in channel_data) {
                    const { pga, pgv, pgd, data_array, input_unit } = channel_data[channel];
                    if (!isGroundMotionChannel(channel, input_unit, { pga, pgv, pgd })) {
                        continue;
                    }
                    const peaks = { pga: pga ?? null, pgv: pgv ?? null, pgd: pgd ?? null };
                    motion[channel] = peaks;

                    // Records stay sorted by start; a stream restart clears older entries.
                    const records = (motionHistoryRef.current[channel] ?? []).filter(
                        ({ start }) => start >= expired && start < record_time
                    );
                    const duration =
                        sample_rate > 0 ? (data_array.length * 1000) / sample_rate : 1000;
                    records.push({ start: record_time, end: record_time + duration, ...peaks });
                    history[channel] = records;
                }
                motionHistoryRef.current = history;
                setGroundMotion(motion);

                updatedAtRef.current = record_time;
                setSampleRate((prev) => (prev === sample_rate ? prev : sample_rate));
                setUpdatedAt((prev) =>
                    Math.floor(prev / 1000) === Math.floor(record_time / 1000) ? prev : record_time
                );

                updateChannels(channel_data);
            }
        },
        true
    );
    return {
        sampleRate,
        updatedAt,
        channelUnits,
        latestValues,
        activeChannels,
        groundMotion,
        getGroundMotionRecords,
        waveformRefs,
        spectrogramRefs,
        updatedAtRef,
        readyState
    };
};
