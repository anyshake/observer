import { createSlice } from "@reduxjs/toolkit";

export interface Station {
	readonly station: string;
	readonly network: string;
	readonly location: string;
	readonly channel: string;
	readonly initialized: boolean;
}

const initialStation: Station = {
	station: "SHAKE",
	network: "AS",
	location: "00",
	channel: "EH",
	initialized: false
};

const slice = createSlice({
	name: "station",
	initialState: { station: initialStation },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			state.station = payload;
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
