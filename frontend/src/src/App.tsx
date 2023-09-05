import { Suspense, Component, ReactElement } from "react";
import {
    HashRouter,
    Routes as Switch,
    Route,
    BrowserRouter,
} from "react-router-dom";
import withRouter from "./helpers/router/withRouter";
import ROUTER_CONFIG from "./config/router";
import GLOBAL_CONFIG from "./config/global";
import NotFound from "./components/NotFound";
import getAsciiArt from "./helpers/app/getAsciiArt";
import PageLoader from "./components/PageLoader";

export default class App extends Component {
    componentDidMount(): void {
        const asciiArt = getAsciiArt();
        const { version, release } = GLOBAL_CONFIG.app_settings;
        console.info(`%c${asciiArt}`, "color: #0891b2;");
        console.info(`%cRelease: ${version}-${release}`, "color: #0369a1;");
    }

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
                <Suspense fallback={<PageLoader />}>
                    <Routes />
                </Suspense>
            </HashRouter>
        ) : (
            <BrowserRouter>
                <Suspense fallback={<PageLoader />}>
                    <Routes />
                </Suspense>
            </BrowserRouter>
        );
    }
}
