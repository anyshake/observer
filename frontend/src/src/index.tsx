import { createRoot } from "react-dom/client";
import App from "./App";
import { StyledEngineProvider } from "@mui/material/styles";
import "./index.css";

const root = createRoot(document.getElementById("root") as HTMLElement);
root.render(
    <StyledEngineProvider injectFirst>
        <App />
    </StyledEngineProvider>
);
