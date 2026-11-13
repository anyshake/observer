// Values arrive in SI units (m/s², m/s, m) and are shown in gal, cm/s and cm,
// which keep typical station readings readable.
export interface IGroundMotion {
    readonly pga?: number | null;
    readonly pgv?: number | null;
    readonly pgd?: number | null;
}

// Peaks of one record, covering [start, end) in milliseconds.
export interface IGroundMotionRecord extends IGroundMotion {
    readonly start: number;
    readonly end: number;
}

export const toCentimetre = (value?: number | null) =>
    value === null || value === undefined || !Number.isFinite(value) ? null : value * 100;

// Fixed four decimals keep the columns from jumping while values update.
export const formatGroundMotion = (value?: number | null) => {
    const centimetre = toCentimetre(value);
    return centimetre === null ? '-' : centimetre.toFixed(4);
};

export const hasGroundMotion = (motion?: IGroundMotion | null) =>
    !!motion && [motion.pga, motion.pgv, motion.pgd].some((value) => Number.isFinite(value));

// Returns the record covering time, using binary search over sorted record starts.
export const findGroundMotion = (
    records: readonly IGroundMotionRecord[] | undefined,
    time: number
) => {
    if (!records?.length) {
        return null;
    }
    let low = 0;
    let high = records.length - 1;
    while (low < high) {
        const middle = (low + high + 1) >> 1;
        if (records[middle].start <= time) {
            low = middle;
        } else {
            high = middle - 1;
        }
    }
    const record = records[low];
    return record.start <= time && time < record.end ? record : null;
};
