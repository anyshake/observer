export const setLocalStorage = <T,>(key: string, value: T) => {
    localStorage.setItem(
        key,
        typeof value === 'object' ? JSON.stringify(value) : (value as unknown as string)
    );
};
