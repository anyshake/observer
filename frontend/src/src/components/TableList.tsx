import { DataGrid, GridColDef, GridRowSelectionModel, GridValidRowModel } from "@mui/x-data-grid";
import * as dataGridLocales from "@mui/x-data-grid/locales";
import { Localization } from "@mui/x-data-grid/utils/getGridLocalization";

import { i18nConfig } from "../config/i18n";

interface TableListProps {
	readonly locale: keyof typeof i18nConfig.resources;
	readonly columns: GridColDef<GridValidRowModel>[];
	readonly data: GridValidRowModel[];
	readonly onSelect?: (rows: GridRowSelectionModel) => void;
	readonly sortDirection?: "asc" | "desc";
	readonly sortField?: string;
}

export const TableList = ({
	locale,
	columns,
	data,
	onSelect,
	sortField,
	sortDirection
}: TableListProps) => {
	// Get locale of MUI DataGrid component
	const themeRecords = Object.entries(dataGridLocales).reduce(
		(acc, [locale, value]) => {
			acc[locale] = value;
			return acc;
		},
		{} as Record<string, object>
	);
	let locale4Component = locale.replaceAll(/[^a-z0-9]/gi, "");
	if (!themeRecords[locale4Component]) {
		locale4Component = "enUS";
	}

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
				"& .MuiDataGrid-cell": {
					paddingLeft: !onSelect ? 4 : 0
				},
				"& .MuiDataGrid-columnHeaderDraggableContainer": {
					paddingLeft: !onSelect ? 3 : 0
				}
			}}
			initialState={{
				pagination: { paginationModel: { page: 0, pageSize: 10 } },
				sorting: { sortModel: [{ field: sortField ?? "id", sort: sortDirection ?? "asc" }] }
			}}
			onCellClick={({ field }, event) => {
				if (field === "actions") {
					event.stopPropagation();
				}
			}}
			onRowSelectionModelChange={onSelect}
			pageSizeOptions={[5, 10, 20]}
			checkboxSelection={!!onSelect}
		/>
	);
};
