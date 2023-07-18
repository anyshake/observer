const getSortedArray = (arr: [], key: string, order = "asc"): [] => {
    if (!arr) {
        return [];
    }

    return arr.sort((a, b) => {
        if (order === "asc") {
            return a[key] - b[key];
        }

        return b[key] - a[key];
    });
};

export default getSortedArray;
