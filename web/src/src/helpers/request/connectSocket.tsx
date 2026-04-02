export interface ISocket {
    readonly url: string;
    readonly protocols?: string[];
    readonly onOpen?: (event: Event) => void;
    readonly onError?: (event: Event) => void;
    readonly onClose?: (event: CloseEvent) => void;
    readonly onData: (data: MessageEvent) => void;
}

export const connectSocket = ({ url, protocols, onOpen, onData, onClose, onError }: ISocket) => {
    const websocket = new WebSocket(url, protocols);

    websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        onData({ ...event, data } as MessageEvent);
    };
    websocket.onopen = onOpen ?? null;
    websocket.onclose = onClose ?? null;
    websocket.onerror = onError ?? null;

    return websocket;
};
