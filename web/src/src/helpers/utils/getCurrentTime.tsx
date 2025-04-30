export const getCurrentTime = (baseTime: number, refTime: number) => {
    if (!baseTime || !refTime) {
        return Date.now();
    }

    const elapsed = Date.now() - refTime;
    return baseTime + elapsed;
};
