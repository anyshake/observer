import { ReactNode, lazy } from "react";

const Home = lazy(() => import("../views/Home/Home"));
const Realtime = lazy(() => import("../views/Realtime/Realtime"));
const History = lazy(() => import("../views/History/History"));
const Export = lazy(() => import("../views/Export/Export"));
const Setting = lazy(() => import("../views/Setting/Setting"));

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
        uri: "/export",
        title: "config.router.export",
        node: <Export />,
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
