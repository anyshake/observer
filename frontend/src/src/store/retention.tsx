import { Dispatch, createSlice } from "@reduxjs/toolkit";
import GLOBAL_CONFIG from "../config/global";

const { retention } = GLOBAL_CONFIG.app_settings;

const slice = createSlice({
    name: "retention",
    initialState: {
        retention,
    },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            if (payload > 0 && payload <= 1000) {
                state.retention = payload;
            } else {
                state.retention = retention;
            }
        },
    },
});

const { onUpdate } = slice.actions;
const update = (newRetention: number) => (dispatch: Dispatch) => {
    dispatch(onUpdate(newRetention));
};

export { update };
export default slice.reducer;
