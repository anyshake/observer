export const getVelocityArr = (data: number[], sensitivity: number): number[] => {
    const velocity: number[] = [];
    for (let i of data) {
        velocity.push(i / sensitivity);
    }

    return velocity;
};
