import { BannerProps } from "../../components/Banner";
import { I18nTranslation } from "../../config/i18n";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";

const setBanner = (res?: ApiResponse): BannerProps => {
    // Parse response, empty response means error
    const { error } = res || {};
    const { station, uptime, os } = res?.data || {};
    const { uuid, name } = station || "Unknown";

    // Error banner by default
    let label = { id: "views.home.banner.error.label" } as I18nTranslation;
    let text = { id: "views.home.banner.error.text" } as I18nTranslation;
    // Change to success banner and return if no error
    if (!error) {
        label = {
            id: "views.home.banner.success.label",
            format: { station: name },
        };
        text = {
            id: "views.home.banner.success.text",
            format: { ...os, uptime, uuid },
        };
    }

    return {
        type: error ? "error" : "success",
        label,
        text,
    };
};

export default setBanner;
