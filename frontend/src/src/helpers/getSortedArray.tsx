const getSortedArray = (arr: [], key: string, sortOrder: "asc" | "desc") => {
    if (!arr.length) {
        return [];
    }

    const compare = (a: never, b: never) => {
        if (sortOrder === "desc") {
            return b[key] - a[key];
        } else {
            return a[key] - b[key];
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
