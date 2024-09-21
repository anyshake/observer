import { md, pki, util } from "node-forge";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";

import { Button } from "../../components/Button";
import { Code } from "../../components/Code";
import { Container } from "../../components/Container";
import { Dialog, DialogProps } from "../../components/Dialog";
import { Form, FormProps } from "../../components/Form";
import { Input } from "../../components/Input";
import { Panel } from "../../components/Panel";
import { TableList } from "../../components/TableList";
import {
	apiConfig,
	userCommonResponseModel0,
	userCommonResponseModel1,
	userCommonResponseModel2,
	userCommonResponseModel3
} from "../../config/api";
import { globalConfig } from "../../config/global";
import { i18nConfig } from "../../config/i18n";
import { RouterComponentProps } from "../../config/router";
import { ReduxStoreProps } from "../../config/store";
import { sendPromiseAlert } from "../../helpers/interact/sendPromiseAlert";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { sendUserConfirm } from "../../helpers/interact/sendUserConfirm";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { getTimeString } from "../../helpers/utils/getTimeString";
import { onUpdate as updateDuration } from "../../stores/duration";
import { onUpdate as updateRetention } from "../../stores/retention";

const Settings = ({ locale }: RouterComponentProps) => {
	const { retention, duration } = globalConfig;

	// Read the station inventory from the backend
	const [stationInventory, setStationInventory] = useState<string>();
	const getStationInventory = async () => {
		const { backend, endpoints } = apiConfig;
		const payload = { format: "json" };
		const res = await requestRestApi({
			backend,
			payload,
			endpoint: endpoints.inventory
		});
		if (res?.data) {
			setStationInventory(res.data);
		}
	};
	useEffect(() => {
		getStationInventory();
	}, []);

	// Check if current user is an admin if station is restricted
	const { credential } = useSelector(({ credential }: ReduxStoreProps) => credential);
	const [isAdministrator, setIsAdministrator] = useState(false);
	const getUserProfile = async () => {
		const { backend, endpoints } = apiConfig;
		const res = (await requestRestApi({
			backend,
			endpoint: endpoints.user,
			payload: {
				user_id: "",
				action: "profile",
				nonce: "",
				password: "",
				username: "",
				admin: false
			}
		})) as typeof userCommonResponseModel2;
		if (res?.data) {
			setIsAdministrator(res.data.admin);
		}
	};
	useEffect(() => {
		if (credential.token.length) {
			getUserProfile();
		}
	}, [credential]);

	// Basic states for form, other attributes are set on demand
	const [form, setForm] = useState<FormProps & { values?: Record<string, string | number> }>({
		open: false,
		inputType: "number"
	});
	const handleCloseForm = () => {
		setForm({ ...form, open: false });
	};

	// Handler for refreshing the page
	const handleReload = () => {
		setTimeout(() => {
			window.location.reload();
		}, 2500);
	};

	// Handler for changing waveform retention value
	const { t } = useTranslation();
	const dispatch = useDispatch();
	const handleRetentionChange = (newValue?: string) => {
		if (newValue?.length) {
			const value = parseInt(newValue);
			const { maximum, minimum } = retention;
			if (isNaN(value) || value < minimum || value > maximum) {
				return;
			}
			sendUserAlert(t("views.setting.toasts.retention_set", { current: value }));
			dispatch(updateRetention(value));
			handleCloseForm();
			handleReload();
		} else {
			setForm({
				...form,
				open: true,
				values: { ...retention },
				onSubmit: handleRetentionChange,
				cancelText: "views.setting.forms.waveform_retention.cancel",
				submitText: "views.setting.forms.waveform_retention.submit",
				title: "views.setting.forms.waveform_retention.title",
				content: "views.setting.forms.waveform_retention.content",
				placeholder: "views.setting.forms.waveform_retention.placeholder"
			});
		}
	};

	// Handler for changing default query duration value
	const handleDurationChange = (newValue?: string) => {
		if (newValue?.length) {
			const value = parseInt(newValue);
			const { maximum, minimum } = duration;
			if (isNaN(value) || value < minimum || value > maximum) {
				return;
			}
			sendUserAlert(t("views.setting.toasts.duration_set", { current: value }));
			dispatch(updateDuration(value));
			handleCloseForm();
			handleReload();
		} else {
			setForm({
				...form,
				open: true,
				values: { ...duration },
				onSubmit: handleDurationChange,
				cancelText: "views.setting.forms.query_duration.cancel",
				submitText: "views.setting.forms.query_duration.submit",
				title: "views.setting.forms.query_duration.title",
				content: "views.setting.forms.query_duration.content",
				placeholder: "views.setting.forms.query_duration.placeholder"
			});
		}
	};

	// Handler for purging the cache
	const handlePurgeCache = () => {
		sendUserConfirm(t("views.setting.toasts.confirm_purge"), {
			title: t("views.setting.toasts.confirm_title"),
			confirmText: t("views.setting.toasts.confirm_button"),
			cancelText: t("views.setting.toasts.cancel_button"),
			onConfirmed: () => {
				sendUserAlert(t("views.setting.toasts.cache_purged"));
				localStorage.clear();
				handleReload();
			}
		});
	};

	// Read the current values from the Redux store, and set the initial values
	const { retention: currentRetention } = useSelector(
		({ retention }: ReduxStoreProps) => retention
	);
	const { duration: currentDuration } = useSelector(({ duration }: ReduxStoreProps) => duration);
	const [panels] = useState([
		{
			label: "views.setting.panels.waveform_retention",
			content: "views.setting.contents.waveform_retention",
			button: "views.setting.buttons.waveform_retention",
			className: "bg-teal-700 hover:bg-teal-800",
			onClick: handleRetentionChange,
			values: { current: currentRetention, ...retention }
		},
		{
			label: "views.setting.panels.query_duration",
			content: "views.setting.contents.query_duration",
			button: "views.setting.buttons.query_duration",
			className: "bg-lime-700 hover:bg-lime-800",
			onClick: handleDurationChange,
			values: { current: currentDuration, ...duration }
		},
		{
			label: "views.setting.panels.purge_cache",
			content: "views.setting.contents.purge_cache",
			button: "views.setting.buttons.purge_cache",
			className: "bg-pink-700 hover:bg-pink-800",
			onClick: handlePurgeCache
		}
	]);

	// Get all users for admin account
	const [sysUsers, setSysUsers] = useState<
		{
			admin: boolean;
			user_id: number;
			username: string;
			created_at: string;
			updated_at: string;
			last_login: string;
		}[]
	>([]);
	const getSysUsers = async () => {
		const { backend, endpoints } = apiConfig;
		const res = (await requestRestApi({
			backend,
			endpoint: endpoints.user,
			payload: {
				user_id: "",
				action: "list",
				nonce: "",
				password: "",
				username: "",
				admin: false
			}
		})) as typeof userCommonResponseModel3;
		if (res?.data) {
			const users = res.data
				.sort((a, b) => a.user_id - b.user_id)
				.map((user, index) => ({
					...user,
					id: index + 1
				}));
			setSysUsers(users);
		}
	};
	useEffect(() => {
		if (isAdministrator) {
			getSysUsers();
		}
	}, [isAdministrator]);

	// Get encrypted texts by acquring a public key from the backend
	const getEncryptedText = async (textArr: string[]) => {
		const { backend, endpoints } = apiConfig;
		const res = (await requestRestApi({
			backend,
			endpoint: endpoints.user,
			payload: {
				user_id: "",
				action: "preauth",
				nonce: "",
				password: "",
				username: "",
				admin: false
			}
		})) as typeof userCommonResponseModel1;
		if (res?.data) {
			const publicKey = util.decode64(res?.data.encrypt_key);
			const encryptedArr = textArr.map((text) =>
				util.encode64(pki.publicKeyFromPem(publicKey).encrypt(text, "RSA-OAEP"))
			);
			return { encryptedArr, publicKey };
		}

		return { encryptedArr: [], publicKey: "" };
	};

	// Handler for removing a user
	const handleRemoveUser = async (username: string) => {
		const removeFn = async () => {
			const { encryptedArr, publicKey } = await getEncryptedText([username]);
			const { backend, endpoints } = apiConfig;
			const res = (await requestRestApi({
				backend,
				endpoint: endpoints.user,
				payload: {
					user_id: "",
					action: "remove",
					password: "",
					admin: false,
					username: encryptedArr[0],
					nonce: md.sha1.create().update(publicKey).digest().toHex()
				}
			})) as typeof userCommonResponseModel0;
			if (res?.error) {
				throw new Error(res.message);
			}
			await getSysUsers();
		};
		sendUserConfirm(t("views.setting.toasts.confirm_remove_user", { username }), {
			title: t("views.setting.toasts.confirm_title"),
			confirmText: t("views.setting.toasts.confirm_button"),
			cancelText: t("views.setting.toasts.cancel_button"),
			onConfirmed: () => {
				sendPromiseAlert(
					removeFn(),
					t("views.setting.toasts.is_removing_user"),
					t("views.setting.toasts.remove_user_success"),
					t("views.setting.toasts.remove_user_error")
				);
			}
		});
	};

	// Basic states for dialog, other attributes are set on demand
	const [dialog, setDialog] = useState<DialogProps>({ open: false });
	const handleCloseDialog = () => {
		setDialog({ ...dialog, open: false });
	};

	// Handler for creating a new user
	const handleCreateUser = () => {
		const newUserData = { username: "", password: "", admin: false };
		const submitFn = async () => {
			const { encryptedArr, publicKey } = await getEncryptedText([
				newUserData.username,
				newUserData.password
			]);
			const { backend, endpoints } = apiConfig;
			const res = (await requestRestApi({
				backend,
				endpoint: endpoints.user,
				payload: {
					user_id: "",
					action: "create",
					username: encryptedArr[0],
					password: encryptedArr[1],
					admin: newUserData.admin,
					nonce: md.sha1.create().update(publicKey).digest().toHex()
				}
			})) as typeof userCommonResponseModel0;
			if (res?.error) {
				throw new Error(res.message);
			}
			await getSysUsers();
			setDialog({ ...dialog, open: false });
		};
		setDialog({
			...dialog,
			open: true,
			content: t("views.setting.dialogs.create_user.content"),
			title: t("views.setting.dialogs.create_user.title"),
			submitText: t("views.setting.dialogs.create_user.submit"),
			cancelText: t("views.setting.dialogs.create_user.cancel"),
			children: (
				<div className="flex flex-col md:flex-row gap-4">
					<Input
						label={t("views.setting.dialogs.create_user.username")}
						type="text"
						onValueChange={(value) => (newUserData.username = String(value))}
					/>
					<Input
						label={t("views.setting.dialogs.create_user.password")}
						type="password"
						onValueChange={(value) => (newUserData.password = String(value))}
					/>
					<Input
						label={t("views.setting.dialogs.create_user.admin")}
						type="select"
						selectOptions={[
							{ value: "false", label: t("views.setting.dialogs.create_user.no") },
							{ value: "true", label: t("views.setting.dialogs.create_user.yes") }
						]}
						onValueChange={(value) => (newUserData.admin = value === "true")}
					/>
				</div>
			),
			onSubmit: () =>
				sendPromiseAlert(
					submitFn(),
					t("views.setting.toasts.is_creating_user"),
					t("views.setting.toasts.create_user_success"),
					t("views.setting.toasts.create_user_error")
				)
		});
	};

	// Handler for editing a user
	const handleEditUser = (userId: string, username: string) => {
		const newUserData = { userId, username, admin: false, password: "" };
		const submitFn = async () => {
			const { encryptedArr, publicKey } = await getEncryptedText([
				newUserData.userId,
				newUserData.username,
				newUserData.password
			]);
			const { backend, endpoints } = apiConfig;
			const res = (await requestRestApi({
				backend,
				endpoint: endpoints.user,
				payload: {
					action: "edit",
					user_id: encryptedArr[0],
					username: encryptedArr[1],
					password: newUserData.password.length ? encryptedArr[2] : "",
					admin: newUserData.admin,
					nonce: md.sha1.create().update(publicKey).digest().toHex()
				}
			})) as typeof userCommonResponseModel0;
			if (res?.error) {
				throw new Error(res.message);
			}
			await getSysUsers();
			setDialog({ ...dialog, open: false });
		};
		setDialog({
			...dialog,
			open: true,
			content: t("views.setting.dialogs.edit_user.content"),
			title: t("views.setting.dialogs.edit_user.title"),
			submitText: t("views.setting.dialogs.edit_user.submit"),
			cancelText: t("views.setting.dialogs.edit_user.cancel"),
			children: (
				<div className="flex flex-col md:flex-row gap-4">
					<Input
						label={t("views.setting.dialogs.edit_user.username")}
						type="text"
						defaultValue={username}
						onValueChange={(value) => (newUserData.username = String(value))}
					/>
					<Input
						label={t("views.setting.dialogs.edit_user.password")}
						type="password"
						onValueChange={(value) => (newUserData.password = String(value))}
					/>
					<Input
						label={t("views.setting.dialogs.edit_user.admin")}
						type="select"
						selectOptions={[
							{ value: "false", label: t("views.setting.dialogs.edit_user.no") },
							{ value: "true", label: t("views.setting.dialogs.edit_user.yes") }
						]}
						onValueChange={(value) => (newUserData.admin = value === "true")}
					/>
				</div>
			),
			onSubmit: () =>
				sendPromiseAlert(
					submitFn(),
					t("views.setting.toasts.is_editing_user"),
					t("views.setting.toasts.edit_user_success"),
					t("views.setting.toasts.edit_user_error")
				)
		});
	};

	return (
		<>
			<Container className="gap-4 grid md:grid-cols-3">
				{panels.map(({ label, content, button, className, onClick, values }) => (
					<Panel key={label} className="" label={t(label)}>
						{t(content, { ...values })
							.split("\n")
							.map((item) => (
								<div key={item}>{item}</div>
							))}
						<Button label={t(button)} onClick={onClick} className={className} />
					</Panel>
				))}
				<Form
					{...form}
					onClose={handleCloseForm}
					title={t(form.title ?? "")}
					cancelText={t(form.cancelText ?? "")}
					submitText={t(form.submitText ?? "")}
					placeholder={t(form.placeholder ?? "")}
					content={t(form.content ?? "", { ...form.values })}
				/>
			</Container>
			{isAdministrator && (
				<Panel label={t("views.setting.panels.user_management")}>
					<Container className="max-w-xs space-x-2 flex">
						<Button
							label={t("views.setting.buttons.create_user")}
							className="w-1/2 bg-orange-700 hover:bg-orange-800"
							onClick={() => {
								handleCreateUser();
							}}
						/>
						<Button
							label={t("views.setting.buttons.refresh_user")}
							className="w-1/2 bg-fuchsia-700 hover:bg-fuchsia-800"
							onClick={() => {
								sendPromiseAlert(
									getSysUsers(),
									t("views.setting.toasts.is_refreshing_user"),
									t("views.setting.toasts.refresh_user_success"),
									t("views.setting.toasts.refresh_user_error")
								);
							}}
						/>
					</Container>
					<TableList
						locale={locale ?? i18nConfig.fallback}
						columns={[
							{
								field: "id",
								headerName: t("views.setting.table.columns.id"),
								hideable: false,
								sortable: true,
								minWidth: 120
							},
							{
								field: "user_id",
								headerName: t("views.setting.table.columns.user_id"),
								hideable: false,
								sortable: true,
								minWidth: 230
							},
							{
								field: "username",
								headerName: t("views.setting.table.columns.username"),
								hideable: false,
								sortable: true,
								minWidth: 200
							},
							{
								field: "admin",
								headerName: t("views.setting.table.columns.admin"),
								hideable: true,
								sortable: true,
								minWidth: 150,
								renderCell: ({ value }) =>
									value
										? t("views.setting.table.columns.yes")
										: t("views.setting.table.columns.no")
							},
							{
								field: "created_at",
								headerName: t("views.setting.table.columns.created_at"),
								hideable: true,
								sortable: true,
								minWidth: 230,
								renderCell: ({ value }) => getTimeString(value)
							},
							{
								field: "last_login",
								headerName: t("views.setting.table.columns.last_login"),
								hideable: true,
								sortable: true,
								minWidth: 230,
								renderCell: ({ value }) =>
									value
										? getTimeString(value)
										: t("views.setting.table.columns.never")
							},
							{
								field: "user_ip",
								headerName: t("views.setting.table.columns.user_ip"),
								hideable: true,
								sortable: true,
								minWidth: 180,
								renderCell: ({ value }) =>
									value ? value : t("views.setting.table.columns.never")
							},
							{
								field: "user_agent",
								headerName: t("views.setting.table.columns.user_agent"),
								hideable: true,
								sortable: true,
								minWidth: 200,
								renderCell: ({ value }) =>
									value ? value : t("views.setting.table.columns.never")
							},
							{
								field: "updated_at",
								headerName: t("views.setting.table.columns.updated_at"),
								hideable: true,
								sortable: true,
								minWidth: 230,
								renderCell: ({ value }) =>
									value
										? getTimeString(value)
										: t("views.setting.table.columns.never")
							},
							{
								field: "actions",
								headerName: t("views.setting.table.columns.actions"),
								sortable: false,
								resizable: false,
								minWidth: 150,
								renderCell: ({ row }) => (
									<div className="flex flex-row space-x-4 w-full">
										<button
											className="text-blue-700 dark:text-blue-400 hover:opacity-50"
											onClick={() => {
												handleEditUser(String(row.user_id), row.username);
											}}
										>
											{t("views.setting.table.actions.edit")}
										</button>
										<button
											className="text-red-700 dark:text-red-400 hover:opacity-50"
											onClick={() => {
												handleRemoveUser(row.username);
											}}
										>
											{t("views.setting.table.actions.remove")}
										</button>
									</div>
								)
							}
						]}
						data={sysUsers}
					/>
					<Dialog {...dialog} onClose={handleCloseDialog} />
				</Panel>
			)}
			{stationInventory?.length && (
				<Panel label={t("views.setting.panels.station_inventory")}>
					<Code language="xml" fileName="inventory.xml">
						{stationInventory}
					</Code>
				</Panel>
			)}
		</>
	);
};

export default Settings;
