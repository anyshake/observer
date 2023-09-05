import { StyledEngineProvider } from "@mui/material/styles";
import { PersistGate } from "redux-persist/integration/react";
import store, { REDUX_PRESIST } from "./config/store";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import App from "./App";
import "./index.css";

const root = createRoot(document.getElementById("root") as HTMLElement);
root.render(
    <Provider store={store}>
        <PersistGate loading={null} persistor={REDUX_PRESIST}>
            <StyledEngineProvider injectFirst>
                <App />
            </StyledEngineProvider>
        </PersistGate>
    </Provider>
);
