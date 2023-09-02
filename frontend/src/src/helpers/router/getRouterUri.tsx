import GLOBAL_CONFIG from "../../config/global";

const getRouterUri = () => {
    const { hash, pathname } = window.location;
    const { router } = GLOBAL_CONFIG.app_settings;

    let uri = pathname;
    if (router === "hash" && hash.length > 0) {
        uri = hash.replaceAll("#", "");
    }

    return uri;
};

export default getRouterUri;
