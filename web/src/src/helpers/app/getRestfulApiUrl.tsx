export const getRestfulApiUrl = (resource: string) => {
    const baseHost = import.meta.env.VITE_APP_BACKEND_BASE_HOST;
    const apiBasePath = import.meta.env.VITE_APP_RESTFUL_API_BASE_PATH;
    const isProduction = import.meta.env.MODE === 'production';

    const protocol = `${window.location.protocol}//`;

    if (isProduction) {
        return baseHost
            ? `${apiBasePath}${resource}`
            : `${protocol}${window.location.host}${apiBasePath}${resource}`;
    }

    return `${baseHost}${apiBasePath}${resource}`;
};
