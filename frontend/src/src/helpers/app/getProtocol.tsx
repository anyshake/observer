export const getProtocol = (http: boolean) => {
    if (process.env.NODE_ENV === "production") {
        if (http) {
            return window.location.protocol;
        }
        return window.location.protocol === "https:" ? "wss:" : "ws:";
    }
    if (http) {
        return !!process.env.REACT_APP_BACKEND_TLS ? "https:" : "http:";
    }
    return !!process.env.REACT_APP_BACKEND_TLS ? "wss:" : "ws:";
};
