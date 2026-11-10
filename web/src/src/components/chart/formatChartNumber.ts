export const formatChartNumber = (value: number | null | undefined) => {
    if (value === null || value === undefined || !Number.isFinite(value)) {
        return '-';
    }
    const abs = Math.abs(value);
    if (abs !== 0 && (abs < 1e-3 || abs >= 1e5)) {
        return value.toExponential(2).replace(/\.?0+e/, 'e');
    }
    return String(Number(value.toPrecision(4)));
};
