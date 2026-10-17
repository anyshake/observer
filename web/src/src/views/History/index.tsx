import { mdiBlurLinear, mdiLink, mdiLock, mdiLockOpen, mdiLockReset, mdiWaveform } from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';
import { FFTExecutor } from 'spectrogram-js';

import { LineChart } from '../../components/chart/LineChart';
import { Spectrogram } from '../../components/chart/Spectrogram';
import { DraggableBox } from '../../components/ui/DraggableBox';
import { HistoryConstraints } from '../../config/constraints';
import { IRouterComponent } from '../../config/router';
import { useGetSeismicEventSourceListQuery } from '../../graphql';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { useUrlParams } from '../../helpers/request/useUrlParams';
import { setClipboardText } from '../../helpers/utils/setClipboardText';
import { useRetentionStore } from '../../stores/retention';
import { DataExport } from './DataExport';
import { EventListModal } from './EventListModal';
import { EventSearch } from './EventSearch';
import { RecordSearch } from './RecordSearch';
import { useDataExport } from './useDataExport';
import { useHistoryLayout } from './useHistoryLayout';
import { useSeismicEvents } from './useSeismicEvents';
import { useSeismicRecords } from './useSeismicRecords';

const History = ({ currentLocale }: IRouterComponent) => {
    const { t } = useTranslation();
    const sharedFFTExecutor = useMemo(() => new FFTExecutor(HistoryConstraints.fftSize), []);

    const [displayMode, setDisplayMode] = useState<'waveform' | 'spectrogram'>('waveform');
    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

    const [isBusy, setIsBusy] = useState(true);
    const [searchParams, setSearchParams] = useSearchParams();
    const { start_time: startTimeInSearchParam, end_time: endTimeInSearchParam } = useUrlParams<{
        start_time?: string;
        end_time?: string;
    }>(searchParams);
    const [startTime, setStartTime] = useState(parseInt(startTimeInSearchParam ?? '0'));
    const [endTime, setEndTime] = useState(parseInt(endTimeInSearchParam ?? '0'));
    const { data: getSeismicEventSourceListData, loading: getSeismicEventSourceListLoading } =
        useGetSeismicEventSourceListQuery();
    useEffect(() => {
        if (getSeismicEventSourceListData?.getCurrentTime && !startTime && !endTime) {
            const { getCurrentTime: currentTime } = getSeismicEventSourceListData;
            setStartTime(currentTime - 1000 * 60 * 5);
            setEndTime(currentTime);
            setIsBusy(false);
        } else if (startTime && endTime) {
            setIsBusy(false);
        }
    }, [startTime, endTime, getSeismicEventSourceListData]);

    const events = useSeismicEvents({ setIsBusy, setStartTime, setEndTime });

    const dataExport = useDataExport({ startTime, endTime, setIsBusy });

    const { sampleRate, chartData, spectrogramSeries, activeChannels, handleSearchRecords } =
        useSeismicRecords({ startTime, endTime, setIsBusy });
    const {
        locks,
        toggleLock,
        getInitialLayout,
        handleResetLayout,
        handleDragStop,
        handleResizeStop,
        handleSpectrogramUpdate
    } = useHistoryLayout(activeChannels);
    const { retention } = useRetentionStore();
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart

    const handleShareLink = useCallback(() => {
        const params = new URLSearchParams(searchParams);
        params.set('start_time', startTime.toString());
        params.set('end_time', endTime.toString());
        setSearchParams(params);
        setClipboardText(window.location.href);
        sendUserAlert(t('views.History.share_link.link_copied'));
    }, [t, searchParams, startTime, endTime, setSearchParams]);

    return (
        <div className="app-page">
            <div className="app-history-search-grid">
                <RecordSearch
                    currentLocale={currentLocale}
                    startTime={startTime}
                    endTime={endTime}
                    setStartTime={setStartTime}
                    setEndTime={setEndTime}
                    isBusy={isBusy}
                    handleSearchRecords={handleSearchRecords}
                />

                <DataExport isBusy={isBusy} dataExport={dataExport} />

                <EventSearch
                    currentLocale={currentLocale}
                    isBusy={isBusy}
                    getSeismicEventSourceListData={getSeismicEventSourceListData}
                    getSeismicEventSourceListLoading={getSeismicEventSourceListLoading}
                    events={events}
                />
            </div>

            <div className="app-toolbar">
                <button
                    className="btn btn-sm app-button-toolbar"
                    onClick={() => toggleLock(HistoryConstraints.id)}
                >
                    {locks[HistoryConstraints.id] ? (
                        <>
                            <Icon className="app-icon" path={mdiLock} size={0.7} />
                            <span>{t('views.History.layout_locker.unlock_button')}</span>
                        </>
                    ) : (
                        <>
                            <Icon className="app-icon" path={mdiLockOpen} size={0.7} />
                            <span>{t('views.History.layout_locker.lock_button')}</span>
                        </>
                    )}
                </button>
                <button className="btn btn-sm app-button-toolbar" onClick={handleResetLayout}>
                    <Icon className="app-icon" path={mdiLockReset} size={0.7} />
                    <span>{t('views.History.reset_layout.reset_button')}</span>
                </button>
                {sampleRate >= HistoryConstraints.freqRange[1] * 2 && (
                    <button
                        className="btn btn-sm app-button-toolbar"
                        onClick={handleToggleDisplayMode}
                    >
                        <Icon
                            className="app-icon"
                            path={displayMode === 'waveform' ? mdiBlurLinear : mdiWaveform}
                            size={0.7}
                        />
                        <span>
                            {displayMode === 'waveform'
                                ? t('views.History.display_mode.spectrogram_mode')
                                : t('views.History.display_mode.waveform_mode')}
                        </span>
                    </button>
                )}
                <button className="btn btn-sm" onClick={handleShareLink}>
                    <Icon className="app-icon" path={mdiLink} size={0.7} />
                    <span>{t('views.History.share_link.share_button')}</span>
                </button>
            </div>

            <div className="app-chart-workspace">
                {Object.keys(activeChannels)
                    .sort((a, b) => {
                        // Move active chart to the end
                        if (a === activeChart) {
                            return 1;
                        }
                        // Move active chart to the end
                        if (b === activeChart) {
                            return -1;
                        }
                        return 0;
                    })
                    .map((channel) => {
                        const initialLayout = getInitialLayout(
                            activeChannels[channel].id,
                            activeChannels[channel].index
                        );
                        return (
                            <DraggableBox
                                key={channel}
                                layout={initialLayout}
                                locked={locks[HistoryConstraints.id]}
                                constraints={HistoryConstraints}
                                onDragStart={() => setActiveChart(channel)}
                                onDragStop={(x, y) =>
                                    handleDragStop(
                                        activeChannels[channel].id,
                                        activeChannels[channel].index,
                                        x,
                                        y
                                    )
                                }
                                onResizeStop={(width, height) =>
                                    handleResizeStop(
                                        activeChannels[channel].id,
                                        activeChannels[channel].index,
                                        width,
                                        height
                                    )
                                }
                            >
                                <div
                                    className={
                                        displayMode === 'waveform'
                                            ? 'app-chart-visible'
                                            : 'app-hidden'
                                    }
                                >
                                    <LineChart
                                        minSpanValue={HistoryConstraints.minSpanValue}
                                        height={'100%'}
                                        yPosition="right"
                                        title={channel}
                                        zoom={true}
                                        animation={false}
                                        data={chartData[channel]}
                                    />
                                </div>
                                <div
                                    className={
                                        displayMode === 'spectrogram'
                                            ? 'app-chart-visible'
                                            : 'app-hidden'
                                    }
                                >
                                    <Spectrogram
                                        title={channel}
                                        duration={retention}
                                        overlap={HistoryConstraints.overlap}
                                        freqRange={HistoryConstraints.freqRange}
                                        windowSize={HistoryConstraints.windowSize}
                                        maxDB={initialLayout.spectrogram.maxDB}
                                        minDB={initialLayout.spectrogram.minDB}
                                        colorMap={initialLayout.spectrogram.colorMap}
                                        fftExecutor={sharedFFTExecutor}
                                        data={spectrogramSeries[channel] ?? []}
                                        sampleRate={sampleRate}
                                        paused={displayMode !== 'spectrogram'}
                                        onSpectrogramUpdate={(minDB, maxDB, colorMap) =>
                                            handleSpectrogramUpdate(
                                                activeChannels[channel].id,
                                                activeChannels[channel].index,
                                                minDB,
                                                maxDB,
                                                colorMap
                                            )
                                        }
                                    />
                                </div>
                            </DraggableBox>
                        );
                    })}
            </div>

            <EventListModal currentLocale={currentLocale} events={events} />
        </div>
    );
};

export default History;
