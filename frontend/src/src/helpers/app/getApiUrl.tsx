import API_CONFIG from "../../config/api";
import GLOBAL_CONFIG from "../../config/global";

const getApiUrl = (tag: string): string => {
    const { version, prefix } = GLOBAL_CONFIG.api_settings;
    const uri = API_CONFIG.find((config) => config.tag === tag)?.uri as string;

    return `${prefix}/${version}${uri}`;
};

export default getApiUrl;
