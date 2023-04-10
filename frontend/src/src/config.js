const AppConfig = {
    backend: {
        host: `172.17.138.214`,
        port: `8073`,
        version: `v1`,
        tls: false,
    },
    api: {
        station: {
            uri: `station`,
            type: `http`,
            method: `get`,
        },
        history: {
            uri: `history`,
            type: `http`,
            method: `get`,
        },
        socket: {
            uri: `socket`,
            type: `websocket`,
            method: `arraybuffer`,
        },
    },
};

export default AppConfig;
