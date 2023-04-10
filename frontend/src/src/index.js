import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import { Provider } from "react-redux";
import store from "./helpers/redux/store";
import {
    hookAddEventListener,
    hookRemoveEventListener,
} from "./helpers/events/eventListenerHook";

hookAddEventListener();
hookRemoveEventListener();

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
    // <React.StrictMode>
    <Provider store={store}>
        <App />
    </Provider>
    // </React.StrictMode>
);
