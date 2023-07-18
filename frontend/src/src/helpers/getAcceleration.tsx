const getAcceleration = (data: number[], interval: number): number[] => {
    const acceleration = [];
    for (let i = 0; i < data.length; i++) {
        if (!i) {
            acceleration.push(0);
            continue;
        }

        acceleration.push((data[i] - data[i - 1]) / interval);
    }

    return acceleration;
};

export default getAcceleration;
