import { Dispatch, createSlice } from "@reduxjs/toolkit";
import { Station } from "../config/station";

const initialStation: Station = {
    station: "SHAKE",
    network: "AS",
    location: "00",
};

const slice = createSlice({
    name: "station",
    initialState: { station: initialStation },
    reducers: {
        onUpdate: (state, action) => {
            const { payload } = action;
            state.station = payload;
        },
    },
});

const { onUpdate } = slice.actions;
const update = (newStation: Station) => (dispatch: Dispatch) => {
    dispatch(onUpdate(newStation));
};

export { update };
export default slice.reducer;
