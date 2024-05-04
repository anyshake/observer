import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import ExportIcon from "../../assets/icons/download-solid.svg";
import { Container } from "../../components/Container";
import { Panel } from "../../components/Panel";
import { Progress } from "../../components/Progress";
import { Table, TableData, TableProps } from "../../components/Table";
import { apiConfig } from "../../config/api";
import { useInterval } from "../../helpers/hook/useInterval";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { getExportsUpdates } from "./getExportsUpdates";
import { handleSetTable } from "./handleSetTable";

const Export = () => {
	const { t } = useTranslation();

	// Little hack to force update the component
	const [, forceUpdateKey] = useState(false);
	const forceUpdate = () => {
		forceUpdateKey((prev) => !prev);
	};

	const exportTasksRef = useRef<
		Record<
			string,
			{
				fileName: string;
				progress: number;
				abortController?: AbortController;
			}
		>
	>({});

	const handleExportProgress = (name: string, progress: number) => {
		exportTasksRef.current[name] = {
			...exportTasksRef.current[name],
			progress
		};
		if (progress === 100) {
			delete exportTasksRef.current[name];
			sendUserAlert(t("views.export.toasts.export_mseed_success"));
		}
		forceUpdate();
	};

	const handleExportRequest = async (data: TableData) => {
		const { name } = data;
		if (name in exportTasksRef.current) {
			return;
		}

		const abortController = new AbortController();
		exportTasksRef.current[name] = {
			fileName: name as string,
			abortController,
			progress: 0
		};

		forceUpdate();
		sendUserAlert(t("views.export.toasts.is_exporting_mseed"));

		const { backend, endpoints } = apiConfig;
		await requestRestApi<
			typeof endpoints.mseed.model.request,
			typeof endpoints.mseed.model.response.common,
			typeof endpoints.mseed.model.response.error
		>({
			blobOptions: {
				fileName: String(name),
				onDownload: ({ progress }) => {
					handleExportProgress(name as string, (progress ?? 0) * 100);
				}
			},
			payload: { action: "export", name: name as string },
			endpoint: endpoints.mseed,
			abortController,
			timeout: 3600,
			backend
		});
	};

	const handleExportCancel = (fileName: string) => {
		const { abortController } = exportTasksRef.current[fileName];
		abortController?.abort();
		delete exportTasksRef.current[fileName];
		forceUpdate();
	};

	useEffect(
		() => () => {
			Object.values(exportTasksRef.current).forEach(({ abortController }) => {
				abortController?.abort();
			});
		},
		[]
	);

	const [table, setTable] = useState<TableProps>({
		columns: [
			{ key: "name", label: "views.export.table.columns.name" },
			{ key: "size", label: "views.export.table.columns.size" },
			{ key: "time", label: "views.export.table.columns.time" },
			{ key: "ttl", label: "views.export.table.columns.ttl" }
		],
		actions: [
			{
				icon: ExportIcon,
				onClick: handleExportRequest,
				label: "views.export.table.actions.export"
			}
		],
		rowsLimit: 10,
		loadMoreText: "views.export.table.load_more",
		placeholder: "views.export.table.placeholder.is_fetching_mseed"
	});

	useInterval(
		() => {
			getExportsUpdates((res) => {
				handleSetTable(res, setTable);
			});
		},
		5000,
		true
	);

	return (
		<Container>
			<Panel label={t("views.export.panels.file_list")}>
				{Object.values(exportTasksRef.current).map(({ fileName, progress }) => (
					<Progress
						key={fileName}
						value={progress}
						label={fileName}
						onCancel={() => {
							handleExportCancel(fileName);
						}}
					/>
				))}
				<Table
					{...table}
					placeholder={t(table.placeholder)}
					loadMoreText={t(table.loadMoreText ?? "")}
					columns={table.columns.map((column) => ({
						...column,
						label: t(column.label)
					}))}
					actions={table.actions.map((action) => ({
						...action,
						label: t(action.label)
					}))}
				/>
			</Panel>
		</Container>
	);
};

export default Export;
