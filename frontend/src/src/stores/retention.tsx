import { createSlice } from "@reduxjs/toolkit";

import { globalConfig } from "../config/global";

const { retention } = globalConfig;
const { default: initialRetention } = retention;

const slice = createSlice({
	name: "retention",
	initialState: { retention: initialRetention },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			if (payload > retention.minimum && payload <= retention.maximum) {
				state.retention = payload;
			} else {
				state.retention = retention.default;
			}
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
