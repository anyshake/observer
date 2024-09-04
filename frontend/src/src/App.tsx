import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch } from "react-redux";
import { useLocation } from "react-router-dom";

import { Container } from "./components/Container";
import { Footer } from "./components/Footer";
import { Header } from "./components/Header";
import { Navbar } from "./components/Navbar";
import { RouterView } from "./components/RouterView";
import { Scroller } from "./components/Scroller";
import { Sidebar } from "./components/Sidebar";
import { Skeleton } from "./components/Skeleton";
import { apiConfig } from "./config/api";
import { globalConfig } from "./config/global";
import i18n, { i18nConfig } from "./config/i18n";
import { menuConfig } from "./config/menu";
import { routerConfig } from "./config/router";
import { hideLoading } from "./helpers/app/hideLoading";
import { getCurrentLocale } from "./helpers/i18n/getCurrentLocale";
import { setUserLocale } from "./helpers/i18n/setUserLocale";
import { requestRestApi } from "./helpers/request/requestRestApi";
import { onUpdate as UpdateADC } from "./stores/adc";
import { onUpdate as UpdateGeophone } from "./stores/geophone";
import { onUpdate as UpdateStation } from "./stores/station";

const App = () => {
	const { t } = useTranslation();
	const { routes, basename } = routerConfig;
	const { fallback, resources } = i18nConfig;
	const { name, title, author, repository, homepage, footer } = globalConfig;

	useEffect(() => {
		hideLoading();
		// eslint-disable-next-line no-console
		console.log(`%c${process.env.BUILD_TAG ?? "custom build"}`, "color: #0369a1;");
	}, []);

	const { pathname } = useLocation();
	const [currentTitle, setCurrentTitle] = useState(title);
	const [currentLocale, setCurrentLocale] = useState(fallback);

	const setCurrentLocaleToState = async () => {
		setCurrentLocale(await getCurrentLocale(i18n));
	};

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
		void setCurrentLocaleToState();
		const subtitle = getCurrentTitle();
		setCurrentTitle(subtitle);
		document.title = `${subtitle} - ${title}`;
	}, [t, getCurrentTitle, title]);

	const dispatch = useDispatch();

	const getStationAttributes = useCallback(async () => {
		const { backend, endpoints } = apiConfig;
		const res = await requestRestApi({
			backend,
			timeout: 30,
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
		void getStationAttributes();
	}, [getStationAttributes]);

	const handleSwitchLocale = (locale: string) => {
		void setUserLocale(i18n, locale);
	};

	const locales = Object.entries(resources).reduce(
		(acc, [key, value]) => {
			acc[key] = value.label;
			return acc;
		},
		{} as Record<string, string>
	);

	return (
		<Container toaster={true}>
			<Header
				title={name}
				locales={locales}
				currentLocale={currentLocale}
				onSwitchLocale={handleSwitchLocale}
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
		</Container>
	);
};

export default App;
