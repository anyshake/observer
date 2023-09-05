import ROUTER_CONFIG from "../../config/router";
import getRouterUri from "./getRouterUri";

const getRouterTitle = (): string => {
    const uri = getRouterUri();
    for (let i of ROUTER_CONFIG) {
        if (i.uri === uri) {
            return i.title;
        }
    }

    return "config.router.error";
};

export default getRouterTitle;
