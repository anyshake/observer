import { i18n } from "i18next";
import getRouterTitle from "../router/getRouterTitle";
import GLOBAL_CONFIG from "../../config/global";

const toggleI18n = async (i18n: Promise<i18n>, lang: string): Promise<void> => {
    const instance = await i18n;
    await instance.changeLanguage(lang);

    // Update document title
    const subtitle = getRouterTitle();
    const { title } = GLOBAL_CONFIG.app_settings;
    document.title = `${instance.t(subtitle)} | ${instance.t(title)}`;
};

export default toggleI18n;
