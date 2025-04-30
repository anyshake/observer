export const getCtaMessageDataUrl = (resource: string, locale: string) => {
    const apiPath = import.meta.env.VITE_APP_CTA_MESSAGE_DATA_BASE_URL;
    return `${apiPath}/${resource}/${locale}.json`;
};
