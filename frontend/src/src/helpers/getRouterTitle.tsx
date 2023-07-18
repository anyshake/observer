import ROUTER_CONFIG from "../config/router";
import getRouterUri from "./getRouterUri";

const getRouterTitle = (): string => {
    const uri = getRouterUri();
    for (let i of ROUTER_CONFIG) {
        if (i.uri === uri) {
            return i.title;
        }
    }

    return "发生错误";
};

export default getRouterTitle;
