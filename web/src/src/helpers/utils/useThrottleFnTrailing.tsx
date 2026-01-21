import { useCallback, useEffect, useRef } from 'react';

type Fn<Args extends readonly unknown[]> = (...args: Args) => void;

export function useThrottleFnTrailing<Args extends readonly unknown[]>(
    fn: Fn<Args>,
    delay: number
): (...args: Args) => void {
    const fnRef = useRef(fn);
    const lastTimeRef = useRef<number>(0);
    const timerRef = useRef<number | null>(null);

    useEffect(() => {
        fnRef.current = fn;
    }, [fn]);

    const throttled = useCallback(
        (...args: Args) => {
            const now = Date.now();
            const remaining = delay - (now - lastTimeRef.current);

            if (remaining <= 0) {
                if (timerRef.current !== null) {
                    clearTimeout(timerRef.current);
                    timerRef.current = null;
                }
                lastTimeRef.current = now;
                fnRef.current(...args);
            } else if (timerRef.current === null) {
                timerRef.current = window.setTimeout(() => {
                    lastTimeRef.current = Date.now();
                    timerRef.current = null;
                    fnRef.current(...args);
                }, remaining);
            }
        },
        [delay]
    );

    return throttled;
}
