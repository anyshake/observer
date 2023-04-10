const getApiUrl = ({ host, port, api, version, tls, type }) => {
    const baseUrl = `${host}:${port}/api/${version}/${api}`;
    switch (type) {
        case `http`:
            if (tls) {
                return `https://${baseUrl}`;
            }
            return `http://${baseUrl}`;

        case `websocket`:
            if (tls) {
                return `wss://${baseUrl}`;
            }
            return `ws://${baseUrl}`;

        default:
            return null;
    }
};

export default getApiUrl;
