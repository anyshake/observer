import { RefObject, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { getTimeString } from '../../helpers/utils/getTimeString';

export const useRecording = (updatedAtRef: RefObject<number>) => {
    const { t } = useTranslation();
    const [recordingState, setRecordingState] = useState({
        isRecording: false,
        startTime: 0,
        endTime: 0
    });
    const [recordList, setRecordList] = useState<[number, number, string][]>([]);
    const handleToggleRecording = useCallback(() => {
        const stamp = updatedAtRef.current;
        setRecordingState((prev) => {
            const { isRecording, startTime } = prev;
            if (!isRecording) {
                if (stamp === 0) {
                    return prev;
                }
                sendUserAlert(
                    t('views.RealTime.record_data.start_recording', {
                        startedAt: getTimeString(stamp)
                    })
                );
                return { ...prev, isRecording: true, startTime: stamp };
            }

            const endTime = stamp;
            if (startTime !== endTime) {
                const search = new URLSearchParams({
                    start_time: startTime.toString(),
                    end_time: endTime.toString()
                });
                setRecordList((prevList) => [...prevList, [startTime, endTime, search.toString()]]);
                sendUserAlert(
                    t('views.RealTime.record_data.link_created', {
                        endedAt: getTimeString(endTime),
                        duration: ((endTime - startTime) / 1000).toFixed(1)
                    })
                );
            } else {
                sendUserAlert(t('views.RealTime.record_data.no_data_recorded'), true);
            }

            return { ...prev, isRecording: false, endTime };
        });
    }, [t, updatedAtRef]);

    const handleRemoveRecord = useCallback(
        (index: number) => {
            setRecordList((prevList) => prevList.filter((_, i) => i !== index));
        },
        [setRecordList]
    );

    return { recordingState, recordList, handleToggleRecording, handleRemoveRecord };
};
