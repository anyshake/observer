const setLocalStorage = (key: string, value: any, json: boolean): void => {
    localStorage.setItem(key, json ? JSON.stringify(value) : value);
};

export default setLocalStorage;
