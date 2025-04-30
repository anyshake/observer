export const getGraphQlApiUrl = () => {
    const baseHost = import.meta.env.VITE_APP_BACKEND_BASE_HOST;
    const apiPath = import.meta.env.VITE_APP_GRAPHQL_API_ENDPOINT;
    const isProduction = import.meta.env.MODE === 'production';

    const protocol = `${window.location.protocol}//`;

    if (isProduction) {
        return baseHost ? `${baseHost}${apiPath}` : `${protocol}${window.location.host}${apiPath}`;
    }

    return `${baseHost}${apiPath}`;
};
