import { useCallback, useEffect, useRef, useState } from 'react';

import { connectSocket, ISocket } from './connectSocket';

export const useSocket = (options: ISocket, reconnect: boolean) => {
    const optionsRef = useRef<ISocket | null>(options);
    const socketRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    const [readyState, setReadyState] = useState<WebSocket['readyState']>(WebSocket.CONNECTING);

    const sendMessage = useCallback((message: unknown, json?: boolean) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(json ? JSON.stringify(message) : String(message));
        }
    }, []);

    useEffect(() => {
        const connect = () => {
            if (!optionsRef.current || socketRef.current) {
                return;
            }

            socketRef.current = connectSocket({
                ...optionsRef.current,
                onClose: (e) => {
                    optionsRef.current?.onClose?.(e);
                    setReadyState(WebSocket.CLOSED); // Update readyState on close
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
                    setReadyState(WebSocket.CLOSED); // Update readyState on error
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
                onOpen: () => {
                    setReadyState(WebSocket.OPEN); // Update readyState on open
                },
                onData: (data) => {
                    optionsRef.current?.onData?.(data);
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

    return { sendMessage, readyState };
};
