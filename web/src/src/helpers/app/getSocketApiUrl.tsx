export const getSocketApiUrl = (token?: string) => {
    const baseHost = import.meta.env.VITE_APP_BACKEND_BASE_HOST;
    const wsEndpoint = import.meta.env.VITE_APP_WEBSOCKET_API_ENDPOINT;
    const isProduction = import.meta.env.MODE === 'production';

    const getWebSocketProtocol = (url: string) => (url.startsWith('https:') ? 'wss://' : 'ws://');

    if (isProduction) {
        const protocol = getWebSocketProtocol(baseHost ? baseHost : window.location.protocol);
        return baseHost
            ? `${baseHost.replace(/^https?:\/\//, protocol)}${wsEndpoint}${token ? `?token=${token}` : ''}`
            : `${protocol}${window.location.host}${wsEndpoint}${token ? `?token=${token}` : ''}`;
    }

    const parsedBaseHost = new URL(baseHost);
    const protocol = getWebSocketProtocol(parsedBaseHost.protocol);
    return `${baseHost.replace(/^https?:\/\//, protocol)}${wsEndpoint}${token ? `?token=${token}` : ''}`;
};
