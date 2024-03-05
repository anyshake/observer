import { createSlice } from "@reduxjs/toolkit";

export interface ADC {
    readonly resolution: number;
    readonly fullscale: number;
    readonly initialized: boolean;
}

const initialADC: ADC = { fullscale: 5, resolution: 24, initialized: false };

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

export const { onUpdate } = slice.actions;
export default slice.reducer;
