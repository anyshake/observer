import { ReactNode, Suspense, useEffect } from 'react';
import { Route, Routes, useLocation } from 'react-router-dom';

import { localeConfig, Translation } from '../config/locale';
import { IRoute, IRouterComponent } from '../config/router';

interface IRouterView {
    readonly currentLocale: keyof typeof localeConfig.resources;
    readonly routes: Record<string, IRoute>;
    readonly routerProps: IRouterComponent;
    readonly appName: Translation;
    readonly suspense: ReactNode;
}

export const RouterView = ({
    currentLocale,
    routes,
    suspense,
    appName,
    routerProps
}: IRouterView) => {
    const { pathname } = useLocation();

    // Set the document title based on the current route
    useEffect(() => {
        const routeTitle = Object.values(routes).find(({ uri }) => pathname === uri)?.title;
        const title = routeTitle?.[currentLocale] ?? routes.default.title?.[currentLocale];
        document.title = `${title} - ${appName[currentLocale]}`;
    }, [routes, appName, pathname, currentLocale]);

    return (
        <Suspense key={pathname} fallback={suspense}>
            <Routes>
                {Object.values(routes).map(({ uri, element: Element }, index) => (
                    <Route key={index} element={<Element {...routerProps} />} path={`${uri}`} />
                ))}
            </Routes>
        </Suspense>
    );
};
