import { Dispatch, createSlice } from "@reduxjs/toolkit";
import { ADC } from "../config/adc";

const initialADC: ADC = {
    fullscale: 0,
    resolution: -1,
};

const slice = createSlice({
    name: "adc",
    initialState: { adc: initialADC },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            state.adc = payload;
        },
    },
});

const { onUpdate } = slice.actions;
const update = (newADC: ADC) => (dispatch: Dispatch) => {
    dispatch(onUpdate(newADC));
};

export { update };
export default slice.reducer;
