export const useUrlParams = <T = unknown,>(searchParams: URLSearchParams) => {
    const result: Record<string, string | string[] | undefined> = {};

    for (const [key, value] of searchParams.entries()) {
        if (searchParams.getAll(key).length > 1) {
            result[key] = searchParams.getAll(key);
        } else {
            result[key] = value;
        }
    }

    return result as T;
};
