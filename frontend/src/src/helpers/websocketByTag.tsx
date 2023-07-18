import API_CONFIG from "../config/api";
import getApiUrl from "./getApiUrl";
import getBackend from "./getBackend";

interface WebsocketByTag {
    readonly tag: string;
    readonly onOpen?: (event: Event) => void;
    readonly onData: (event: MessageEvent) => void;
    readonly onClose?: (event: CloseEvent) => void;
    readonly onError?: (event: Event) => void;
}

const websocketByTag = ({
    tag,
    onOpen,
    onData,
    onClose,
    onError,
}: WebsocketByTag): WebSocket | void => {
    try {
        const type = API_CONFIG.find((config) => config.tag === tag)?.type;
        if (type !== "ws") {
            throw new Error("non-websocket protocol is not supported");
        }

        const uri = getApiUrl(tag);
        const backend = `ws:${getBackend()}`;
        const ws = new WebSocket(`${backend}${uri}`);

        ws.onmessage = onData;
        ws.onopen = onOpen || (() => {});
        ws.onclose = onClose || (() => {});
        ws.onerror = onError || (() => {});

        return ws;
    } catch (err: unknown) {
        if (onError) {
            onError(new Event("error occurred"));
        }

        return;
    }
};

export default websocketByTag;
