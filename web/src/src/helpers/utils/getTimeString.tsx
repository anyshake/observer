import { format } from 'date-fns';

export const getTimeString = (ts: number) => {
    return format(ts, 'yyyy-MM-dd HH:mm:ss');
};
