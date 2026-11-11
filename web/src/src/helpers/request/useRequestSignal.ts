import { useCallback, useEffect, useRef } from 'react';

// Lazy queries and imperative REST requests can outlive their subscribers.
export const useRequestSignal = () => {
    const controller = useRef<AbortController | null>(null);
    const getRequestSignal = useCallback(() => {
        controller.current ??= new AbortController();
        return controller.current.signal;
    }, []);

    useEffect(() => {
        if (controller.current?.signal.aborted) {
            controller.current = new AbortController();
        }
        return () => controller.current?.abort();
    }, []);

    return getRequestSignal;
};
