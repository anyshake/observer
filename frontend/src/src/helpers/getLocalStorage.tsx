const getLocalStorage = (key: string, fallback: string): string => {
    if (localStorage.getItem(key)) {
        return localStorage.getItem(key) as string;
    }

    localStorage.setItem(key, fallback);
    return fallback;
};

export default getLocalStorage;
