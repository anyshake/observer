import { ReactNode } from 'react';
import { BrowserRouter, HashRouter } from 'react-router-dom';

export type RouterMode = 'hash' | 'history';

export interface IRouterWrapper {
    readonly mode: 'hash' | 'history';
    readonly basename: string;
    readonly children: ReactNode;
}

export const RouterWrapper = ({ mode, children }: IRouterWrapper) =>
    mode === 'hash' ? (
        <HashRouter>{children}</HashRouter>
    ) : (
        <BrowserRouter>{children}</BrowserRouter>
    );
