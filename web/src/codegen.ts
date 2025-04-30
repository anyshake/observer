/* eslint-disable @typescript-eslint/no-explicit-any */
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
    schema: {
        [`${(import.meta as any).env.VITE_APP_BACKEND_BASE_HOST}${(import.meta as any).env.VITE_APP_GRAPHQL_API_ENDPOINT}`]:
            {
                headers: {
                    Authorization: `Bearer ${(import.meta as any).env.VITE_APP_GRAPHQL_API_ENDPOINT_DEV_TOKEN}`
                }
            }
    },
    overwrite: true,
    documents: 'src/**/*.graphql',
    generates: {
        'src/graphql.tsx': {
            plugins: ['typescript', 'typescript-operations', 'typescript-react-apollo'],
            config: {
                withHooks: true,
                preResolveTypes: true,
                withHOC: false,
                withComponent: false,
                skipTypename: false,
                apolloReactCommonImportFrom: '@apollo/client',
                apolloReactHooksImportFrom: '@apollo/client',
                scalars: { Int64: 'number', Int32: 'number' }
            }
        }
    }
};

export default config;
