import { mdiServer } from '@mdi/js';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { LineChart } from '../../components/chart/LineChart';
import { Card } from '../../components/widget/Card';
import { HomeConstraints } from '../../config/constraints';
import type { GetHomeDataQuery } from '../../graphql';

type Series = Array<[number, number | null]>;

const appendPoint = (prev: Series, point: [number, number]): Series => {
    const last = [...prev].reverse().find(([, value]) => value !== null);
    const next: Series =
        last && (point[0] - last[0]) / 1000 > HomeConstraints.maxGapSeconds
            ? [[last[0] + 1, null], point]
            : [point];
    return [...prev, ...next].slice(-HomeConstraints.lineChartRetention);
};

const formatUptime = (ms: number) => ({
    days: Math.floor(ms / 86400000),
    hours: Math.floor((ms % 86400000) / 3600000),
    minutes: Math.floor((ms % 3600000) / 60000)
});

export const SystemUsage = ({ data }: { data: GetHomeDataQuery }) => {
    const { t } = useTranslation();
    const [metric, setMetric] = useState<'cpu' | 'memory'>('cpu');
    const [history, setHistory] = useState<{ cpu: Series; memory: Series }>({
        cpu: [],
        memory: []
    });
    const lastTime = useRef(0);

    useEffect(() => {
        if (data.getCurrentTime <= lastTime.current) {
            return;
        }
        lastTime.current = data.getCurrentTime;
        const { cpu, memory } = data.getSystemStatus;
        setHistory((prev) => ({
            cpu: appendPoint(prev.cpu, [data.getCurrentTime, cpu]),
            memory: appendPoint(prev.memory, [data.getCurrentTime, memory])
        }));
    }, [data]);

    const system = data.getSystemStatus;
    const meters = useMemo(
        () => [
            { key: 'cpu' as const, label: t('views.Home.cpu_usage.title'), value: system.cpu },
            {
                key: 'memory' as const,
                label: t('views.Home.memory_usage.title'),
                value: system.memory
            },
            { key: 'disk' as const, label: t('views.Home.disk_usage.title'), value: system.disk }
        ],
        [t, system]
    );

    return (
        <Card
            title={t('views.Home.system_resources.title')}
            iconPath={mdiServer}
            className="app-home-system-body"
        >
            <div className="app-home-meters">
                {meters.map(({ key, label, value }) => {
                    const selectable = key !== 'disk';
                    const percent = Math.min(100, Math.max(0, value));
                    return (
                        <button
                            key={key}
                            type="button"
                            disabled={!selectable}
                            aria-pressed={selectable ? metric === key : undefined}
                            className={`app-home-meter ${metric === key ? 'app-home-meter-active' : ''}`}
                            onClick={() => selectable && setMetric(key)}
                        >
                            <span className="app-home-meter-row">
                                <span className="app-home-meter-label">{label}</span>
                                <span className="app-home-meter-value">{value.toFixed(1)}%</span>
                            </span>
                            <span className="app-home-meter-track">
                                <span
                                    className={`app-home-meter-bar ${percent >= 90 ? 'app-home-meter-bar-high' : ''}`}
                                    style={{ width: `${percent}%` }}
                                />
                            </span>
                        </button>
                    );
                })}
                <p className="app-home-footnote">
                    {t('views.Home.system_resources.uptime', formatUptime(system.uptime))}
                </p>
            </div>

            <div className="app-home-system-chart">
                {history[metric].length > 1 ? (
                    <div className="app-home-chart-fill">
                        <LineChart
                            height="100%"
                            lineWidth={2}
                            yMin={0}
                            yMax={100}
                            yInterval={25}
                            data={history[metric]}
                        />
                    </div>
                ) : (
                    <div className="app-home-chart-empty">
                        <span className="loading loading-dots loading-sm" />
                        <span>{t('views.Home.system_resources.collecting')}</span>
                    </div>
                )}
            </div>
        </Card>
    );
};
