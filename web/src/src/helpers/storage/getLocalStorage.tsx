export const getLocalStorage = <T,>(key: string, fallback: T): T => {
    try {
        const storedItem = localStorage.getItem(key);
        if (storedItem !== null) {
            return typeof fallback === 'object'
                ? (JSON.parse(storedItem) as T)
                : (storedItem as unknown as T);
        }
    } catch {
        /* empty */
    }

    return fallback;
};
