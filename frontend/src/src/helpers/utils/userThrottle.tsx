/* eslint-disable @typescript-eslint/no-explicit-any */
export const userThrottle = (fn: (...args: any[]) => void, wait: number) => {
	let timerId: number | null = null;
	return (...args: any[]) => {
		if (timerId !== null) {
			return;
		}
		timerId = window.setTimeout(() => {
			timerId = null;
			fn(...args);
		}, wait);
	};
};
