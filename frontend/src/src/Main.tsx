import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";
import { useLocation } from "react-router-dom";

import { Container } from "./components/Container";
import { Footer } from "./components/Footer";
import { Header } from "./components/Header";
import { Navbar } from "./components/Navbar";
import { RouterView } from "./components/RouterView";
import { Scroller } from "./components/Scroller";
import { Sidebar } from "./components/Sidebar";
import { Skeleton } from "./components/Skeleton";
import { apiConfig, authCommonResponseModel2 } from "./config/api";
import { globalConfig } from "./config/global";
import { menuConfig } from "./config/menu";
import { routerConfig } from "./config/router";
import { ReduxStoreProps } from "./config/store";
import { sendUserAlert } from "./helpers/interact/sendUserAlert";
import { sendUserConfirm } from "./helpers/interact/sendUserConfirm";
import { requestRestApi } from "./helpers/request/requestRestApi";
import { onUpdate as UpdateADC } from "./stores/adc";
import { onUpdate as UpdateCredential } from "./stores/credential";
import { onUpdate as UpdateGeophone } from "./stores/geophone";
import { onUpdate as UpdateStation } from "./stores/station";

interface MainProps {
	currentLocale: string;
	locales: Record<string, string>;
	onSwitchLocale: (newLocale: string) => void;
	onLoginStateChange: (alive: boolean) => void;
}

export const Main = ({ currentLocale, locales, onSwitchLocale, onLoginStateChange }: MainProps) => {
	const dispatch = useDispatch();
	const { name, title, author, repository, homepage, footer } = globalConfig;
	const { routes, basename } = routerConfig;
	const { pathname } = useLocation();

	// Refresh token before 1 hour of expiration
	const { credential } = useSelector(({ credential }: ReduxStoreProps) => credential);
	const refreshToken = useCallback(async () => {
		if (credential?.token.length && credential.expires_at > Date.now()) {
			const { backend, endpoints } = apiConfig;
			const res = (await requestRestApi({
				backend,
				endpoint: endpoints.auth,
				payload: { action: "refresh", nonce: "", credential: "" }
			})) as typeof authCommonResponseModel2;
			if (res.data) {
				const { token, expires_at } = res.data;
				dispatch(UpdateCredential({ token, expires_at }));
			} else {
				onLoginStateChange(false);
			}
		} else if (credential?.token.length) {
			onLoginStateChange(false);
		}
	}, [credential, onLoginStateChange, dispatch]);
	useEffect(() => {
		const refreshThreshold = 1800 * 1000; // refresh before 30 minutes of expiration
		let nextRefresh = credential.expires_at - Date.now() - refreshThreshold;
		if (nextRefresh < 0) {
			// Attempt to refresh in 5 seconds
			nextRefresh = 5000;
		}
		const interval = setInterval(refreshToken, nextRefresh);
		return () => clearInterval(interval);
	}, [credential, refreshToken]);

	// Get page title in the current locale
	const [currentTitle, setCurrentTitle] = useState(title);
	const getCurrentTitle = useCallback(() => {
		for (const key in routes) {
			const { prefix, uri, suffix } = routes[key];
			const fullPath = `${prefix}${uri}${suffix}`;
			if (pathname === fullPath) {
				return routes[key].title[currentLocale];
			}
		}
		return routes.default.title[currentLocale];
	}, [routes, pathname, currentLocale]);
	useEffect(() => {
		const subtitle = getCurrentTitle();
		document.title = `${subtitle} - ${title}`;
		setCurrentTitle(subtitle);
	}, [getCurrentTitle, title]);

	// Fetch station attributes
	const getStationAttributes = useCallback(async () => {
		const { backend, endpoints } = apiConfig;
		const res = await requestRestApi({
			backend,
			endpoint: endpoints.station
		});
		if (res?.data) {
			const initialized = true;
			const { sensitivity, frequency } = res.data.sensor;
			dispatch(UpdateGeophone({ sensitivity, frequency, initialized }));
			const { resolution, fullscale } = res.data.sensor;
			dispatch(UpdateADC({ resolution, fullscale, initialized }));
			const { station, network, location, channel } = res.data.stream;
			dispatch(UpdateStation({ station, network, location, initialized, channel }));
		}
	}, [dispatch]);
	useEffect(() => {
		getStationAttributes();
	}, [getStationAttributes]);

	// Handler for logout button
	const { t } = useTranslation();
	const handleLogout = () => {
		sendUserConfirm(t("main.toasts.confirm_logout"), {
			title: t("main.toasts.confirm_title"),
			confirmText: t("main.toasts.confirm_button"),
			cancelText: t("main.toasts.cancel_button"),
			onConfirmed: () => {
				onLoginStateChange(false);
				sendUserAlert(t("main.toasts.logout_success"));
			}
		});
	};

	return (
		<>
			<Header
				title={name}
				locales={locales}
				currentLocale={currentLocale}
				onSwitchLocale={onSwitchLocale}
				{...(credential?.token.length && { onLogout: handleLogout })}
			/>
			<Sidebar title={name} links={menuConfig} currentLocale={currentLocale} />

			<Container main={true}>
				<Navbar pathname={pathname} basename={basename} title={currentTitle} />
				<RouterView
					routes={routes}
					suspense={<Skeleton />}
					routerProps={{
						locale: currentLocale
					}}
				/>
			</Container>

			<Scroller threshold={100} />
			<Footer
				text={footer}
				author={author}
				homepage={homepage}
				repository={repository}
				currentLocale={currentLocale}
			/>
		</>
	);
};
