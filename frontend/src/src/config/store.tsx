import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from "redux";
import geophone from "../store/geophone";
import { Geophone } from "./geophone";
import adc from "../store/adc";
import { ADC } from "./adc";

const reducer = combineReducers({ adc, geophone });
const REDUX_STORE = configureStore({ reducer });
process.env.NODE_ENV !== "production" &&
    REDUX_STORE.subscribe(() => {
        console.log(REDUX_STORE.getState());
    });

export interface ReduxStoreProps {
    readonly adc: ReduxStore["adc"];
    readonly geophone: ReduxStore["geophone"];
    readonly updateADC: (adc: ADC) => void;
    readonly updateGeophone: (geophone: Geophone) => void;
}

export type ReduxStore = ReturnType<typeof reducer>;
export default REDUX_STORE;
