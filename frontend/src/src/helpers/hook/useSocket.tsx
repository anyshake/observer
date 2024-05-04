import { useEffect, useRef } from "react";

import { connectSocket, SocketOptions } from "../request/connectSocket";

export const useSocket = <APIRequest, APICommonResponse>(
	options: SocketOptions<APIRequest, APICommonResponse>,
	reconnect = true
) => {
	const optionsRef = useRef<SocketOptions<APIRequest, APICommonResponse> | null>(options);
	const socketRef = useRef<WebSocket | null>(null);
	const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

	useEffect(() => {
		const connect = () => {
			if (!optionsRef.current || socketRef.current) {
				return;
			}

			socketRef.current = connectSocket<APIRequest, APICommonResponse>({
				...optionsRef.current,
				onClose: (e) => {
					optionsRef.current?.onClose?.(e);
					if (reconnect) {
						if (reconnectTimeoutRef.current) {
							clearTimeout(reconnectTimeoutRef.current);
						}
						reconnectTimeoutRef.current = setTimeout(() => {
							socketRef.current = null;
							connect();
						}, 1000);
					} else {
						socketRef.current = null;
					}
				},
				onError: (e) => {
					optionsRef.current?.onError?.(e);
					if (reconnect) {
						if (reconnectTimeoutRef.current) {
							clearTimeout(reconnectTimeoutRef.current);
						}
						reconnectTimeoutRef.current = setTimeout(() => {
							socketRef.current = null;
							connect();
						}, 1000);
					} else {
						socketRef.current = null;
					}
				}
			});
		};
		connect();

		return () => {
			if (reconnectTimeoutRef.current) {
				clearTimeout(reconnectTimeoutRef.current);
			}
			socketRef.current?.close();
			optionsRef.current = null;
			socketRef.current = null;
		};
	}, [reconnect]);
};
