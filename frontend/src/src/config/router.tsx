import { ReactNode, lazy } from "react";

const Home = lazy(() => import("../views/Home/Index"));
const Realtime = lazy(() => import("../views/Realtime/Index"));
const History = lazy(() => import("../views/History/Index"));
const Setting = lazy(() => import("../views/Setting/Index"));

const ROUTER_CONFIG: RouterConfig[] = [
    {
        uri: "/",
        title: "测站状态",
        node: <Home />,
    },
    {
        uri: "/realtime",
        title: "实时波形",
        node: <Realtime />,
    },
    {
        uri: "/history",
        title: "历史查询",
        node: <History />,
    },
    {
        uri: "/setting",
        title: "面板设定",
        node: <Setting />,
    },
];

export interface RouterConfig {
    readonly uri: string;
    readonly title: string;
    readonly node: ReactNode;
}

export default ROUTER_CONFIG;
