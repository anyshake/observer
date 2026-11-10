import { RetentionConstraints } from '../../config/constraints';

export const normalizeRetention = (
    retention: unknown,
    fallback: number = RetentionConstraints.defaultValue
) => {
    const { minValue, maxValue, step } = RetentionConstraints;
    const value =
        typeof retention === 'number' || (typeof retention === 'string' && retention.trim())
            ? Number(retention)
            : NaN;
    if (!Number.isFinite(value)) {
        return fallback;
    }
    return Math.min(maxValue, Math.max(minValue, Math.round(value / step) * step));
};
