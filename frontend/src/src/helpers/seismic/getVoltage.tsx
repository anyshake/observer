const getVoltage = (
    arr: number[],
    resolution: number,
    fullscale: number
): number[] => {
    const factor = 2 ** (resolution - 1);
    const voltage: number[] = [];
    for (let i of arr) {
        voltage.push((fullscale / factor) * i);
    }

    return voltage;
};

export default getVoltage;
