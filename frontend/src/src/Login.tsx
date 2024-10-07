import { mdiAccount, mdiEarth, mdiKey, mdiShieldCheck } from "@mdi/js";
import Icon from "@mdi/react";
import { ErrorMessage, Field, Form, Formik } from "formik";
import { md, pki, util } from "node-forge";
import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";

import { Container } from "./components/Container";
import { apiConfig, authCommonResponseModel0, authCommonResponseModel2 } from "./config/api";
import { globalConfig } from "./config/global";
import { ReduxStoreProps } from "./config/store";
import { sendPromiseAlert } from "./helpers/interact/sendPromiseAlert";
import { requestRestApi } from "./helpers/request/requestRestApi";
import logo from "./logo.png";
import { onUpdate as UpdateCredential } from "./stores/credential";

interface LoginProps {
	currentLocale: string;
	locales: Record<string, string>;
	onSwitchLocale: (newLocale: string) => void;
	onLoginStateChange: (alive: boolean) => void;
}

export const Login = ({
	currentLocale,
	locales,
	onSwitchLocale,
	onLoginStateChange
}: LoginProps) => {
	// Set page title
	useEffect(() => {
		document.title = globalConfig.title;
	}, []);

	// State for pre-authentication data (captcha, encrypt key, etc.)
	const [preAuthTTL, setPreAuthTTL] = useState(0);
	const [preAuthData, setPreAuthData] = useState({
		encrypt_key: "",
		captcha_id: "",
		captcha_img: ""
	});
	const { t } = useTranslation();
	const getPreAuthData = useCallback(
		async (notify: boolean) => {
			const request = (throwError: boolean) => {
				const { backend, endpoints } = apiConfig;
				return requestRestApi({
					backend,
					throwError,
					endpoint: endpoints.auth,
					payload: { action: "preauth", nonce: "", credential: "" }
				});
			};
			const res = (
				notify
					? await sendPromiseAlert(
							request(true),
							t("login.toasts.is_refreshing_captcha"),
							t("login.toasts.refresh_captcha_success"),
							t("login.toasts.refresh_captcha_error")
						)
					: await request(false)
			) as typeof authCommonResponseModel0;
			if (res.data) {
				const { ttl, encrypt_key, captcha_id, captcha_img } = res.data;
				setPreAuthData({ encrypt_key, captcha_id, captcha_img });
				setPreAuthTTL(ttl);
			}
		},
		[t]
	);
	useEffect(() => {
		// Refresh captcha and encrypt key after TTL
		if (preAuthTTL) {
			const interval = setInterval(() => {
				getPreAuthData(true);
			}, preAuthTTL);
			return () => clearInterval(interval);
		}
		getPreAuthData(false);
	}, [preAuthData, preAuthTTL, getPreAuthData]);

	// Handle login submission
	const dispatch = useDispatch();
	const handleLoginSubmit = async (username: string, password: string, captcha: string) => {
		// Load public key and encrypt credential
		const publicKey = util.decode64(preAuthData.encrypt_key);
		const credential = pki.publicKeyFromPem(publicKey).encrypt(
			JSON.stringify({
				username,
				password,
				timestamp: Date.now(),
				captcha_solution: captcha,
				captcha_id: preAuthData.captcha_id
			}),
			"RSA-OAEP"
		);

		// Send login request
		const { backend, endpoints } = apiConfig;
		const res = (await requestRestApi({
			backend,
			endpoint: endpoints.auth,
			payload: {
				action: "login",
				credential: util.encode64(credential),
				nonce: md.sha1.create().update(publicKey).digest().toHex()
			}
		})) as typeof authCommonResponseModel2;
		if (!res.data) {
			throw new Error(res.message);
		}

		// Save token and expires_at to redux store
		const { token, expires_at } = res.data;
		dispatch(UpdateCredential({ token, expires_at }));
	};

	// Handle credential state change, e.g. login success or token expired
	const { credential } = useSelector(({ credential }: ReduxStoreProps) => credential);
	useEffect(() => {
		if (credential.token.length && credential.expires_at > Date.now()) {
			onLoginStateChange(true);
		}
	}, [onLoginStateChange, credential, t]);

	return (
		<Container className="min-h-screen p-20 px-4 flex flex-col from-purple-500 to-blue-500 bg-gradient-to-br">
			<div className="bg-white w-full max-w-md md:max-w-xl p-12 rounded-lg shadow-xl m-auto border">
				<img src={logo} alt="Login" className="size-24 md:size-32 mx-auto mb-6" />
				<Formik
					enableReinitialize
					initialValues={{ username: "", password: "", captcha: "" }}
					validate={(values) => {
						const errors: {
							username?: string;
							password?: string;
							captcha?: string;
						} = {};
						if (!values.username.length) {
							errors.username = t("login.forms.username.error");
						}
						if (!values.password.length) {
							errors.password = t("login.forms.password.error");
						}
						if (!values.captcha.length) {
							errors.captcha = t("login.forms.captcha.error");
						}
						return errors;
					}}
					onSubmit={({ username, password, captcha }, { setSubmitting }) => {
						sendPromiseAlert(
							handleLoginSubmit(username, password, captcha),
							t("login.toasts.is_logging_in"),
							t("login.toasts.login_success"),
							t("login.toasts.login_error"),
							false
						).catch(() => {
							getPreAuthData(false);
							setSubmitting(false);
						});
					}}
				>
					{({ isSubmitting }) => (
						<Form className="space-y-6">
							<div>
								<label
									htmlFor="username"
									className="flex items-center text-gray-600"
								>
									<Icon className="mr-2" path={mdiAccount} size={0.8} />
									{t("login.forms.username.title")}
								</label>
								<ErrorMessage
									className="text-red-500 text-sm"
									name="username"
									component="div"
								/>
								<Field
									id="username"
									className="py-2 px-3 border border-gray-300 focus:outline-none focus:ring focus:ring-opacity-50 rounded-md shadow-sm disabled:bg-gray-100 mt-1 block w-full"
									type="text"
									name="username"
									placeholder={t("login.forms.username.placeholder")}
								/>
							</div>

							<div>
								<label
									htmlFor="password"
									className="flex items-center text-gray-600"
								>
									<Icon className="mr-2" path={mdiKey} size={0.8} />
									{t("login.forms.password.title")}
								</label>
								<ErrorMessage
									className="text-red-500 text-sm"
									name="password"
									component="div"
								/>
								<Field
									id="password"
									className="py-2 px-3 border border-gray-300 focus:outline-none focus:ring focus:ring-opacity-50 rounded-md shadow-sm disabled:bg-gray-100 mt-1 block w-full"
									type="password"
									name="password"
									placeholder={t("login.forms.password.placeholder")}
								/>
							</div>

							<div>
								<label
									htmlFor="captcha"
									className="flex items-center text-gray-600"
								>
									<Icon className="mr-2" path={mdiShieldCheck} size={0.8} />
									{t("login.forms.captcha.title")}
								</label>
								<ErrorMessage
									className="text-red-500 text-sm"
									name="captcha"
									component="div"
								/>
								<div className="flex justify-between space-x-2">
									<Field
										id="captcha"
										className="py-2 px-3 border border-gray-300 focus:outline-none focus:ring focus:ring-opacity-50 rounded-md shadow-sm disabled:bg-gray-100 mt-1 block w-full"
										type="text"
										name="captcha"
                                        disabled={!preAuthData.captcha_img.length}
										placeholder={t("login.forms.captcha.placeholder")}
									/>
									{preAuthData.captcha_img && (
										<img
											className="self-center w-24 md:w-32 cursor-pointer"
											src={`data:image/png;base64,${preAuthData.captcha_img}`}
											onClick={() => {
												getPreAuthData(true);
											}}
											alt=""
										/>
									)}
								</div>
							</div>

							<button
								className="bg-purple-500 hover:bg-purple-700 w-full text-white font-medium text-sm shadow-lg rounded-lg py-2 transition-all disabled:cursor-wait"
								type="submit"
								disabled={isSubmitting}
							>
								{t("login.forms.submit.title")}
							</button>

							<div className="flex items-center justify-center text-gray-500 text-xs">
								<Icon path={mdiEarth} size={0.6} />
								<select
									className="bg-transparent focus:outline-none text-center truncate cursor-pointer hover:opacity-80"
									onChange={({ target }) => onSwitchLocale(target.value)}
									value={currentLocale}
								>
									<option disabled>Choose Language</option>
									{Object.entries(locales).map(([key, value]) => (
										<option key={key} value={key}>
											{value}
										</option>
									))}
								</select>
							</div>
						</Form>
					)}
				</Formik>
			</div>
		</Container>
	);
};
