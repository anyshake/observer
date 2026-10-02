import { useEffect, useRef } from 'react';

export interface IEventStreamMessage {
    readonly data: string;
    readonly event: string;
    readonly id?: string;
}

interface IEventStreamOptions {
    readonly url: string;
    readonly token: string;
    readonly onEvent: (message: IEventStreamMessage) => void;
    readonly onUnauthorized?: () => void;
}

const initialRetryDelay = 1000;
const maxRetryDelay = 30000;

const parseEvent = (block: string): IEventStreamMessage | null => {
    let event = 'message';
    let id: string | undefined;
    const data: string[] = [];

    block.split('\n').forEach((line) => {
        if (!line.length || line.startsWith(':')) {
            return;
        }

        const separatorIndex = line.indexOf(':');
        const field = separatorIndex === -1 ? line : line.slice(0, separatorIndex);
        let value = separatorIndex === -1 ? '' : line.slice(separatorIndex + 1);
        if (value.startsWith(' ')) {
            value = value.slice(1);
        }

        switch (field) {
            case 'event':
                event = value;
                break;
            case 'id':
                id = value;
                break;
            case 'data':
                data.push(value);
                break;
        }
    });

    if (!data.length) {
        return null;
    }

    return { data: data.join('\n'), event, id };
};

export const useEventStream = ({ url, token, onEvent, onUnauthorized }: IEventStreamOptions) => {
    const onEventRef = useRef(onEvent);
    const onUnauthorizedRef = useRef(onUnauthorized);

    useEffect(() => {
        onEventRef.current = onEvent;
        onUnauthorizedRef.current = onUnauthorized;
    }, [onEvent, onUnauthorized]);

    useEffect(() => {
        if (!token.length) {
            return;
        }

        let stopped = false;
        let retryDelay = initialRetryDelay;
        let retryTimer: ReturnType<typeof setTimeout> | undefined;
        let abortController: AbortController | undefined;

        function scheduleReconnect() {
            if (stopped) {
                return;
            }

            retryTimer = setTimeout(() => {
                void connect();
            }, retryDelay);
            retryDelay = Math.min(retryDelay * 2, maxRetryDelay);
        }

        async function connect() {
            abortController = new AbortController();

            try {
                const response = await fetch(url, {
                    headers: {
                        Accept: 'text/event-stream',
                        Authorization: `Bearer ${token}`
                    },
                    signal: abortController.signal
                });

                if (response.status === 401) {
                    onUnauthorizedRef.current?.();
                    return;
                }
                if (!response.ok || !response.body) {
                    throw new Error(`event stream request failed with status ${response.status}`);
                }

                retryDelay = initialRetryDelay;
                const reader = response.body.getReader();
                const decoder = new TextDecoder();
                let buffer = '';

                while (!stopped) {
                    const { done, value } = await reader.read();
                    if (done) {
                        break;
                    }

                    buffer += decoder.decode(value, { stream: true });
                    buffer = buffer.replace(/\r\n/g, '\n');

                    let boundaryIndex = buffer.indexOf('\n\n');
                    while (boundaryIndex !== -1) {
                        const message = parseEvent(buffer.slice(0, boundaryIndex));
                        buffer = buffer.slice(boundaryIndex + 2);
                        if (message) {
                            onEventRef.current(message);
                        }
                        boundaryIndex = buffer.indexOf('\n\n');
                    }
                }

                scheduleReconnect();
            } catch (error) {
                if (!(error instanceof DOMException && error.name === 'AbortError')) {
                    scheduleReconnect();
                }
            }
        }

        void connect();

        return () => {
            stopped = true;
            abortController?.abort();
            if (retryTimer !== undefined) {
                clearTimeout(retryTimer);
            }
        };
    }, [token, url]);
};
