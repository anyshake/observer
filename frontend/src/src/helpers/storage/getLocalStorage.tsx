const getLocalStorage = (key: string, fallback: any, json: boolean): any => {
    try {
        const storedItem = localStorage.getItem(key);
        if (storedItem !== null) {
            if (!json) {
                return storedItem;
            }

            return JSON.parse(storedItem);
        }
    } catch (err) {
        localStorage.setItem(key, json ? JSON.stringify(fallback) : fallback);
    }

    return fallback;
};

export default getLocalStorage;
