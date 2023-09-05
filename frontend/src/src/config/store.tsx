import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from "redux";
import { Geophone } from "./geophone";
import { ADC } from "./adc";
import adc from "../store/adc";
import geophone from "../store/geophone";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";
import scale from "../store/scale";
import { IntensityStandardProperty } from "../helpers/seismic/intensityStandard";

const scalePersistConfig = persistReducer(
    {
        storage,
        key: "scale",
        whitelist: ["scale"],
    },
    scale
);
const reducer = combineReducers({
    adc,
    geophone,
    scale: scalePersistConfig,
});

const REDUX_STORE = configureStore({
    reducer,
    middleware: (getDefaultMiddleware) => {
        return getDefaultMiddleware({
            serializableCheck: false,
        });
    },
});
const REDUX_PRESIST = persistStore(REDUX_STORE);

export interface ReduxStoreProps {
    readonly adc: ReduxStore["adc"];
    readonly scale: ReduxStore["scale"];
    readonly geophone: ReduxStore["geophone"];
    readonly updateADC?: (adc: ADC) => void;
    readonly updateGeophone?: (geophone: Geophone) => void;
    readonly updateScale?: (scale: IntensityStandardProperty) => void;
}

export type ReduxStore = ReturnType<typeof reducer>;
export default REDUX_STORE;
export { REDUX_PRESIST };
