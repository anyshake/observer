const getCounts = (data: number[]): number[] => {
    const sortedData = data.slice().sort((a, b) => a - b);

    let median: number;
    if (sortedData.length % 2) {
        const middleIndex = Math.floor(sortedData.length / 2);
        median = sortedData[middleIndex];
    } else {
        const middleIndex1 = sortedData.length / 2 - 1;
        const middleIndex2 = sortedData.length / 2;
        median = (sortedData[middleIndex1] + sortedData[middleIndex2]) / 2;
    }

    return data.map((value) => value - median);
};

export default getCounts;
