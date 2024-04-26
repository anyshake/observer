export const getRelease = () => {
    return process.env.REACT_APP_RELEASE ?? "unknown";
};
