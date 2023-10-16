import { Dispatch, createSlice } from "@reduxjs/toolkit";
import GLOBAL_CONFIG from "../config/global";

const { duration } = GLOBAL_CONFIG.app_settings;

const slice = createSlice({
    name: "duration",
    initialState: {
        duration,
    },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            if (payload > 0 && payload <= 3600) {
                state.duration = payload;
            } else {
                state.duration = duration;
            }
        },
    },
});

const { onUpdate } = slice.actions;
const update = (newDuration: number) => (dispatch: Dispatch) => {
    dispatch(onUpdate(newDuration));
};

export { update };
export default slice.reducer;
