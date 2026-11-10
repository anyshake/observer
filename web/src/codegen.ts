import type { CodegenConfig } from '@graphql-codegen/cli';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const schemaPath = (import.meta as any).env.VITE_APP_GRAPHQL_SCHEMA_PATH;
if (!schemaPath) {
    throw new Error('VITE_APP_GRAPHQL_SCHEMA_PATH must be set to the backend GraphQL schema path.');
}

const config: CodegenConfig = {
    schema: schemaPath,
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
