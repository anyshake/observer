export const getVoltageArr = (arr: number[], resolution: number, fullscale: number): number[] => {
	const factor = 2 ** (resolution - 1);
	const voltage: number[] = [];
	for (const i of arr) {
		voltage.push((fullscale / factor) * i);
	}

	return voltage;
};
