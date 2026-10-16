import { Fancybox } from '@fancyapps/ui';
import { useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { getRestfulApiUrl } from '../../helpers/app/getRestfulApiUrl';
import { ApiClient } from '../../helpers/request/ApiClient';

export const useAssetActions = () => {
    const { t } = useTranslation();
    const getAssetFetchLink = useCallback(
        (apiUrl: string, token: string, namespace: string, filePath: string) => {
            const urlObj = new URL(apiUrl);
            urlObj.searchParams.set('token', token);
            urlObj.searchParams.set('namespace', namespace);
            urlObj.searchParams.set('file_path', filePath);
            return urlObj.toString();
        },
        []
    );

    const handleDownloadAsset = useCallback(
        async (namespace: string, fileName: string, filePath: string) => {
            const apiUrl = getRestfulApiUrl('/files');
            await sendPromiseAlert(
                ApiClient.request<string>({
                    url: apiUrl,
                    method: 'post',
                    data: { file_path: filePath }
                }),
                t('views.Download.request_file.requesting', { fileName }),
                (res) => {
                    const fetchUrl = getAssetFetchLink(apiUrl, res!.data!, namespace, filePath);
                    window.open(fetchUrl, '_blank');
                    return t('views.Download.request_file.success', { fileName });
                },
                (error) => t('views.Download.request_file.error', { fileName, error })
            );
        },
        [t, getAssetFetchLink]
    );

    const handlePreviewAsset = useCallback(
        async (namespace: string, fileName: string, filePath: string) => {
            const requestFn = async () => {
                const apiUrl = getRestfulApiUrl('/files');
                const res = await ApiClient.request<string>({
                    url: apiUrl,
                    method: 'post',
                    data: { file_path: filePath }
                });
                const fetchUrl = getAssetFetchLink(apiUrl, res!.data!, namespace, filePath);
                const blobObj = await ApiClient.getBlob({ url: fetchUrl });
                const blobUrl = URL.createObjectURL(blobObj!);
                Fancybox.show([{ src: blobUrl, type: 'image' }], {
                    on: { close: () => URL.revokeObjectURL(blobUrl) }
                });
            };
            await sendPromiseAlert(
                requestFn(),
                t('views.Download.preview_file.requesting', { fileName }),
                t('views.Download.preview_file.success', { fileName }),
                (error) => t('views.Download.preview_file.error', { fileName, error })
            );
        },
        [t, getAssetFetchLink]
    );
    useEffect(
        () => () => {
            Fancybox.destroy();
        },
        []
    );

    return { handleDownloadAsset, handlePreviewAsset };
};
