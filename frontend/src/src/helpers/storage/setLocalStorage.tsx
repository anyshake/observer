const setLocalStorage = <T,>(key: string, value: T, json: boolean): void => {
    localStorage.setItem(key, json ? JSON.stringify(value) : (value as string));
};

export default setLocalStorage;
