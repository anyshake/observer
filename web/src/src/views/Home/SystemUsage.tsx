import { mdiChip, mdiHarddisk, mdiMemory } from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { GaugeChart } from '../../components/chart/GaugeChart';
import { LineChart } from '../../components/chart/LineChart';
import { Skeleton } from '../../components/ui/Skeleton';
import { HomeConstraints } from '../../config/constraints';
import type { GetHomeDataQuery } from '../../graphql';

export const SystemUsage = ({ data: getHomeDataData }: { data: GetHomeDataQuery | undefined }) => {
    const { t } = useTranslation();
    const lineChartConfig = useMemo(
        () => ({
            height: 250,
            lineWidth: 5,
            yMax: 100,
            yMin: 0,
            yInterval: 20,
            yPosition: 'left' as const
        }),
        []
    );

    const [memoryData, setMemoryData] = useState<Array<[number, number | null]>>([]);
    const [cpuData, setCpuData] = useState<Array<[number, number | null]>>([]);
    const memoryDataRef = useRef(memoryData);
    const cpuDataRef = useRef(cpuData);
    const processChartData = useCallback(
        (
            prevData: Array<[number, number | null]>,
            newPoint: [number, number],
            retention: number,
            maxGap: number
        ): Array<[number, number | null]> => {
            if (prevData.length === 0) {
                return [newPoint];
            }

            const lastValidPoint = [...prevData].reverse().find((point) => point[1] !== null);
            if (!lastValidPoint) {
                return [...prevData.slice(-retention), newPoint];
            }

            const timeDiff = (newPoint[0] - lastValidPoint[0]) / 1000;
            if (timeDiff > maxGap) {
                return [...prevData.slice(-retention), [lastValidPoint[0] + 1, null], newPoint];
            }

            return [...prevData.slice(-retention), newPoint];
        },
        []
    );
    useEffect(() => {
        if (getHomeDataData) {
            const newCpuData = processChartData(
                cpuDataRef.current,
                [getHomeDataData.getCurrentTime, getHomeDataData.getSystemStatus.cpu ?? 0],
                HomeConstraints.lineChartRetention,
                HomeConstraints.maxGapSeconds
            );
            cpuDataRef.current = newCpuData;
            setCpuData(newCpuData);

            const newMemoryData = processChartData(
                memoryDataRef.current,
                [getHomeDataData.getCurrentTime, getHomeDataData.getSystemStatus.memory ?? 0],
                HomeConstraints.lineChartRetention,
                HomeConstraints.maxGapSeconds
            );
            memoryDataRef.current = newMemoryData;
            setMemoryData(newMemoryData);
        }
    }, [getHomeDataData, processChartData]);

    return (
        <div className="app-usage-grid">
            <div className="card app-panel app-usage-line-card">
                <h2 className="app-section-title">
                    <Icon className="app-icon" path={mdiChip} size={1} />
                    <span>{t('views.Home.cpu_usage.title')}</span>
                </h2>
                {cpuData.length > 1 ? (
                    <LineChart {...lineChartConfig} data={cpuData} />
                ) : (
                    <Skeleton height={lineChartConfig.height} />
                )}
            </div>
            <div className="card app-panel app-usage-line-card">
                <h2 className="app-section-title">
                    <Icon className="app-icon" path={mdiMemory} size={1} />
                    <span>{t('views.Home.memory_usage.title')}</span>
                </h2>
                {memoryData.length > 1 ? (
                    <LineChart {...lineChartConfig} data={memoryData} />
                ) : (
                    <Skeleton height={lineChartConfig.height} />
                )}
            </div>
            <div className="card app-panel app-usage-gauge-card">
                <h2 className="app-section-title">
                    <Icon className="app-icon" path={mdiHarddisk} size={1} />
                    <span>{t('views.Home.disk_usage.title')}</span>
                </h2>
                <GaugeChart
                    height={lineChartConfig.height}
                    valueMax={100}
                    animation={true}
                    value={getHomeDataData?.getSystemStatus.disk ?? 0}
                />
            </div>
        </div>
    );
};
