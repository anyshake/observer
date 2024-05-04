import { createSlice } from "@reduxjs/toolkit";

import { globalConfig } from "../config/global";

const { duration } = globalConfig;
const { default: initialDuration } = duration;

const slice = createSlice({
	name: "duration",
	initialState: { duration: initialDuration },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			if (payload > duration.minimum && payload <= duration.maximum) {
				state.duration = payload;
			} else {
				state.duration = duration.default;
			}
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
