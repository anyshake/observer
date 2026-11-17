export type HealthLevel = 'ok' | 'warning' | 'error' | 'unknown';

export interface IHealth {
    readonly level: HealthLevel;
    readonly value: string;
    readonly detail: string;
}

export const healthStatusClass: Record<HealthLevel, string> = {
    ok: 'status-success',
    warning: 'status-warning',
    error: 'status-error',
    unknown: 'status-neutral'
};
