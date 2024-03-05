export const getAccelerationArr = (
    data: number[],
    intervalMS: number
): number[] => {
    const acceleration = [];
    for (let i = 0; i < data.length; i++) {
        if (!i) {
            acceleration.push(0);
            continue;
        }

        acceleration.push((data[i] - data[i - 1]) / (intervalMS / 1000));
    }

    return acceleration;
};
