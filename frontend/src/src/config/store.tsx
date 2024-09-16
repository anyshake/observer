import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from "redux";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";

import adc, { ADC } from "../stores/adc";
import credential, { Credential } from "../stores/credential";
import duration from "../stores/duration";
import geophone, { Geophone } from "../stores/geophone";
import retention from "../stores/retention";
import scale from "../stores/scale";
import station, { Station } from "../stores/station";

const scalePersistConfig = persistReducer({ storage, key: "scale", whitelist: ["scale"] }, scale);
const durationPersistConfig = persistReducer(
	{ storage, key: "duration", whitelist: ["duration"] },
	duration
);
const retentionPersistConfig = persistReducer(
	{ storage, key: "retention", whitelist: ["retention"] },
	retention
);
const credentialPersistConfig = persistReducer(
	{ storage, key: "credential", whitelist: ["credential"] },
	credential
);

const reducer = combineReducers({
	adc,
	geophone,
	station,
	scale: scalePersistConfig,
	duration: durationPersistConfig,
	retention: retentionPersistConfig,
	credential: credentialPersistConfig
});
const REDUX_STORE = configureStore({
	reducer,
	middleware: (getDefaultMiddleware) => getDefaultMiddleware({ serializableCheck: false })
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
	readonly credential: ReturnType<typeof credential>;
	readonly updateCredential: (credential: Credential) => void;
}
export default REDUX_STORE;
