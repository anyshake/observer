import { createSlice } from "@reduxjs/toolkit";

export interface Geophone {
    readonly ehz: number;
    readonly ehe: number;
    readonly ehn: number;
    readonly initialized: boolean;
}

const initialGeophone: Geophone = {
    ehz: 0.288,
    ehe: 0.288,
    ehn: 0.288,
    initialized: false,
};

const slice = createSlice({
    name: "geophone",
    initialState: { geophone: initialGeophone },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            state.geophone = payload;
        },
    },
});

export const { onUpdate } = slice.actions;
export default slice.reducer;
