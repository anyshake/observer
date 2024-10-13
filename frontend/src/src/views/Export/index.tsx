import "@fancyapps/ui/dist/fancybox/fancybox.css";

import { Fancybox } from "@fancyapps/ui";
import { mdiMagnify } from "@mdi/js";
import { Icon } from "@mdi/react";
import { GridValidRowModel } from "@mui/x-data-grid";
import { Field, Form, Formik } from "formik";
import { useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "../../components/Button";
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
	// miniSEED table data and its fetch handler
	const [miniSeedTableData, setMiniSeedTableData] = useState<GridValidRowModel[]>([]);
	const getMiniSeedList = async () => {
		const { endpoints, backend } = apiConfig;
		const miniSeedRes = await requestRestApi({
			backend,
			throwError: true,
			endpoint: endpoints.miniseed,
			payload: { action: "list", name: "" }
		});
		setMiniSeedTableData(
			miniSeedRes?.data
				?.sort((a, b) => Math.floor(b.time / 1000) - Math.floor(a.time / 1000))
				.map((item, index) => ({ id: index + 1, ...item })) ?? []
		);
	};

	// miniSEED search result and its filter
	const [miniSeedSearchResult, setMiniSeedSearchResult] = useState<GridValidRowModel[]>([]);
	const handleSearchMiniSeed = (keyword: string) => {
		const date = new Date(keyword);
		const result = isNaN(date.getTime())
			? miniSeedTableData.filter(({ name }) => name.toUpperCase().includes(keyword))
			: miniSeedTableData.filter(({ time }) => time >= date.getTime());
		if (result.length !== 0) {
			sendUserAlert(t("views.export.toasts.search_n_results", { count: result.length }));
			setMiniSeedSearchResult(result);
		} else {
			sendUserAlert(t("views.export.toasts.search_no_result"), true);
		}
	};
	const handleResetMiniSeedFilter = () => {
		if (miniSeedSearchResult.length) {
			sendUserAlert(t("views.export.toasts.search_filter_reset"));
			setMiniSeedSearchResult([]);
		}
	};

	// HeliCorder table data and its fetch handler
	const [heliCorderTableData, setHeliCorderTableData] = useState<GridValidRowModel[]>([]);
	const getHeliCorderList = async () => {
		const { endpoints, backend } = apiConfig;
		const heliCorderRes = await requestRestApi({
			backend,
			throwError: true,
			endpoint: endpoints.helicorder,
			payload: { action: "list", name: "" }
		});
		setHeliCorderTableData(
			heliCorderRes?.data
				?.sort((a, b) => Math.floor(b.time / 1000) - Math.floor(a.time / 1000))
				.map((item, index) => ({ id: index + 1, ...item })) ?? []
		);
	};

	// HeliCorder search result and its filter
	const [heliCorderSearchResult, setHeliCorderSearchResult] = useState<GridValidRowModel[]>([]);
	const handleSearchHeliCorder = (keyword: string) => {
		const date = new Date(keyword);
		const result = isNaN(date.getTime())
			? heliCorderTableData.filter(({ name }) => name.toUpperCase().includes(keyword))
			: heliCorderTableData.filter(({ time }) => time >= date.getTime());
		if (result.length !== 0) {
			sendUserAlert(t("views.export.toasts.search_n_results", { count: result.length }));
			setHeliCorderSearchResult(result);
		} else {
			sendUserAlert(t("views.export.toasts.search_no_result"), true);
		}
	};
	const handleResetHeliCorderFilter = () => {
		if (heliCorderSearchResult.length) {
			sendUserAlert(t("views.export.toasts.search_filter_reset"));
			setHeliCorderSearchResult([]);
		}
	};

	// Update file list every 10 seconds
	const { t } = useTranslation();
	const getFileList = useCallback(async () => {
		await getMiniSeedList();
		await getHeliCorderList();
	}, []);
	useEffect(() => {
		sendPromiseAlert(
			getFileList(),
			t("views.export.toasts.is_fetching_files"),
			t("views.export.toasts.fetch_files_success"),
			t("views.export.toasts.fetch_files_error")
		);
	}, [getFileList, t]);
	useInterval(
		async () => {
			try {
				await getFileList();
			} catch {
				/* empty */
			}
		},
		10000,
		false
	);

	// Ref for storing any HTTP request tasks, with progress and abort controller
	const requestTasksRef = useRef<
		Record<string, { fileName: string; progress: number; abortController?: AbortController }>
	>({});
	useEffect(
		() => () => {
			Object.values(requestTasksRef.current).forEach(({ abortController }) => {
				abortController?.abort();
			});
		},
		[]
	);

	// Handler for cancelling request task
	const handleRequestTask = (fileName: string) => {
		const { abortController } = requestTasksRef.current[fileName];
		abortController?.abort();
		delete requestTasksRef.current[fileName];
		forceUpdateComponent();
	};

	// Little hack to force update the component
	const [, forceUpdateKey] = useState(false);
	const forceUpdateComponent = () => {
		forceUpdateKey((prev) => !prev);
	};

	// Handler for setting request task progress
	const handleRequestProgress = (name: string, progress: number, onComplete?: () => void) => {
		requestTasksRef.current[name] = {
			...requestTasksRef.current[name],
			progress
		};
		if (progress === 100) {
			delete requestTasksRef.current[name];
			onComplete?.();
		}
		forceUpdateComponent();
	};

	// Handler for making export miniSEED request
	const handleExportMiniSeed = async (fileName: string) => {
		if (fileName in requestTasksRef.current) {
			return;
		}

		const abortController = new AbortController();
		requestTasksRef.current[fileName] = {
			fileName,
			abortController,
			progress: 0
		};

		forceUpdateComponent();
		sendUserAlert(t("views.export.toasts.is_exporting_mseed"));

		const { backend, endpoints } = apiConfig;
		await requestRestApi<
			typeof endpoints.miniseed.model.request,
			typeof endpoints.miniseed.model.response.common,
			typeof endpoints.miniseed.model.response.error
		>({
			blobOptions: {
				fileName,
				onDownload: ({ progress }) => {
					handleRequestProgress(fileName, (progress ?? 0) * 100, () => {
						sendUserAlert(t("views.export.toasts.export_mseed_success"));
					});
				}
			},
			payload: { action: "export", name: fileName },
			endpoint: endpoints.miniseed,
			abortController,
			timeout: 3600,
			backend
		});
	};

	// Handler for previewing helicorder
	const handlePreviewHeliCorder = async (fileName: string) => {
		if (fileName in requestTasksRef.current) {
			return;
		}

		const abortController = new AbortController();
		requestTasksRef.current[fileName] = {
			fileName,
			abortController,
			progress: 0
		};

		forceUpdateComponent();
		sendUserAlert(t("views.export.toasts.is_previewing_helicorder"));

		const { backend, endpoints } = apiConfig;
		await requestRestApi<
			typeof endpoints.helicorder.model.request,
			typeof endpoints.helicorder.model.response.common,
			typeof endpoints.helicorder.model.response.error
		>({
			payload: { action: "export", name: fileName },
			endpoint: endpoints.helicorder,
			timeout: 3600,
			backend,
			blobOptions: {
				onDownload: ({ progress }) => {
					handleRequestProgress(fileName, (progress ?? 0) * 100);
				},
				onComplete: async (response) => {
					sendUserAlert(t("views.export.toasts.preview_helicorder_success"));
					const imageData = await response.text();
					const blob = new Blob([imageData], { type: "image/svg+xml" });
					const blobUrl = URL.createObjectURL(blob);
					Fancybox.show([{ src: blobUrl, type: "image" }], {
						on: { close: () => URL.revokeObjectURL(blobUrl) }
					});
				}
			},
			abortController
		});
	};
	useEffect(
		() => () => {
			Fancybox.destroy();
		},
		[]
	);

	// Handler for downloading helicorder
	const handleDownloadHeliCorder = async (fileName: string) => {
		if (fileName in requestTasksRef.current) {
			return;
		}

		const abortController = new AbortController();
		requestTasksRef.current[fileName] = {
			fileName,
			abortController,
			progress: 0
		};

		forceUpdateComponent();
		sendUserAlert(t("views.export.toasts.is_downloading_helicorder"));

		const { backend, endpoints } = apiConfig;
		await requestRestApi<
			typeof endpoints.helicorder.model.request,
			typeof endpoints.helicorder.model.response.common,
			typeof endpoints.helicorder.model.response.error
		>({
			blobOptions: {
				fileName,
				onDownload: ({ progress }) => {
					handleRequestProgress(fileName, (progress ?? 0) * 100);
				},
				onComplete: () => {
					sendUserAlert(t("views.export.toasts.download_helicorder_success"));
				}
			},
			payload: { action: "export", name: fileName },
			endpoint: endpoints.helicorder,
			timeout: 3600,
			backend,
			abortController
		});
	};

	return (
		<Container>
			<Panel label={t("views.export.panels.miniseed_list")}>
				{Object.values(requestTasksRef.current)
					.filter(({ fileName }) => fileName.endsWith(".mseed"))
					.map(({ fileName, progress }) => (
						<Progress
							key={fileName}
							value={progress}
							label={fileName}
							onCancel={() => {
								handleRequestTask(fileName);
							}}
						/>
					))}
				<Container className="flex flex-col sm:flex-row justify-between gap-6">
					<Container className="flex flex-row space-x-4 sm:whitespace-nowrap">
						<Button
							className="px-4 bg-indigo-700 hover:bg-indigo-800"
							label={t("views.export.buttons.refresh_list")}
							onClick={async () => {
								await sendPromiseAlert(
									getMiniSeedList(),
									t("views.export.toasts.is_fetching_files"),
									t("views.export.toasts.fetch_files_success"),
									t("views.export.toasts.fetch_files_error")
								);
							}}
						/>
						<Button
							className="px-4 bg-yellow-700 hover:bg-yellow-800"
							label={t("views.export.buttons.reset_filter")}
							onClick={handleResetMiniSeedFilter}
						/>
					</Container>
					<Formik
						initialValues={{ keyword: "" }}
						onSubmit={({ keyword }, { setSubmitting }) => {
							handleSearchMiniSeed(keyword.toLocaleUpperCase());
							setSubmitting(false);
						}}
					>
						{({ isSubmitting }) => (
							<Form className="flex flex-row space-x-2">
								<Field
									type="search"
									name="keyword"
									className="ps-3 w-full min-w-32 md:w-64 py-2 text-sm text-gray-900 border focus:outline-none border-gray-300 rounded-lg bg-gray-50 focus:border-blue-500"
									placeholder={t(
										"views.export.forms.search_miniseed.placeholder"
									)}
									required
								/>
								<button
									className="text-white bg-blue-700 hover:bg-blue-800 focus:outline-none font-medium rounded-lg text-sm p-2 disabled:cursor-not-allowed"
									disabled={isSubmitting}
									type="submit"
								>
									<Icon className="text-white" path={mdiMagnify} size={1} />
								</button>
							</Form>
						)}
					</Formik>
				</Container>
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
							minWidth: 200,
							renderCell: ({ value }) => `${(value / 1024 / 1024).toFixed(2)} MB`
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
										className="text-blue-700 hover:opacity-50"
										onClick={() => {
											handleExportMiniSeed(row.name);
										}}
									>
										{t("views.export.table.actions.export")}
									</button>
								</div>
							)
						}
					]}
					data={miniSeedSearchResult.length ? miniSeedSearchResult : miniSeedTableData}
				/>
			</Panel>

			<Panel label={t("views.export.panels.helicorder_list")}>
				{Object.values(requestTasksRef.current)
					.filter(({ fileName }) => fileName.endsWith(".svg"))
					.map(({ fileName, progress }) => (
						<Progress
							key={fileName}
							value={progress}
							label={fileName}
							onCancel={() => {
								handleRequestTask(fileName);
							}}
						/>
					))}
				<Container className="flex flex-col sm:flex-row justify-between gap-6">
					<Container className="flex flex-row space-x-4 sm:whitespace-nowrap">
						<Button
							className="px-4 bg-indigo-700 hover:bg-indigo-800"
							label={t("views.export.buttons.refresh_list")}
							onClick={async () => {
								await sendPromiseAlert(
									getHeliCorderList(),
									t("views.export.toasts.is_fetching_files"),
									t("views.export.toasts.fetch_files_success"),
									t("views.export.toasts.fetch_files_error")
								);
							}}
						/>
						<Button
							className="px-4 bg-yellow-700 hover:bg-yellow-800"
							label={t("views.export.buttons.reset_filter")}
							onClick={handleResetHeliCorderFilter}
						/>
					</Container>
					<Formik
						initialValues={{ keyword: "" }}
						onSubmit={({ keyword }, { setSubmitting }) => {
							handleSearchHeliCorder(keyword.toLocaleUpperCase());
							setSubmitting(false);
						}}
					>
						{({ isSubmitting }) => (
							<Form className="flex flex-row space-x-2">
								<Field
									type="search"
									name="keyword"
									className="ps-3 w-full min-w-32 md:w-64 py-2 text-sm text-gray-900 border focus:outline-none border-gray-300 rounded-lg bg-gray-50 focus:border-blue-500"
									placeholder={t(
										"views.export.forms.search_helicorder.placeholder"
									)}
									required
								/>
								<button
									className="text-white bg-blue-700 hover:bg-blue-800 focus:outline-none font-medium rounded-lg text-sm p-2 disabled:cursor-not-allowed"
									disabled={isSubmitting}
									type="submit"
								>
									<Icon className="text-white" path={mdiMagnify} size={1} />
								</button>
							</Form>
						)}
					</Formik>
				</Container>
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
							minWidth: 200,
							renderCell: ({ value }) => `${(value / 1024 / 1024).toFixed(2)} MB`
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
							minWidth: 200,
							renderCell: ({ value }) => {
								if (value === -1) {
									return "∞";
								}
								if (value < 0) {
									return "0";
								}
								return String(value);
							}
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
										className="text-blue-700 hover:opacity-50"
										onClick={() => {
											handlePreviewHeliCorder(row.name);
										}}
									>
										{t("views.export.table.actions.preview")}
									</button>
									<button
										className="text-blue-700 hover:opacity-50"
										onClick={() => {
											handleDownloadHeliCorder(row.name);
										}}
									>
										{t("views.export.table.actions.download")}
									</button>
								</div>
							)
						}
					]}
					data={
						heliCorderSearchResult.length ? heliCorderSearchResult : heliCorderTableData
					}
				/>
			</Panel>
		</Container>
	);
};

export default Export;
