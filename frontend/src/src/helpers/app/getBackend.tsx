const API_FOR_DEVELOPMENT = `172.17.138.214:8073`;

const getBackend = (): string => {
    const production = process.env.NODE_ENV === "production";
    return production
        ? `//${window.location.host}`
        : `//${API_FOR_DEVELOPMENT}`;
};

export default getBackend;
