import * as OregonDSPTop from "oregondsp";

export enum FilterPassband {
	LOW_PASS,
	HIGH_PASS,
	BAND_PASS
}

export interface FilterOptions {
	readonly poles: number;
	readonly sampleRate: number;
	readonly lowFreqCorner: number;
	readonly highFreqCorner: number;
	readonly passbandType: FilterPassband;
}

export const getFilteredCounts = (record: number[], options: FilterOptions) => {
	const { passbandType, poles, sampleRate, lowFreqCorner, highFreqCorner } = options;

	let passband: OregonDSPTop.com.oregondsp.signalProcessing.filter.iir.PassbandType;
	if (passbandType === FilterPassband.LOW_PASS) {
		passband = OregonDSPTop.com.oregondsp.signalProcessing.filter.iir.PassbandType.LOWPASS;
	} else if (passbandType === FilterPassband.HIGH_PASS) {
		passband = OregonDSPTop.com.oregondsp.signalProcessing.filter.iir.PassbandType.HIGHPASS;
	} else {
		passband = OregonDSPTop.com.oregondsp.signalProcessing.filter.iir.PassbandType.BANDPASS;
	}

	const butterworth = new OregonDSPTop.com.oregondsp.signalProcessing.filter.iir.Butterworth(
		poles,
		passband,
		lowFreqCorner,
		highFreqCorner,
		1 / sampleRate
	);
	const float32ArrayBuffer = new Float32Array(record);
	butterworth.filterInPlace(float32ArrayBuffer);

	return Array.from(float32ArrayBuffer);
};
