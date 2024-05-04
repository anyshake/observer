export const getDayOfYear = (date: Date) => {
	const start = new Date(date.getUTCFullYear(), 0, 0);
	const diff = date.getTime() - start.getTime();
	const oneDay = 1000 * 60 * 60 * 24;
	return Math.floor(diff / oneDay);
};
