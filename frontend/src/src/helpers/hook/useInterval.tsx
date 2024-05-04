import { useEffect, useRef } from "react";

export const useInterval = (
	callback: () => void,
	delay: number | null | false,
	immediate?: boolean
) => {
	// eslint-disable-next-line @typescript-eslint/no-empty-function
	const savedCallback = useRef(() => {});

	useEffect(() => {
		savedCallback.current = callback;
	});

	useEffect(() => {
		if (!immediate || delay === null || delay === false) {
			return;
		}
		savedCallback.current();
	}, [immediate, delay]);

	useEffect(() => {
		if (delay === null || delay === false) {
			return;
		}
		const tick = () => {
			savedCallback.current();
		};
		const id = setInterval(tick, delay);
		return () => {
			clearInterval(id);
		};
	}, [delay]);
};
