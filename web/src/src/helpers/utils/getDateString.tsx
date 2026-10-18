import { format } from 'date-fns';

export const getDateString = (ts: number) => {
    return format(ts, 'yyyy-MM-dd');
};
