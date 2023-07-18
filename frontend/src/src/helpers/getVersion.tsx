const getVersion = (): string => {
    return process.env.REACT_APP_VERSION || "unknown";
};

export default getVersion;
