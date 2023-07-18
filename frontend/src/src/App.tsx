import { Suspense, Component, ReactElement } from "react";
import {
    HashRouter,
    Routes as Switch,
    Route,
    BrowserRouter,
} from "react-router-dom";
import withRouter from "./helpers/withRouter";
import ROUTER_CONFIG from "./config/router";
import GLOBAL_CONFIG from "./config/global";
import Loader from "./views/Loader";
import NotFound from "./views/NotFound";

export default class App extends Component {
    render() {
        const { router: routerMode } = GLOBAL_CONFIG.app_settings;
        const router = (props: any): ReactElement => {
            const { location } = props;
            return (
                <Switch location={location}>
                    <Route element={<NotFound />} path="*" />
                    {ROUTER_CONFIG.map(({ node, uri }, index) => (
                        <Route key={index} element={node} path={uri} />
                    ))}
                </Switch>
            );
        };

        const Routes = withRouter(router);
        return routerMode === "hash" ? (
            <HashRouter>
                <Suspense fallback={<Loader />}>
                    <Routes />
                </Suspense>
            </HashRouter>
        ) : (
            <BrowserRouter>
                <Suspense fallback={<Loader />}>
                    <Routes />
                </Suspense>
            </BrowserRouter>
        );
    }
}
