import { createSlice } from "@reduxjs/toolkit";

export interface Geophone {
	readonly frequency: number;
	readonly sensitivity: number;
	readonly initialized: boolean;
}

const initialGeophone: Geophone = {
	frequency: 4.5,
	sensitivity: 28.8,
	initialized: false
};

const slice = createSlice({
	name: "geophone",
	initialState: { geophone: initialGeophone },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			state.geophone = payload;
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
