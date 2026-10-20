import {
    DataGrid,
    GridColDef,
    GridRowParams,
    GridRowSelectionModel,
    GridValidRowModel
} from '@mui/x-data-grid';
import * as dataGridLocales from '@mui/x-data-grid/locales';

import { localeConfig } from '../../config/locale';

interface TableListProps {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly columns: GridColDef<GridValidRowModel>[];
    readonly data: GridValidRowModel[];
    readonly onSelect?: (rows: GridValidRowModel[]) => void;
    readonly rowSelectionModel?: GridRowSelectionModel;
    readonly isRowSelectable?: (params: GridRowParams) => boolean;
    readonly sortDirection?: 'asc' | 'desc';
    readonly sortField?: string;
}

export const TableList = ({
    currentLocale,
    columns,
    data,
    onSelect,
    rowSelectionModel,
    isRowSelectable,
    sortField = 'id',
    sortDirection
}: TableListProps) => {
    const locale =
        localeConfig.resources[currentLocale] ?? localeConfig.resources[localeConfig.fallback];
    const locale4Component = dataGridLocales[locale.key4mui.dataGrid];

    return (
        <DataGrid
            localeText={locale4Component.components.MuiDataGrid.defaultProps.localeText}
            columns={columns}
            rows={data}
            className={`app-data-grid ${!onSelect ? 'app-data-grid-unselectable' : ''}`}
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
                    const row = data.find((v) => v.id === id);
                    if (row) {
                        rows.push(row);
                    }
                });
                onSelect?.(rows);
            }}
            rowSelectionModel={rowSelectionModel}
            isRowSelectable={isRowSelectable}
            pageSizeOptions={[6, 12, 18, 24]}
            checkboxSelection={!!onSelect}
        />
    );
};
