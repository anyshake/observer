import { createSlice } from "@reduxjs/toolkit";

export interface Sensor {
	readonly resolution: number;
	readonly velocity: boolean;
	readonly initialized: boolean;
}

const initialSensor: Sensor = {
	resolution: 24,
	velocity: true,
	initialized: false
};

const slice = createSlice({
	name: "sensor",
	initialState: { sensor: initialSensor },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			state.sensor = payload;
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
