import { format } from "date-fns";

export const getTimeString = (ts: number) => {
	const date = new Date(ts);
	return format(date, "yyyy-MM-dd HH:mm:ss");
};
