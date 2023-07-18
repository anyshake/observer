const getBackend = (): string => {
    const production = process.env.NODE_ENV === "production";
    return production
        ? `//${window.location.host}`
        : `//127.0.0.1:8073`;
};

export default getBackend;
