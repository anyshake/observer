import { lazy } from "react";

const Routes = [
    {
        index: true,
        timer: 1000,
        component: "./app/v1/stationInfo",
        path: "/",
    },
    {
        timer: 1000,
        component: "./app/v1/realtimeWaveform",
        path: "/waveform",
    },
    {
        timer: 1000,
        component: "./app/v1/historyWaveform",
        path: "/history",
    },
];

const routerConfig = new Array(Routes.length);
Routes.forEach((item, index) => {
    routerConfig[index] = lazy(async () => {
        const [moduleExports] = await Promise.all([
            import(`${item.component}`),
            new Promise((resolve) => setTimeout(resolve, item.timer)),
        ]);
        return moduleExports;
    });
});

export { Routes };
export default routerConfig;
