import { createSlice } from "@reduxjs/toolkit";

export interface Stream {
	readonly station: string;
	readonly network: string;
	readonly location: string;
	readonly channel: string;
	readonly initialized: boolean;
}

const initialStation: Stream = {
	station: "SHAKE",
	network: "AS",
	location: "00",
	channel: "EH",
	initialized: false
};

const slice = createSlice({
	name: "stream",
	initialState: { stream: initialStation },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			state.stream = payload;
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
