const getSortedArray = (
    arr: [],
    key: string,
    keyType: "number" | "datetime",
    sortOrder: "asc" | "desc"
) => {
    if (!arr.length) {
        return [];
    }

    const compare = (a: never, b: never) => {
        if (sortOrder === "desc") {
            return keyType === "datetime"
                ? new Date(b[key]).getTime() - new Date(a[key]).getTime()
                : b[key] - a[key];
        } else {
            return keyType === "datetime"
                ? new Date(a[key]).getTime() - new Date(b[key]).getTime()
                : a[key] - b[key];
        }
    };

    const n = arr.length;
    for (let i = 0; i < n - 1; i++) {
        for (let j = 0; j < n - i - 1; j++) {
            if (compare(arr[j], arr[j + 1]) > 0) {
                const temp = arr[j];
                arr[j] = arr[j + 1];
                arr[j + 1] = temp;
            }
        }
    }

    return arr;
};

export default getSortedArray;
