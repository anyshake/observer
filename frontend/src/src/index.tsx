import { createRoot } from "react-dom/client";
import App from "./App";
import { StyledEngineProvider } from "@mui/material/styles";
import { Provider } from "react-redux";
import store from "./config/store";
import "./index.css";

const root = createRoot(document.getElementById("root") as HTMLElement);
root.render(
    <Provider store={store}>
        <StyledEngineProvider injectFirst>
            <App />
        </StyledEngineProvider>
    </Provider>
);
