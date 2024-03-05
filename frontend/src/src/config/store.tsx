import { combineReducers } from "redux";
import storage from "redux-persist/lib/storage";
import { configureStore } from "@reduxjs/toolkit";
import { persistReducer, persistStore } from "redux-persist";
import scale from "../stores/scale";
import adc, { ADC } from "../stores/adc";
import duration from "../stores/duration";
import retention from "../stores/retention";
import station, { Station } from "../stores/station";
import geophone, { Geophone } from "../stores/geophone";

const scalePersistConfig = persistReducer(
    { storage, key: "scale", whitelist: ["scale"] },
    scale
);
const durationPersistConfig = persistReducer(
    { storage, key: "duration", whitelist: ["duration"] },
    duration
);
const retentionPersistConfig = persistReducer(
    { storage, key: "retention", whitelist: ["retention"] },
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
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: false,
        }),
});

export const REDUX_PRESIST = persistStore(REDUX_STORE);
export type ReduxStore = ReturnType<typeof reducer>;
export interface ReduxStoreProps {
    readonly adc: ReturnType<typeof adc>;
    readonly updateADC: (adc: ADC) => void;
    readonly scale: ReturnType<typeof scale>;
    readonly updateScale: (scale: string) => void;
    readonly station: ReturnType<typeof station>;
    readonly updateStation: (station: Station) => void;
    readonly duration: ReturnType<typeof duration>;
    readonly updateDuration: (duration: number) => void;
    readonly geophone: ReturnType<typeof geophone>;
    readonly updateGeophone: (geophone: Geophone) => void;
    readonly retention: ReturnType<typeof retention>;
    readonly updateRetention: (retention: number) => void;
}
export default REDUX_STORE;
