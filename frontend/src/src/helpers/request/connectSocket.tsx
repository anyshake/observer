import { Endpoint } from "../../config/api";
import { getProtocol } from "../app/getProtocol";

export interface SocketOptions<APIRequest, APICommonResponse> {
	readonly backend: string;
	readonly endpoint: Endpoint<APIRequest, APICommonResponse>;
	readonly onOpen?: (event: Event) => void;
	readonly onError?: (event: Event) => void;
	readonly onClose?: (event: CloseEvent) => void;
	readonly onData: (event: MessageEvent<APICommonResponse>) => void;
}

export const connectSocket = <APIRequest, APICommonResponse>({
	backend,
	endpoint,
	onOpen,
	onData,
	onClose,
	onError
}: SocketOptions<APIRequest, APICommonResponse>) => {
	try {
		if (endpoint.type !== "socket") {
			throw new Error("non-websocket protocol is not supported");
		}

		const protocol = getProtocol(false);
		const reqPath = `${protocol}//${backend}${endpoint.path}`;
		const websocket = new WebSocket(reqPath);

		websocket.onmessage = (event) => {
			const data = JSON.parse(event.data);
			onData({ ...event, data } as MessageEvent<APICommonResponse>);
		};
		websocket.onopen = onOpen ?? null;
		websocket.onclose = onClose ?? null;
		websocket.onerror = onError ?? null;

		return websocket;
	} catch {
		return null;
	}
};
