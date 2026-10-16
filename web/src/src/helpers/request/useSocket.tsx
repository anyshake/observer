import { useCallback, useEffect, useRef, useState } from 'react';

import { connectSocket, ISocket } from './connectSocket';

const connectionTimeout = 10000;
const initialRetryDelay = 1000;
const maxRetryDelay = 10000;

export const useSocket = (options: ISocket, reconnect: boolean) => {
    const optionsRef = useRef(options);
    const socketRef = useRef<WebSocket | null>(null);
    const [readyState, setReadyState] = useState<WebSocket['readyState']>(WebSocket.CONNECTING);
    optionsRef.current = options;

    const { url } = options;
    const protocols = JSON.stringify(options.protocols ?? []);
    const heartbeatMessage = options.heartbeat?.message;
    const heartbeatResponse = options.heartbeat?.response;
    const heartbeatInterval = options.heartbeat?.interval ?? 10000;
    const heartbeatTimeout = options.heartbeat?.timeout ?? 5000;

    const sendMessage = useCallback((message: unknown, json?: boolean) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(json ? JSON.stringify(message) : String(message));
        }
    }, []);

    useEffect(() => {
        let stopped = false;
        let retryDelay = initialRetryDelay;
        let retryTimer: ReturnType<typeof setTimeout> | undefined;
        let connectTimer: ReturnType<typeof setTimeout> | undefined;
        let heartbeatTimer: ReturnType<typeof setTimeout> | undefined;

        const isCurrent = (socket: WebSocket) => !stopped && socketRef.current === socket;

        function disconnect() {
            clearTimeout(retryTimer);
            clearTimeout(connectTimer);
            clearTimeout(heartbeatTimer);
            retryTimer = undefined;
            connectTimer = undefined;
            heartbeatTimer = undefined;

            const socket = socketRef.current;
            socketRef.current = null;
            if (socket) {
                // Closing an old socket must not change the next connection's state.
                socket.onopen = null;
                socket.onclose = null;
                socket.onerror = null;
                socket.onmessage = null;
                socket.close();
            }
        }

        function scheduleReconnect() {
            disconnect();
            if (stopped) {
                return;
            }

            setReadyState(WebSocket.CLOSED);
            if (!reconnect || !navigator.onLine) {
                return;
            }

            retryTimer = setTimeout(() => {
                retryTimer = undefined;
                connect();
            }, retryDelay);
            retryDelay = Math.min(retryDelay * 2, maxRetryDelay);
        }

        function scheduleHeartbeat(socket: WebSocket, delay = heartbeatInterval) {
            clearTimeout(heartbeatTimer);
            if (heartbeatMessage === undefined || heartbeatResponse === undefined) {
                return;
            }

            heartbeatTimer = setTimeout(() => {
                if (!isCurrent(socket) || socket.readyState !== WebSocket.OPEN) {
                    return;
                }

                heartbeatTimer = setTimeout(() => {
                    if (isCurrent(socket)) {
                        scheduleReconnect();
                    }
                }, heartbeatTimeout);
                try {
                    socket.send(heartbeatMessage);
                } catch {
                    scheduleReconnect();
                }
            }, delay);
        }

        function connect() {
            if (stopped || socketRef.current) {
                return;
            }

            setReadyState(WebSocket.CONNECTING);
            try {
                const socket = connectSocket({
                    url,
                    protocols: JSON.parse(protocols),
                    onClose: (event) => {
                        if (!isCurrent(socket)) {
                            return;
                        }
                        scheduleReconnect();
                        optionsRef.current.onClose?.(event);
                    },
                    onError: (event) => {
                        if (!isCurrent(socket)) {
                            return;
                        }
                        scheduleReconnect();
                        optionsRef.current.onError?.(event);
                    },
                    onOpen: (event) => {
                        if (!isCurrent(socket)) {
                            return;
                        }
                        clearTimeout(connectTimer);
                        connectTimer = undefined;
                        setReadyState(WebSocket.OPEN);
                        scheduleHeartbeat(socket);
                        optionsRef.current.onOpen?.(event);
                    },
                    onData: (event) => {
                        if (!isCurrent(socket)) {
                            return;
                        }
                        retryDelay = initialRetryDelay;
                        scheduleHeartbeat(socket);
                        if (heartbeatResponse !== undefined && event.data === heartbeatResponse) {
                            return;
                        }
                        optionsRef.current.onData(event);
                    }
                });
                socketRef.current = socket;
                connectTimer = setTimeout(() => {
                    if (isCurrent(socket)) {
                        scheduleReconnect();
                    }
                }, connectionTimeout);
            } catch {
                scheduleReconnect();
                optionsRef.current.onError?.(new Event('error'));
            }
        }

        function handleOnline() {
            if (!reconnect || stopped) {
                return;
            }
            retryDelay = initialRetryDelay;
            disconnect();
            connect();
        }

        function handleOffline() {
            if (reconnect) {
                disconnect();
                setReadyState(WebSocket.CLOSED);
            }
        }

        function handleVisibilityChange() {
            if (document.visibilityState !== 'visible' || !navigator.onLine || !reconnect) {
                return;
            }
            const socket = socketRef.current;
            if (socket?.readyState === WebSocket.OPEN) {
                scheduleHeartbeat(socket, 0);
            } else if (!socket) {
                handleOnline();
            }
        }

        window.addEventListener('online', handleOnline);
        window.addEventListener('offline', handleOffline);
        document.addEventListener('visibilitychange', handleVisibilityChange);
        connect();

        return () => {
            stopped = true;
            window.removeEventListener('online', handleOnline);
            window.removeEventListener('offline', handleOffline);
            document.removeEventListener('visibilitychange', handleVisibilityChange);
            disconnect();
        };
    }, [
        url,
        protocols,
        reconnect,
        heartbeatMessage,
        heartbeatResponse,
        heartbeatInterval,
        heartbeatTimeout
    ]);

    return { sendMessage, readyState };
};
