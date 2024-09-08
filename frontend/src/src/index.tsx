import "./config/i18n";
import "./index.css";

import { StyledEngineProvider } from "@mui/material/styles";
import ReactDOM from "react-dom/client";
import { ErrorBoundary } from "react-error-boundary";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";

import App from "./App";
import { RouterWrapper } from "./components/RouterWrapper";
import { routerConfig } from "./config/router";
import store, { REDUX_PRESIST } from "./config/store";
import { Error } from "./Error";

const root = ReactDOM.createRoot(document.getElementById("root")!);
root.render(
	<ErrorBoundary  fallback={<Error />}>
		<Provider store={store}>
			<PersistGate loading={null} persistor={REDUX_PRESIST}>
				<StyledEngineProvider injectFirst>
					<RouterWrapper mode={routerConfig.mode} basename={routerConfig.basename}>
						<App />
					</RouterWrapper>
				</StyledEngineProvider>
			</PersistGate>
		</Provider>
	</ErrorBoundary>
);
