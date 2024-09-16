import { createSlice } from "@reduxjs/toolkit";

export interface Credential {
	readonly token: string;
	readonly expires_at: number;
}

export const initialCredential: Credential = { token: "", expires_at: 0 };

const slice = createSlice({
	name: "credential",
	initialState: { credential: initialCredential },
	reducers: {
		onUpdate: (state, action) => {
			const { payload } = action;
			state.credential = payload;
		}
	}
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
