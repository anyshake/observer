import './index.css';

import { ApolloClient, ApolloProvider, InMemoryCache } from '@apollo/client';
import { createRoot } from 'react-dom/client';
import { ErrorBoundary } from 'react-error-boundary';

import App from './App.tsx';
import { ErrorPage } from './components/ui/ErrorPage.tsx';
import { RouterWrapper } from './components/ui/RouterWrapper.tsx';
import { routerConfig } from './config/router';
import { getGraphQlApiUrl } from './helpers/app/getGraphQlApiUrl.tsx';
import { ApiClient } from './helpers/request/ApiClient.tsx';

const graphQlClient = new ApolloClient({
    defaultOptions: {
        watchQuery: { fetchPolicy: 'network-only' },
        query: { fetchPolicy: 'network-only' }
    },
    cache: new InMemoryCache(),
    link: ApiClient.createGraphQlLink(getGraphQlApiUrl())
});

createRoot(document.getElementById('root')!).render(
    <ErrorBoundary fallback={<ErrorPage />}>
        <ApolloProvider client={graphQlClient}>
            <RouterWrapper mode={routerConfig.mode} basename={routerConfig.basename}>
                <App />
            </RouterWrapper>
        </ApolloProvider>
    </ErrorBoundary>
);
