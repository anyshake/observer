export interface ISocket {
    readonly url: string;
    readonly onOpen?: (event: Event) => void;
    readonly onError?: (event: Event) => void;
    readonly onClose?: (event: CloseEvent) => void;
    readonly onData: (data: MessageEvent) => void;
}

export const connectSocket = ({ url, onOpen, onData, onClose, onError }: ISocket) => {
    try {
        const websocket = new WebSocket(url);

        websocket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            onData({ ...event, data } as MessageEvent);
        };
        websocket.onopen = onOpen ?? null;
        websocket.onclose = onClose ?? null;
        websocket.onerror = onError ?? null;

        return websocket;
    } catch {
        return null;
    }
};
