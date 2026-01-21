import { DataGrid, GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import * as dataGridLocales from '@mui/x-data-grid/locales';
import { Localization } from '@mui/x-data-grid/utils/getGridLocalization';
import { useEffect, useMemo, useState } from 'react';

import { localeConfig } from '../../config/locale';

interface TableListProps {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly columns: GridColDef<GridValidRowModel>[];
    readonly data: GridValidRowModel[];
    readonly onSelect?: (rows: GridValidRowModel[]) => void;
    readonly sortDirection?: 'asc' | 'desc';
    readonly sortField?: string;
}

export const TableList = ({
    currentLocale,
    columns,
    data,
    onSelect = () => {},
    sortField = 'id',
    sortDirection
}: TableListProps) => {
    const themeRecords = useMemo(
        () =>
            Object.entries(dataGridLocales).reduce(
                (acc, [locale, value]) => {
                    acc[locale] = value;
                    return acc;
                },
                {} as Record<string, object>
            ),
        []
    );

    const [locale4Component, setLocale4Component] = useState('enUS');
    useEffect(() => {
        const componentLocale = currentLocale.replace(/[^a-z0-9]/gi, '');
        setLocale4Component(themeRecords[componentLocale] ? componentLocale : 'enUS');
    }, [currentLocale, themeRecords]);

    return (
        <DataGrid
            localeText={
                (themeRecords[locale4Component] as Localization).components.MuiDataGrid.defaultProps
                    .localeText
            }
            columns={columns}
            rows={data}
            sx={{
                minHeight: 300,
                minWidth: 10,
                '& .MuiDataGrid-cell': {
                    paddingLeft: !onSelect ? 4 : 0
                },
                '& .MuiDataGrid-columnHeaderDraggableContainer': {
                    paddingLeft: !onSelect ? 3 : 0
                }
            }}
            initialState={{
                pagination: { paginationModel: { page: 0, pageSize: 6 } },
                sorting: { sortModel: [{ field: sortField, sort: sortDirection ?? 'asc' }] }
            }}
            onCellClick={({ field }, event) => {
                if (field === 'actions') {
                    event.stopPropagation();
                }
            }}
            onRowSelectionModelChange={(idx) => {
                const rows: GridValidRowModel[] = [];
                idx.forEach((id) => {
                    const row = data.find((v) => String(v[sortField]) === String(id));
                    if (row) {
                        rows.push(row);
                    }
                });
                onSelect(rows);
            }}
            pageSizeOptions={[6, 12, 18, 24]}
            checkboxSelection={!!onSelect}
        />
    );
};
