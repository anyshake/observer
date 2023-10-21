const API_FOR_DEVELOPMENT = `wolfx.p.sdrotg.com`;

const getBackend = (): string => {
    const production = process.env.NODE_ENV === "production";
    return production
        ? `//${window.location.host}`
        : `//${API_FOR_DEVELOPMENT}`;
};

export default getBackend;
