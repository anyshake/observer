import { Footer } from "./components/Footer";
import { Header } from "./components/Header";
import { RouterView } from "./components/RouterView";
import { Sidebar } from "./components/Sidebar";
import { Skeleton } from "./components/Skeleton";
import { Container } from "./components/Container";
import { routerConfig } from "./config/router";
import { useCallback, useEffect, useState } from "react";
import { getCurrentLocale } from "./helpers/i18n/getCurrentLocale";
import i18n, { i18nConfig } from "./config/i18n";
import { globalConfig } from "./config/global";
import { setUserLocale } from "./helpers/i18n/setUserLocale";
import { useTranslation } from "react-i18next";
import { menuConfig } from "./config/menu";
import { getAsciiArt } from "./helpers/app/getAsciiArt";
import { Navbar } from "./components/Navbar";
import { useLocation } from "react-router-dom";
import { useDispatch } from "react-redux";
import { requestRestApi } from "./helpers/request/requestRestApi";
import { apiConfig } from "./config/api";
import { onUpdate as UpdateADC } from "./stores/adc";
import { onUpdate as UpdateGeophone } from "./stores/geophone";
import { Scroller } from "./components/Scroller";
import { onUpdate as UpdateStation } from "./stores/station";

const App = () => {
    const { t } = useTranslation();
    const { routes, basename } = routerConfig;
    const { fallback, resources } = i18nConfig;
    const {
        name,
        title,
        author,
        repository,
        homepage,
        footer,
        version,
        release,
    } = globalConfig;

    useEffect(() => {
        document.querySelector(".public-loading")?.remove();
        const asciiArt = getAsciiArt();
        console.info(`%c${asciiArt}`, "color: #0891b2;");
        console.info(`%c${version}-${release}`, "color: #0369a1;");
    }, [version, release]);

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
            endpoint: endpoints.station,
        });
        if (!!res?.data) {
            const initialized = true;
            const { sensitivity, frequency } = res.data.geophone;
            dispatch(UpdateGeophone({ sensitivity, frequency, initialized }));
            const { resolution, fullscale } = res.data.adc;
            dispatch(UpdateADC({ resolution, fullscale, initialized }));
            const { station, network, location } = res.data.station;
            dispatch(
                UpdateStation({ station, network, location, initialized })
            );
        }
    }, [dispatch]);

    useEffect(() => {
        void getStationAttributes();
    }, [getStationAttributes]);

    const handleSwitchLocale = (locale: string) => {
        void setUserLocale(i18n, locale);
    };

    const locales = Object.entries(resources).reduce((acc, [key, value]) => {
        acc[key] = value.label;
        return acc;
    }, {} as Record<string, string>);

    return (
        <Container toaster={true}>
            <Header
                title={name}
                locales={locales}
                currentLocale={currentLocale}
                onSwitchLocale={handleSwitchLocale}
            />
            <Sidebar
                title={name}
                links={menuConfig}
                currentLocale={currentLocale}
            />

            <Container main={true}>
                <Navbar
                    pathname={pathname}
                    basename={basename}
                    title={currentTitle}
                />
                <RouterView
                    routes={routes}
                    suspense={<Skeleton />}
                    routerProps={{
                        locale: currentLocale,
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
