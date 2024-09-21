import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from "redux";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";

import credential, { Credential } from "../stores/credential";
import duration from "../stores/duration";
import retention from "../stores/retention";
import sensor, { Sensor } from "../stores/sensor";
import stream, { Stream } from "../stores/stream";

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
	sensor,
	stream,
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
	readonly sensor: ReturnType<typeof sensor>;
	readonly updateSensor: (sensor: Sensor) => void;
	readonly stream: ReturnType<typeof stream>;
	readonly updateStream: (stream: Stream) => void;
	readonly duration: ReturnType<typeof duration>;
	readonly updateDuration: (duration: number) => void;
	readonly retention: ReturnType<typeof retention>;
	readonly updateRetention: (retention: number) => void;
	readonly credential: ReturnType<typeof credential>;
	readonly updateCredential: (credential: Credential) => void;
}
export default REDUX_STORE;
