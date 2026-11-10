export const getMapGlyphsUrl = () => {
    const baseHost = import.meta.env.VITE_APP_BACKEND_BASE_HOST;
    const apiEndpoint = import.meta.env.VITE_APP_MAPTILES_API_ENDPOINT;
    const apiPath = `${apiEndpoint}?type=glyph&fontstack={fontstack}&range={range}`;
    const isProduction = import.meta.env.MODE === 'production';

    const protocol = `${window.location.protocol}//`;

    if (isProduction) {
        return baseHost ? `${baseHost}${apiPath}` : `${protocol}${window.location.host}${apiPath}`;
    }

    return `${baseHost}${apiPath}`;
};
