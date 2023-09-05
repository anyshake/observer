import I18N_CONFIG from "../../config/i18n";
import getLocalStorage from "../storage/getLocalStorage";

const getLanguage = (): string => {
    const { key, fallback } = I18N_CONFIG;
    return getLocalStorage(key, fallback, false);
};

export default getLanguage;
