import { createSlice } from "@reduxjs/toolkit";
import { fallbackScale } from "../config/global";

const { value: initScale } = fallbackScale.property();

const slice = createSlice({
    name: "scale",
    initialState: { scale: initScale },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            state.scale = payload;
        },
    },
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
