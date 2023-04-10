const AppConfig = {
    frontend: {
        router: "hash",
        version: "v1.0.0d",
        title: `G-Observer 测站面板`,
    },
    backend: {
        host:
            process.env.NODE_ENV === "production"
                ? window.location.hostname
                : `172.17.138.214`,
        port:
            process.env.NODE_ENV === "production"
                ? window.location.port
                : `8073`,
        version: `v1`,
        tls: window.location.protocol === "https:" ? true : false,
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
    },
    sidebar: [
        {
            tag: "index",
            link: "/",
            title: "面板主页",
            icon: (
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 448 512"
                    className="w-4 h-4"
                    fill="currentColor"
                >
                    <path d="M64 32C28.7 32 0 60.7 0 96v64c0 35.3 28.7 64 64 64H448c35.3 0 64-28.7 64-64V96c0-35.3-28.7-64-64-64H64zm280 72a24 24 0 1 1 0 48 24 24 0 1 1 0-48zm48 24a24 24 0 1 1 48 0 24 24 0 1 1 -48 0zM64 288c-35.3 0-64 28.7-64 64v64c0 35.3 28.7 64 64 64H448c35.3 0 64-28.7 64-64V352c0-35.3-28.7-64-64-64H64zm280 72a24 24 0 1 1 0 48 24 24 0 1 1 0-48zm56 24a24 24 0 1 1 48 0 24 24 0 1 1 -48 0z" />
                </svg>
            ),
        },
        {
            tag: "waveform",
            link: "/waveform",
            title: "实时波形",
            icon: (
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 640 512"
                    className="w-4 h-4"
                    fill="currentColor"
                >
                    <path d="M128 64c0-17.7 14.3-32 32-32H320c17.7 0 32 14.3 32 32V416h96V256c0-17.7 14.3-32 32-32H608c17.7 0 32 14.3 32 32s-14.3 32-32 32H512V448c0 17.7-14.3 32-32 32H320c-17.7 0-32-14.3-32-32V96H192V256c0 17.7-14.3 32-32 32H32c-17.7 0-32-14.3-32-32s14.3-32 32-32h96V64z" />
                </svg>
            ),
        },
        {
            tag: "history",
            link: "/history",
            title: "历史数据",
            icon: (
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 448 512"
                    className="w-4 h-4"
                    fill="currentColor"
                >
                    <path d="M96 0C60.7 0 32 28.7 32 64V288H144c6.1 0 11.6 3.4 14.3 8.8L176 332.2l49.7-99.4c2.7-5.4 8.3-8.8 14.3-8.8s11.6 3.4 14.3 8.8L281.9 288H352c8.8 0 16 7.2 16 16s-7.2 16-16 16H272c-6.1 0-11.6-3.4-14.3-8.8L240 275.8l-49.7 99.4c-2.7 5.4-8.3 8.8-14.3 8.8s-11.6-3.4-14.3-8.8L134.1 320H32V448c0 35.3 28.7 64 64 64H352c35.3 0 64-28.7 64-64V160H288c-17.7 0-32-14.3-32-32V0H96zM288 0V128H416L288 0z" />
                </svg>
            ),
        },
    ],
};

export default AppConfig;
