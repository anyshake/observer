import { ReactNode, Suspense, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Route, Routes, useLocation } from 'react-router-dom';

import { IRoute, IRouterComponent } from '../../config/router';

interface IRouterView {
    readonly routes: Record<string, IRoute>;
    readonly routerProps: IRouterComponent;
    readonly appName: string;
    readonly suspense: ReactNode;
}

export const RouterView = ({ routes, suspense, appName, routerProps }: IRouterView) => {
    const { pathname } = useLocation();
    const { t } = useTranslation();

    // Set the document title based on the current route
    useEffect(() => {
        const routeTitle = Object.values(routes).find(({ uri }) => pathname === uri)?.title;
        const title = routeTitle ?? routes.default.title;
        document.title = `${t(title)} - ${t(appName)}`;
    }, [routes, appName, pathname, t]);

    return (
        <Suspense key={pathname} fallback={suspense}>
            <Routes>
                {Object.values(routes).map(({ uri, element: Element }, index) => (
                    <Route key={index} element={<Element {...routerProps} />} path={uri} />
                ))}
            </Routes>
        </Suspense>
    );
};
