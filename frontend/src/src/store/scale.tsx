import { Dispatch, createSlice } from "@reduxjs/toolkit";
import { fallbackScale } from "../config/global";

const { value: scale } = fallbackScale.property();

const slice = createSlice({
    name: "scale",
    initialState: {
        scale,
    },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            state.scale = payload;
        },
    },
});

const { onUpdate } = slice.actions;
const update =
    (newScale: string) => (dispatch: Dispatch) => {
        dispatch(onUpdate(newScale));
    };

export { update };
export default slice.reducer;
