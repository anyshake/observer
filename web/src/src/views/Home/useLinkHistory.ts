import { useEffect, useMemo, useRef, useState } from 'react';

import { HomeConstraints } from '../../config/constraints';
import type { GetHomeDataQuery } from '../../graphql';

interface ISample {
    readonly time: number;
    readonly rate: number;
    readonly errors: number;
    readonly gap: boolean;
}

// Derives packet throughput and new reception errors from successive counter polls.
export const useLinkHistory = (data: GetHomeDataQuery | undefined) => {
    const previous = useRef<{ time: number; frames: number; errors: number } | null>(null);
    const [samples, setSamples] = useState<ISample[]>([]);

    useEffect(() => {
        if (!data) {
            return;
        }
        const time = data.getCurrentTime;
        const { frames, errors } = data.getDeviceStatus;
        const last = previous.current;
        previous.current = { time, frames, errors };
        if (!last || time <= last.time) {
            return;
        }
        if (frames < last.frames || errors < last.errors) {
            // Counters were reset, e.g. after the hardware link restarted.
            setSamples([]);
            return;
        }

        const seconds = (time - last.time) / 1000;
        const sample = {
            time,
            rate: (frames - last.frames) / seconds,
            errors: errors - last.errors,
            gap: seconds > HomeConstraints.maxGapSeconds
        };
        setSamples((prev) => [...prev, sample].slice(-HomeConstraints.linkHistoryRetention));
    }, [data]);

    const throughput = useMemo(
        () =>
            samples.flatMap(
                ({ time, rate, gap }, index): Array<[number, number | null]> =>
                    gap && index > 0
                        ? [
                              [samples[index - 1].time + 1, null],
                              [time, rate]
                          ]
                        : [[time, rate]]
            ),
        [samples]
    );
    const recentErrors = useMemo(
        () => samples.reduce((total, { errors }) => total + errors, 0),
        [samples]
    );
    const currentRate = samples.length > 0 ? samples[samples.length - 1].rate : null;

    return { throughput, recentErrors, currentRate };
};
