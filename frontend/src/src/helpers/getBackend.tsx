const getBackend = (): string => {
    const production = process.env.NODE_ENV === "production";
    return production
        ? `//${window.location.host}`
        // : `//wolfx.p.sdrotg.com`;
        : `//172.17.138.214:8073`;
};

export default getBackend;
