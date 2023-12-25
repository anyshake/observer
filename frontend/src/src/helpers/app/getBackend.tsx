const API_FOR_DEVELOPMENT = `127.0.0.1:8073`;

const getBackend = (): string => {
    const production = process.env.NODE_ENV === "production";
    return production
        ? `//${window.location.host}`
        : `//${API_FOR_DEVELOPMENT}`;
};

export default getBackend;
