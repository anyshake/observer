import { i18n } from "i18next";

export const getCurrentLocale = async (i18n: Promise<i18n>) => {
    const lang = (await i18n).language;
    if (lang === "en") {
        return "en-US";
    }

    return lang;
};
