import i18n, { i18n as I18n } from "i18next";
import BackendAdapter from "i18next-http-backend";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

const createI18n = async (
    loadPath: string,
    fallbackLng: string,
    storageKey: string
): Promise<I18n> => {
    const detector = new LanguageDetector(null, {
        lookupLocalStorage: storageKey,
    });
    await i18n
        .use(initReactI18next)
        .use(BackendAdapter)
        .use(detector)
        .init({
            fallbackLng,
            detection: {
                caches: ["localStorage"],
                order: ["localStorage", "navigator"],
            },
            interpolation: {
                escapeValue: true,
            },
            backend: { loadPath },
        });

    return i18n;
};

export default createI18n;
