import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch } from "react-redux";

import { Container } from "./components/Container";
import { apiConfig, authCommonResponseModel1 } from "./config/api";
import i18n, { i18nConfig } from "./config/i18n";
import { hideLoading } from "./helpers/app/hideLoading";
import { getCurrentLocale } from "./helpers/i18n/getCurrentLocale";
import { setUserLocale } from "./helpers/i18n/setUserLocale";
import { requestRestApi } from "./helpers/request/requestRestApi";
import { Login } from "./Login";
import { Main } from "./Main";
import { initialCredential, onUpdate as UpdateCredential } from "./stores/credential";

const App = () => {
	// Check if the user needs to login before loading the main page
	const [appInspection, setAppInspection] = useState({
		hasRestrict: false,
		hasLoggedIn: false,
		loadPage: false
	});
	const dispatch = useDispatch();
	const getAppInspection = useCallback(async () => {
		const { backend, endpoints } = apiConfig;
		const res = (await requestRestApi({
			backend,
			endpoint: endpoints.auth,
			payload: { action: "inspect", nonce: "", credential: "" }
		})) as typeof authCommonResponseModel1;
		if (res.data) {
			// Clear credential store if not restricted
			if (!res.data?.restrict) {
				dispatch(UpdateCredential(initialCredential));
			}
			setAppInspection({
				hasLoggedIn: !(res.data?.restrict ?? false),
				hasRestrict: res.data?.restrict,
				loadPage: true
			});
		}
	}, [dispatch]);
	useEffect(() => {
		getAppInspection();
	}, [getAppInspection]);

	// Remove spinner after inspection
	useEffect(() => {
		if (appInspection.loadPage && appInspection.hasLoggedIn) {
			// eslint-disable-next-line no-console
			console.log(`%c${process.env.BUILD_TAG ?? "custom build"}`, "color: #0369a1;");
		}
		if (appInspection.loadPage) {
			hideLoading();
		}
	}, [appInspection]);

	// Handler for login state change
	const handleLoginStateChange = (alive: boolean) => {
		if (!alive) {
			dispatch(UpdateCredential(initialCredential));
		}
		setAppInspection({ ...appInspection, hasLoggedIn: alive, loadPage: true });
	};

	// Get current locale from i18n
	const [currentLocale, setCurrentLocale] = useState(i18nConfig.fallback);
	const { t } = useTranslation();
	useEffect(() => {
		const setCurrentLocaleToState = async () => {
			setCurrentLocale(await getCurrentLocale(i18n));
		};
		setCurrentLocaleToState();
	}, [t]);

	// Locale resources and switcher
	const locales = Object.entries(i18nConfig.resources).reduce(
		(acc, [key, value]) => {
			acc[key] = value.label;
			return acc;
		},
		{} as Record<string, string>
	);
	const handleSwitchLocale = (newLocale: string) => {
		setUserLocale(i18n, newLocale);
	};

	return (
		<Container toaster={true}>
			{appInspection.loadPage ? (
				appInspection.hasLoggedIn ? (
					<Main
						onLoginStateChange={handleLoginStateChange}
						onSwitchLocale={handleSwitchLocale}
						currentLocale={currentLocale}
						locales={locales}
					/>
				) : (
					<Login
						onLoginStateChange={handleLoginStateChange}
						onSwitchLocale={handleSwitchLocale}
						currentLocale={currentLocale}
						locales={locales}
					/>
				)
			) : (
				<></>
			)}
		</Container>
	);
};

export default App;
