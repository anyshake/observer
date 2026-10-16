export type CleanupMode = 'all' | 'range';
export type CleanupTargetId = 'seisRecords' | 'miniSeed' | 'helicorder';

export interface CleanupRange {
    startDate?: number;
    endDate?: number;
}

export const initialCleanupRanges: Record<CleanupTargetId, CleanupRange> = {
    seisRecords: {},
    miniSeed: {},
    helicorder: {}
};

export const initialCleanupModes: Record<CleanupTargetId, CleanupMode> = {
    seisRecords: 'all',
    miniSeed: 'all',
    helicorder: 'all'
};

export const getErrorMessage = (error: unknown) => {
    if (error instanceof Error) {
        return error.message;
    }

    return String(error);
};

export const cleanupTargetByKind = (kind?: string): CleanupTargetId | undefined => {
    if (kind?.startsWith('seis_records')) {
        return 'seisRecords';
    }
    if (kind?.startsWith('miniseed')) {
        return 'miniSeed';
    }
    if (kind?.startsWith('helicorder')) {
        return 'helicorder';
    }
    return undefined;
};
