import { Dispatch, createSlice } from "@reduxjs/toolkit";
import { Geophone } from "../config/geophone";

const initialGeophone: Geophone = {
    ehz: 0,
    ehe: 0,
    ehn: 0,
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

const { onUpdate } = slice.actions;
const update = (newGeophone: Geophone) => (dispatch: Dispatch) => {
    dispatch(onUpdate(newGeophone));
};

export { update };
export default slice.reducer;
