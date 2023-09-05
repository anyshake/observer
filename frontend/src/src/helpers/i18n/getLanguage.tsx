import I18N_CONFIG, { I18nConfigItem } from "../../config/i18n";
import getLocalStorage from "../storage/getLocalStorage";
import setLocalStorage from "../storage/setLocalStorage";

const getLanguage = (): string => {
    const { list, fallback, key } = I18N_CONFIG;
    const currentLang = getLocalStorage(key, "unknown", false);

    // Use fallback if current language is not set
    if (currentLang === "unknown") {
        setLocalStorage(key, fallback, false);
        return fallback;
    }
    // Match current language with available languages
    const matchedLanguage = list.find((lang: I18nConfigItem) => {
        return lang.value === currentLang;
    });
    // Return fallback if language is not available
    return matchedLanguage ? matchedLanguage.value : fallback;
};

export default getLanguage;
