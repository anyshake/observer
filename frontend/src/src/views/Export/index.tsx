import { GridValidRowModel } from "@mui/x-data-grid";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { Container } from "../../components/Container";
import { Panel } from "../../components/Panel";
import { Progress } from "../../components/Progress";
import { TableList } from "../../components/TableList";
import { apiConfig } from "../../config/api";
import { i18nConfig } from "../../config/i18n";
import { RouterComponentProps } from "../../config/router";
import { useInterval } from "../../helpers/hook/useInterval";
import { sendPromiseAlert } from "../../helpers/interact/sendPromiseAlert";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { getTimeString } from "../../helpers/utils/getTimeString";

const Export = ({ locale }: RouterComponentProps) => {
	// Update MiniSEED file list every 5 seconds
	const { t } = useTranslation();
	const [tableData, setTableData] = useState<GridValidRowModel[]>([]);
	const getFileList = async () => {
		const { endpoints, backend } = apiConfig;
		const res = await requestRestApi({
			backend,
			endpoint: endpoints.mseed,
			payload: { action: "list", name: "" }
		});
		const tableData = res?.data
			?.sort((a, b) => Math.floor(b.time / 1000) - Math.floor(a.time / 1000))
			.map((item, index) => ({ id: index + 1, ...item }));
		setTableData(tableData ?? []);
	};
	useEffect(() => {
		sendPromiseAlert(
			getFileList(),
			t("views.export.toasts.is_fetching_files"),
			t("views.export.toasts.fetch_files_success"),
			t("views.export.toasts.fetch_files_error")
		);
	}, [t]);
	useInterval(getFileList, 5000, false);

	// Ref for storing export tasks, with progress and abort controller
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
	useEffect(
		() => () => {
			Object.values(exportTasksRef.current).forEach(({ abortController }) => {
				abortController?.abort();
			});
		},
		[]
	);

	// Little hack to force update the component
	const [, forceUpdateKey] = useState(false);
	const forceUpdateComponent = () => {
		forceUpdateKey((prev) => !prev);
	};

	// Handler for displaying export progress
	const handleExportProgress = (name: string, progress: number) => {
		exportTasksRef.current[name] = {
			...exportTasksRef.current[name],
			progress
		};
		if (progress === 100) {
			delete exportTasksRef.current[name];
			sendUserAlert(t("views.export.toasts.export_mseed_success"));
		}
		forceUpdateComponent();
	};

	// Handler for making export request
	const handleExportRequest = async (filename: string) => {
		if (filename in exportTasksRef.current) {
			return;
		}

		const abortController = new AbortController();
		exportTasksRef.current[filename] = {
			fileName: filename,
			abortController,
			progress: 0
		};

		forceUpdateComponent();
		sendUserAlert(t("views.export.toasts.is_exporting_mseed"));

		const { backend, endpoints } = apiConfig;
		await requestRestApi<
			typeof endpoints.mseed.model.request,
			typeof endpoints.mseed.model.response.common,
			typeof endpoints.mseed.model.response.error
		>({
			blobOptions: {
				fileName: filename,
				onDownload: ({ progress }) => {
					handleExportProgress(filename, (progress ?? 0) * 100);
				}
			},
			payload: { action: "export", name: filename },
			endpoint: endpoints.mseed,
			abortController,
			timeout: 3600,
			backend
		});
	};

	// Handler for cancelling export request
	const handleExportCancel = (fileName: string) => {
		const { abortController } = exportTasksRef.current[fileName];
		abortController?.abort();
		delete exportTasksRef.current[fileName];
		forceUpdateComponent();
	};

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
				<TableList
					locale={locale ?? i18nConfig.fallback}
					columns={[
						{
							field: "id",
							headerName: t("views.export.table.columns.id"),
							hideable: false,
							sortable: true,
							minWidth: 120
						},
						{
							field: "name",
							headerName: t("views.export.table.columns.name"),
							hideable: false,
							sortable: true,
							minWidth: 350
						},
						{
							field: "size",
							headerName: t("views.export.table.columns.size"),
							hideable: true,
							sortable: true,
							minWidth: 200
						},
						{
							field: "time",
							headerName: t("views.export.table.columns.time"),
							hideable: true,
							sortable: true,
							minWidth: 230,
							renderCell: ({ value }) => getTimeString(value)
						},
						{
							field: "ttl",
							headerName: t("views.export.table.columns.ttl"),
							hideable: true,
							sortable: false,
							minWidth: 200
						},
						{
							field: "actions",
							headerName: t("views.export.table.columns.actions"),
							sortable: false,
							resizable: false,
							minWidth: 150,
							renderCell: ({ row }) => (
								<div className="flex flex-row space-x-4 w-full">
									<button
										className="text-blue-700 dark:text-blue-400 hover:opacity-50"
										onClick={() => {
											handleExportRequest(row.name);
										}}
									>
										{t("views.export.table.actions.export")}
									</button>
								</div>
							)
						}
					]}
					data={tableData}
				/>
			</Panel>
		</Container>
	);
};

export default Export;
