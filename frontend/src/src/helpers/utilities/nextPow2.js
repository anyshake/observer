const nextPow2 = (arr) => {
    let p = 1;
    while (p < arr.length) {
        p *= 2;
    }

    return [...arr, ...Array(p - arr.length).fill(0)];
};

export default nextPow2;
