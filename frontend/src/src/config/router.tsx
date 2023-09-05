import { ReactNode, lazy } from "react";

const Home = lazy(() => import("../views/Home/Index"));
const Realtime = lazy(() => import("../views/Realtime/Index"));
const History = lazy(() => import("../views/History/Index"));
const Setting = lazy(() => import("../views/Setting/Index"));

const ROUTER_CONFIG: RouterConfig[] = [
    {
        uri: "/",
        title: "config.router.index",
        node: <Home />,
    },
    {
        uri: "/realtime",
        title: "config.router.realtime",
        node: <Realtime />,
    },
    {
        uri: "/history",
        title: "config.router.history",
        node: <History />,
    },
    {
        uri: "/setting",
        title: "config.router.setting",
        node: <Setting />,
    },
];

export interface RouterConfig {
    readonly uri: string;
    readonly title: string;
    readonly node: ReactNode;
}

export default ROUTER_CONFIG;
