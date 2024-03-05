import { useEffect, useRef } from "react";
import { SocketOptions, connectSocket } from "../request/connectSocket";

export const useSocket = <APIRequest, APICommonResponse>(
    options: SocketOptions<APIRequest, APICommonResponse>,
    reconnect = true
) => {
    const optionsRef = useRef<SocketOptions<
        APIRequest,
        APICommonResponse
    > | null>(options);
    const socketRef = useRef<WebSocket | null>();

    useEffect(() => {
        const connect = () => {
            if (!optionsRef.current) {
                return;
            }
            socketRef.current = connectSocket<APIRequest, APICommonResponse>({
                ...optionsRef.current,
                onClose: (e) => {
                    socketRef.current = null;
                    optionsRef.current?.onClose?.(e);
                    reconnect && setTimeout(() => connect(), 1000);
                },
                onError: (e) => {
                    socketRef.current = null;
                    optionsRef.current?.onError?.(e);
                    reconnect && setTimeout(() => connect(), 1000);
                },
            });
        };
        connect();

        return () => {
            socketRef.current?.close();
            optionsRef.current = null;
        };
    }, [reconnect]);
};
