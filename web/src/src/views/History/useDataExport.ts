import { type Dispatch, type SetStateAction, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { getRestfulApiUrl } from '../../helpers/app/getRestfulApiUrl';
import { ApiClient } from '../../helpers/request/ApiClient';

export const useDataExport = ({
    startTime,
    endTime,
    setIsBusy
}: {
    startTime: number;
    endTime: number;
    setIsBusy: Dispatch<SetStateAction<boolean>>;
}) => {
    const { t } = useTranslation();
    const [availableExports, setAvailableExports] = useState<{
        channels: string[];
        formats: Record<string, string>;
    }>();
    const [selectedExport, setSelectedExport] = useState<{ channel: string; format: string }>({
        channel: '',
        format: ''
    });
    useEffect(() => {
        const requestFn = async () => {
            setIsBusy(true);
            const { data } = await ApiClient.request<{
                channel_code: string[];
                data_format: Record<string, string>;
            }>({
                url: getRestfulApiUrl('/export'),
                method: 'get',
                ignoreErrors: true
            });
            setIsBusy(false);
            if (data) {
                const { channel_code, data_format } = data;
                setAvailableExports({ channels: channel_code, formats: data_format });
            }
        };
        requestFn();
    }, [setIsBusy]);
    const handleGetData = useCallback(async () => {
        setIsBusy(true);
        await sendPromiseAlert(
            ApiClient.saveAs({
                url: getRestfulApiUrl('/export'),
                method: 'post',
                data: {
                    start_time: startTime,
                    end_time: endTime,
                    data_format: selectedExport.format,
                    channel_code: selectedExport.channel
                }
            }),
            t('views.History.data_export.exporting', {
                channel: selectedExport.channel,
                format: availableExports?.formats[selectedExport.format]
            }),
            t('views.History.data_export.success', {
                channel: selectedExport.channel,
                format: availableExports?.formats[selectedExport.format]
            }),
            (error) =>
                t('views.History.data_export.error', {
                    channel: selectedExport.channel,
                    format: availableExports?.formats[selectedExport.format],
                    error
                })
        );
        setIsBusy(false);
    }, [
        startTime,
        endTime,
        selectedExport.format,
        selectedExport.channel,
        setIsBusy,
        t,
        availableExports?.formats
    ]);

    return { availableExports, selectedExport, setSelectedExport, handleGetData };
};
