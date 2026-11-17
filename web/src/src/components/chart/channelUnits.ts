import { hasGroundMotion, type IGroundMotion } from './groundMotion';

export const getDisplayUnit = (inputUnit?: string | null, sensitivity?: number | null) =>
    inputUnit === 'Pa' &&
    typeof sensitivity === 'number' &&
    Number.isFinite(sensitivity) &&
    sensitivity > 0
        ? 'Pa'
        : '';

export const isGroundMotionChannel = (
    channel: string,
    inputUnit?: string | null,
    motion?: IGroundMotion
) => {
    if (inputUnit) {
        return inputUnit === 'm/s' || inputUnit === 'm/s2';
    }
    return (
        hasGroundMotion(motion) ||
        (!/^CH\d$/.test(channel) && /^.[HLN][ZNE123UVWRTABC]$/.test(channel))
    );
};
