import { GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { getTimeString } from '../../helpers/utils/getTimeString';
import type { useAssetActions } from './useAssetActions';

export const useFileColumns = ({
    handleDownloadAsset,
    handlePreviewAsset
}: ReturnType<typeof useAssetActions>) => {
    const { t } = useTranslation();
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

    return { getColumns };
};
