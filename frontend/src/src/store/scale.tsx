import { Dispatch, createSlice } from "@reduxjs/toolkit";
import { fallbackScale } from "../config/global";
import { IntensityStandardProperty } from "../helpers/seismic/intensityStandard";

const initialScale = fallbackScale.property();

const slice = createSlice({
    name: "scale",
    initialState: {
        scale: initialScale,
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
    (newScale: IntensityStandardProperty) => (dispatch: Dispatch) => {
        dispatch(onUpdate(newScale));
    };

export { update };
export default slice.reducer;
