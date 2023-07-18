const getRelease = (): string => {
    return process.env.REACT_APP_RELEASE || "unknown";
};

export default getRelease;
