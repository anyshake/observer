import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from "redux";
import { Geophone } from "./geophone";
import { ADC } from "./adc";
import adc from "../store/adc";
import geophone from "../store/geophone";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";
import scale from "../store/scale";
import duration from "../store/duration";
import retention from "../store/retention";
import station from "../store/station";
import { Station } from "./station";

const scalePersistConfig = persistReducer(
    {
        storage,
        key: "scale",
        whitelist: ["scale"],
    },
    scale
);
const durationPersistConfig = persistReducer(
    {
        storage,
        key: "duration",
        whitelist: ["duration"],
    },
    duration
);
const retentionPersistConfig = persistReducer(
    {
        storage,
        key: "retention",
        whitelist: ["retention"],
    },
    retention
);

const reducer = combineReducers({
    adc,
    geophone,
    station,
    scale: scalePersistConfig,
    duration: durationPersistConfig,
    retention: retentionPersistConfig,
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
    readonly station: ReduxStore["station"];
    readonly duration: ReduxStore["duration"];
    readonly geophone: ReduxStore["geophone"];
    readonly retention: ReduxStore["retention"];
    readonly updateADC?: (adc: ADC) => void;
    readonly updateScale?: (scale: string) => void;
    readonly updateStation?: (station: Station) => void;
    readonly updateDuration?: (duration: number) => void;
    readonly updateGeophone?: (geophone: Geophone) => void;
    readonly updateRetention?: (retention: number) => void;
}

export type ReduxStore = ReturnType<typeof reducer>;
export default REDUX_STORE;
export { REDUX_PRESIST };
