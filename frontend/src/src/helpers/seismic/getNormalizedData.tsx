import BigNumber from "bignumber.js";

export const getNormalizedData = (data: number[], center: number) => {
	const arrSum = data.reduce((sum, val) => {
		const bigSum = new BigNumber(sum);
		const bigVal = new BigNumber(val);
		return bigSum.plus(bigVal).toNumber();
	});

	const avg = new BigNumber(arrSum).dividedBy(data.length).toNumber();
	return data.map((item) => {
		const bigItem = new BigNumber(item);
		const bigAvg = new BigNumber(avg);
		const bigCenter = new BigNumber(center);
		return bigItem.minus(bigAvg).plus(bigCenter).toNumber();
	});
};
