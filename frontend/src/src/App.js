import React, { Suspense, Component } from "react";
import {
    HashRouter,
    BrowserRouter,
    Routes as Switch,
    Route,
} from "react-router-dom";
import { registerEvents } from "./helpers/events/appEvents";
import withRouter from "./helpers/router/withRouter";
import RouterConfig, { Routes } from "./router";
import Notfound from "./components/Notfound";
import Loader from "./app/loader";
import AppConfig from "./config";

const RouteModule = (props) => {
    return (
        <Switch location={props.location}>
            <Route element={<Notfound />} path="*" />
            {RouterConfig.map((item, index) => {
                const Element = item;
                return (
                    <Route
                        {...(Routes[index].index ? index : null)}
                        path={Routes[index].path}
                        element={<Element />}
                        key={index}
                    />
                );
            })}
        </Switch>
    );
};

export default class App extends Component {
    render() {
        const Routes = withRouter(RouteModule);
        registerEvents({
            eventArray: [
                { trigger: "selectstart", id: "globalApp_userSelectStart" },
                { trigger: "contextmenu", id: "globalApp_userContextMenu" },
            ],
            onEventCallback: (e) => e.preventDefault(),
        });

        if (
            AppConfig.frontend.router === "hash" ||
            AppConfig.frontend.router === "redirect"
        ) {
            return (
                <HashRouter>
                    <Suspense fallback={<Loader />}>
                        <Routes />
                    </Suspense>
                </HashRouter>
            );
        }

        return (
            <BrowserRouter>
                <Suspense fallback={<Loader />}>
                    <Routes />
                </Suspense>
            </BrowserRouter>
        );
    }
}
