const API_CONFIG: ApiConfig[] = [
    {
        tag: "station",
        type: "http",
        method: "get",
        uri: "/station",
    },
    {
        tag: "history",
        type: "http",
        method: "post",
        uri: "/history",
    },
    {
        tag: "trace",
        type: "http",
        method: "post",
        uri: "/trace",
    },
    {
        tag: "mseed",
        type: "http",
        method: "post",
        uri: "/mseed",
    },
    {
        tag: "socket",
        type: "ws",
        uri: "/socket",
    },
];

export interface ApiConfig {
    readonly tag: string;
    readonly uri: string;
    readonly type: "http" | "ws";
    readonly method?: "get" | "post";
}

export default API_CONFIG;
