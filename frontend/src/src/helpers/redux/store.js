import { legacy_createStore as createStore } from "redux";
import reducer from "./reducer";

const store = createStore(reducer);

process.env.NODE_ENV !== "production" &&
    store.subscribe(() => {
        console.log(store.getState());
    });

export default store;
