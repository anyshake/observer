import GLOBAL_CONFIG from "../../config/global";

export interface RouterParam<T> {
    [key: string]: T;
}

const getRouterParam = (): RouterParam<string | number | boolean> => {
    const { hash, search } = window.location;
    const { router } = GLOBAL_CONFIG.app_settings;

    const param = router === "hash" ? hash.split("?")[1] : search.split("?")[1];
    if (param) {
        const paramList = param.split("&");
        const paramObject: RouterParam<string> = {};
        for (let i of paramList) {
            const [key, value] = i.split("=");
            paramObject[key] = value;
        }

        return paramObject;
    }

    return {};
};

export default getRouterParam;
