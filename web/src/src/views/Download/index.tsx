import '@fancyapps/ui/dist/fancybox/fancybox.css';

import { Fancybox } from '@fancyapps/ui';
import { mdiArchive, mdiClose, mdiImageAlbum, mdiMagnify } from '@mdi/js';
import Icon from '@mdi/react';
import { GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Card } from '../../components/Card';
import { TableList } from '../../components/TableList';
import { IRouterComponent } from '../../config/router';
import { useGetFileListDataQuery } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { getRestfulApiUrl } from '../../helpers/app/getRestfulApiUrl';
import { ApiClient } from '../../helpers/request/ApiClient';
import { getTimeString } from '../../helpers/utils/getTimeString';

const Download = ({ currentLocale }: IRouterComponent) => {
    const { t } = useTranslation();
    const { loading: getFileListDataLoading, data: getFileListDataData } = useGetFileListDataQuery({
        pollInterval: 5000
    });

    const [miniSeedList, setMiniSeedList] = useState<
        Array<{
            id: number;
            namespace: string;
            filePath: string;
            fileName: string;
            size: number;
            timestamp: number;
        }>
    >([]);
    useEffect(() => {
        if (getFileListDataData?.getMiniSeedFiles) {
            const rawList = [...getFileListDataData.getMiniSeedFiles]
                .sort((a, b) => b!.modifiedAt - a!.modifiedAt)
                .map((file, index) => ({
                    id: index + 1,
                    namespace: file!.namespace,
                    filePath: file!.filePath,
                    fileName: file!.fileName,
                    size: file!.size,
                    timestamp: file!.modifiedAt
                }));
            setMiniSeedList(rawList.sort((a, b) => b.timestamp - a.timestamp));
        }
    }, [getFileListDataData]);

    const [helicorderList, setHelicorderList] = useState<
        Array<{
            id: number;
            namespace: string;
            filePath: string;
            fileName: string;
            size: number;
            timestamp: number;
        }>
    >([]);
    useEffect(() => {
        if (getFileListDataData?.getHelicorderFiles) {
            const rawList = [...getFileListDataData.getHelicorderFiles]
                .sort((a, b) => b!.modifiedAt - a!.modifiedAt)
                .map((file, index) => ({
                    id: index + 1,
                    namespace: file!.namespace,
                    filePath: file!.filePath,
                    fileName: file!.fileName,
                    size: file!.size,
                    timestamp: file!.modifiedAt
                }));
            setHelicorderList(rawList.sort((a, b) => b.timestamp - a.timestamp));
        }
    }, [getFileListDataData]);

    const [searchMiniSeed, setSearchMiniSeed] = useState('');
    const [searchHelicorder, setSearchHelicorder] = useState('');
    const filteredMiniSeedList = miniSeedList.filter((item) =>
        item.fileName.toLowerCase().includes(searchMiniSeed.toLowerCase())
    );
    const filteredHelicorderList = helicorderList.filter((item) =>
        item.fileName.toLowerCase().includes(searchHelicorder.toLowerCase())
    );

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

    const getColumns = useCallback(
        (hasPreview: boolean): GridColDef<GridValidRowModel>[] => [
            {
                field: 'id',
                headerName: t('views.Download.file_list_columns.id'),
                hideable: false,
                sortable: true,
                minWidth: 150
            },
            {
                field: 'fileName',
                headerName: t('views.Download.file_list_columns.name'),
                hideable: false,
                sortable: true,
                minWidth: 350
            },
            {
                field: 'timestamp',
                headerName: t('views.Download.file_list_columns.time'),
                hideable: true,
                sortable: true,
                minWidth: 250,
                renderCell: (cell: GridValidRowModel) => getTimeString(cell.value)
            },
            {
                field: 'size',
                headerName: t('views.Download.file_list_columns.size'),
                hideable: true,
                sortable: true,
                minWidth: 200,
                renderCell: (cell: GridValidRowModel) =>
                    `${(cell.value / 1024 / 1024).toFixed(2)} MB`
            },
            {
                field: 'actions',
                headerName: t('views.Download.file_list_columns.actions'),
                sortable: false,
                minWidth: 220,
                headerAlign: 'center',
                align: 'center',
                renderCell: ({ row: { fileName, namespace, filePath } }: GridValidRowModel) => (
                    <div className="space-x-4">
                        {hasPreview && (
                            <button
                                className="cursor-pointer text-blue-700 hover:opacity-50"
                                onClick={() => {
                                    handlePreviewAsset(namespace, fileName, filePath);
                                }}
                            >
                                {t('views.Download.file_list_actions.preview')}
                            </button>
                        )}
                        <button
                            className="cursor-pointer text-blue-700 hover:opacity-50"
                            onClick={() => {
                                handleDownloadAsset(namespace, fileName, filePath);
                            }}
                        >
                            {t('views.Download.file_list_actions.download')}
                        </button>
                    </div>
                )
            }
        ],
        [t, handlePreviewAsset, handleDownloadAsset]
    );

    return (
        <div className="container mx-auto space-y-6 p-4">
            <div className="flex w-full flex-col space-y-2">
                <Card title={t('views.Download.miniseed_list.title')} iconPath={mdiArchive}>
                    {getFileListDataLoading ? (
                        <div className="flex items-center justify-center">
                            <span className="loading loading-spinner loading-lg text-gray-400" />
                        </div>
                    ) : (
                        <div className="space-y-4">
                            <label className="input">
                                <Icon className="flex-shrink-0" path={mdiMagnify} size={0.8} />
                                <input
                                    type="search"
                                    className="grow"
                                    value={searchMiniSeed}
                                    onChange={({ currentTarget }) =>
                                        setSearchMiniSeed(currentTarget.value)
                                    }
                                    placeholder={t('views.Download.miniseed_list.search')}
                                />
                                {searchMiniSeed.length > 0 && (
                                    <button
                                        className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                                        onClick={() => setSearchMiniSeed('')}
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiClose}
                                            size={0.8}
                                        />
                                    </button>
                                )}
                            </label>
                            <TableList
                                sortField="timestamp"
                                sortDirection="desc"
                                currentLocale={currentLocale}
                                data={filteredMiniSeedList}
                                columns={getColumns(false)}
                            />
                        </div>
                    )}
                </Card>
                <Card title={t('views.Download.helicorder_list.title')} iconPath={mdiImageAlbum}>
                    {getFileListDataLoading ? (
                        <div className="flex items-center justify-center">
                            <span className="loading loading-spinner loading-lg text-gray-400" />
                        </div>
                    ) : (
                        <div className="space-y-4">
                            <label className="input">
                                <Icon className="flex-shrink-0" path={mdiMagnify} size={0.8} />
                                <input
                                    type="search"
                                    className="grow"
                                    value={searchHelicorder}
                                    onChange={({ currentTarget }) =>
                                        setSearchHelicorder(currentTarget.value)
                                    }
                                    placeholder={t('views.Download.helicorder_list.search')}
                                />
                                {searchHelicorder.length > 0 && (
                                    <button
                                        className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                                        onClick={() => setSearchHelicorder('')}
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiClose}
                                            size={0.8}
                                        />
                                    </button>
                                )}
                            </label>
                            <TableList
                                sortField="timestamp"
                                sortDirection="desc"
                                currentLocale={currentLocale}
                                data={filteredHelicorderList}
                                columns={getColumns(true)}
                            />
                        </div>
                    )}
                </Card>
            </div>
        </div>
    );
};

export default Download;
