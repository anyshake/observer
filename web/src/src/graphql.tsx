import { gql } from '@apollo/client';
import * as ApolloReactCommon from '@apollo/client';
import * as ApolloReactHooks from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Any: { input: any; output: any; }
  Int32: { input: number; output: number; }
  Int64: { input: number; output: number; }
  Map: { input: any; output: any; }
};

export type Mutation = {
  __typename?: 'Mutation';
  createSysUser: Scalars['String']['output'];
  importGlobalConfig: Scalars['Boolean']['output'];
  purgeHelicorderFiles: Scalars['Boolean']['output'];
  purgeMiniSeedFiles: Scalars['Boolean']['output'];
  purgeSeisRecords: Scalars['Boolean']['output'];
  removeSysUser: Scalars['Boolean']['output'];
  restartService: Scalars['Boolean']['output'];
  restoreServiceConfig: Scalars['Boolean']['output'];
  restoreStationConfig: Scalars['Boolean']['output'];
  startService: Scalars['Boolean']['output'];
  stopService: Scalars['Boolean']['output'];
  updateServiceConfig: Scalars['Boolean']['output'];
  updateStationConfig: Scalars['Boolean']['output'];
  updateSysUser: Scalars['Boolean']['output'];
};


export type MutationCreateSysUserArgs = {
  admin: Scalars['Boolean']['input'];
  password: Scalars['String']['input'];
  username: Scalars['String']['input'];
};


export type MutationImportGlobalConfigArgs = {
  data: Scalars['String']['input'];
};


export type MutationRemoveSysUserArgs = {
  userId: Scalars['String']['input'];
};


export type MutationRestartServiceArgs = {
  serviceId: Scalars['String']['input'];
};


export type MutationRestoreServiceConfigArgs = {
  serviceId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationStartServiceArgs = {
  serviceId: Scalars['String']['input'];
};


export type MutationStopServiceArgs = {
  serviceId: Scalars['String']['input'];
};


export type MutationUpdateServiceConfigArgs = {
  key: Scalars['String']['input'];
  serviceId: Scalars['String']['input'];
  val: Scalars['Any']['input'];
};


export type MutationUpdateStationConfigArgs = {
  key: Scalars['String']['input'];
  value: Scalars['Any']['input'];
};


export type MutationUpdateSysUserArgs = {
  admin: Scalars['Boolean']['input'];
  password?: InputMaybe<Scalars['String']['input']>;
  userId: Scalars['String']['input'];
  username: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  exportGlobalConfig: Scalars['String']['output'];
  getApplicationLogs?: Maybe<Array<Scalars['String']['output']>>;
  getCurrentTime: Scalars['Int64']['output'];
  getCurrentUser: SysUser;
  getDeviceConfig: DeviceConfig;
  getDeviceId: Scalars['String']['output'];
  getDeviceInfo: DeviceInfo;
  getDeviceStatus: DeviceStatus;
  getEventSource: Array<Maybe<SeisEventSource>>;
  getEventsBySource: Array<SeisEvent>;
  getHelicorderFiles: Array<Maybe<ServiceAsset>>;
  getMiniSeedFiles: Array<Maybe<ServiceAsset>>;
  getSeisRecordsByTime: Array<Maybe<SeisRecord>>;
  getServiceConfigConstraint: Array<ServiceConfigConstraint>;
  getServiceStatus?: Maybe<Array<ServiceStatus>>;
  getSoftwareVersion: Scalars['String']['output'];
  getStationConfig: Scalars['Map']['output'];
  getStationConfigConstraint: Array<ConfigConstraint>;
  getStationMetadata: Scalars['String']['output'];
  getSysUsers: Array<SysUser>;
  getSystemStatus: SystemStatus;
  isGenuineProduct: Scalars['Boolean']['output'];
};


export type QueryGetEventsBySourceArgs = {
  code: Scalars['String']['input'];
};


export type QueryGetSeisRecordsByTimeArgs = {
  endTime: Scalars['Int64']['input'];
  startTime: Scalars['Int64']['input'];
};


export type QueryGetStationMetadataArgs = {
  format: Scalars['String']['input'];
};

export type _ChannelData = {
  __typename?: '_channelData';
  channelCode: Scalars['String']['output'];
  channelId: Scalars['Int']['output'];
  data: Array<Scalars['Int32']['output']>;
};

export type ConfigConstraint = {
  __typename?: 'configConstraint';
  configType: Scalars['String']['output'];
  currentValue: Scalars['Any']['output'];
  description: Scalars['String']['output'];
  isRequired: Scalars['Boolean']['output'];
  key: Scalars['String']['output'];
  name: Scalars['String']['output'];
  namespace: Scalars['String']['output'];
  options?: Maybe<Scalars['Any']['output']>;
};

export type DeviceConfig = {
  __typename?: 'deviceConfig';
  channelCodes: Array<Scalars['String']['output']>;
  gnssEnabled: Scalars['Boolean']['output'];
  model: Scalars['String']['output'];
  packetInterval: Scalars['Int64']['output'];
  protocol: Scalars['String']['output'];
  sampleRate: Scalars['Int']['output'];
};

export type DeviceInfo = {
  __typename?: 'deviceInfo';
  elevation?: Maybe<Scalars['Float']['output']>;
  latitude?: Maybe<Scalars['Float']['output']>;
  longitude?: Maybe<Scalars['Float']['output']>;
  temperature?: Maybe<Scalars['Float']['output']>;
};

export type DeviceStatus = {
  __typename?: 'deviceStatus';
  errors: Scalars['Int64']['output'];
  frames: Scalars['Int64']['output'];
  messages: Scalars['Int64']['output'];
  startedAt: Scalars['Int64']['output'];
  updatedAt: Scalars['Int64']['output'];
};

export type SeisEvent = {
  __typename?: 'seisEvent';
  depth: Scalars['Float']['output'];
  distance: Scalars['Float']['output'];
  estimation: Array<Scalars['Float']['output']>;
  eventId: Scalars['String']['output'];
  latitude: Scalars['Float']['output'];
  longitude: Scalars['Float']['output'];
  magnitude: Scalars['Map']['output'];
  region: Scalars['String']['output'];
  timestamp: Scalars['Int64']['output'];
  verfied: Scalars['Boolean']['output'];
};

export type SeisEventSource = {
  __typename?: 'seisEventSource';
  country: Scalars['String']['output'];
  defaultLocale: Scalars['String']['output'];
  id: Scalars['String']['output'];
  locales: Scalars['Map']['output'];
};

export type SeisRecord = {
  __typename?: 'seisRecord';
  channelData: Array<_ChannelData>;
  sampleRate: Scalars['Int']['output'];
  timestamp: Scalars['Int64']['output'];
};

export type ServiceAsset = {
  __typename?: 'serviceAsset';
  fileName: Scalars['String']['output'];
  filePath: Scalars['String']['output'];
  modifiedAt: Scalars['Int64']['output'];
  namespace: Scalars['String']['output'];
  size: Scalars['Int64']['output'];
};

export type ServiceConfigConstraint = {
  __typename?: 'serviceConfigConstraint';
  constraints: Array<Maybe<ConfigConstraint>>;
  serviceId: Scalars['String']['output'];
  serviceName: Scalars['String']['output'];
};

export type ServiceStatus = {
  __typename?: 'serviceStatus';
  description: Scalars['String']['output'];
  isRunning: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  restarts: Scalars['Int']['output'];
  serviceId: Scalars['String']['output'];
  startedAt: Scalars['Int64']['output'];
  stoppedAt: Scalars['Int64']['output'];
  updatedAt: Scalars['Int64']['output'];
};

export type SysUser = {
  __typename?: 'sysUser';
  admin: Scalars['Boolean']['output'];
  createdAt: Scalars['Int64']['output'];
  lastLogin: Scalars['Int64']['output'];
  updatedAt: Scalars['Int64']['output'];
  userAgent: Scalars['String']['output'];
  userId: Scalars['String']['output'];
  userIp: Scalars['String']['output'];
  username: Scalars['String']['output'];
};

export type SystemStatus = {
  __typename?: 'systemStatus';
  cpu: Scalars['Float']['output'];
  disk: Scalars['Float']['output'];
  memory: Scalars['Float']['output'];
  uptime: Scalars['Int64']['output'];
};

export type GetSoftwareVersionQueryVariables = Exact<{ [key: string]: never; }>;


export type GetSoftwareVersionQuery = { __typename?: 'Query', getSoftwareVersion: string };

export type IsGenuineProductQueryVariables = Exact<{ [key: string]: never; }>;


export type IsGenuineProductQuery = { __typename?: 'Query', isGenuineProduct: boolean };

export type GetFileListDataQueryVariables = Exact<{ [key: string]: never; }>;


export type GetFileListDataQuery = { __typename?: 'Query', getMiniSeedFiles: Array<{ __typename?: 'serviceAsset', namespace: string, size: number, filePath: string, fileName: string, modifiedAt: number } | null>, getHelicorderFiles: Array<{ __typename?: 'serviceAsset', namespace: string, size: number, filePath: string, fileName: string, modifiedAt: number } | null> };

export type GetSeismicEventSourceListQueryVariables = Exact<{ [key: string]: never; }>;


export type GetSeismicEventSourceListQuery = { __typename?: 'Query', getCurrentTime: number, getEventSource: Array<{ __typename?: 'seisEventSource', id: string, country: string, locales: any, defaultLocale: string } | null> };

export type GetSeismicEventBySourceQueryVariables = Exact<{
  sourceId: Scalars['String']['input'];
}>;


export type GetSeismicEventBySourceQuery = { __typename?: 'Query', getEventsBySource: Array<{ __typename?: 'seisEvent', verfied: boolean, timestamp: number, eventId: string, region: string, depth: number, latitude: number, longitude: number, magnitude: any, distance: number, estimation: Array<number> }> };

export type GetSeismicRecordsQueryVariables = Exact<{
  startTime: Scalars['Int64']['input'];
  endTime: Scalars['Int64']['input'];
}>;


export type GetSeismicRecordsQuery = { __typename?: 'Query', getSeisRecordsByTime: Array<{ __typename?: 'seisRecord', timestamp: number, sampleRate: number, channelData: Array<{ __typename?: '_channelData', channelCode: string, channelId: number, data: Array<number> }> } | null> };

export type GetHomeDataQueryVariables = Exact<{ [key: string]: never; }>;


export type GetHomeDataQuery = { __typename?: 'Query', getCurrentTime: number, getStationConfig: any, getDeviceId: string, getDeviceConfig: { __typename?: 'deviceConfig', packetInterval: number, sampleRate: number, channelCodes: Array<string>, gnssEnabled: boolean, model: string, protocol: string }, getDeviceInfo: { __typename?: 'deviceInfo', latitude?: number | null, longitude?: number | null, elevation?: number | null, temperature?: number | null }, getDeviceStatus: { __typename?: 'deviceStatus', startedAt: number, updatedAt: number, frames: number, errors: number, messages: number }, getServiceStatus?: Array<{ __typename?: 'serviceStatus', serviceId: string, isRunning: boolean, name: string, description: string }> | null, getSystemStatus: { __typename?: 'systemStatus', cpu: number, memory: number, disk: number, uptime: number } };

export type CreateUserMutationVariables = Exact<{
  username: Scalars['String']['input'];
  password: Scalars['String']['input'];
  admin: Scalars['Boolean']['input'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createSysUser: string };

export type ExportGlobalConfigQueryVariables = Exact<{ [key: string]: never; }>;


export type ExportGlobalConfigQuery = { __typename?: 'Query', exportGlobalConfig: string };

export type GetApplicationLogsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetApplicationLogsQuery = { __typename?: 'Query', getApplicationLogs?: Array<string> | null };

export type GetServiceDataQueryVariables = Exact<{ [key: string]: never; }>;


export type GetServiceDataQuery = { __typename?: 'Query', getServiceConfigConstraint: Array<{ __typename?: 'serviceConfigConstraint', serviceName: string, serviceId: string, constraints: Array<{ __typename?: 'configConstraint', key: string, name: string, description: string, configType: string, isRequired: boolean, currentValue: any, options?: any | null } | null> }>, getServiceStatus?: Array<{ __typename?: 'serviceStatus', description: string, serviceId: string, name: string, restarts: number, startedAt: number, stoppedAt: number, updatedAt: number, isRunning: boolean }> | null };

export type GetStationConfigQueryVariables = Exact<{ [key: string]: never; }>;


export type GetStationConfigQuery = { __typename?: 'Query', getStationConfigConstraint: Array<{ __typename?: 'configConstraint', namespace: string, name: string, key: string, description: string, configType: string, isRequired: boolean, currentValue: any, options?: any | null }> };

export type GetStationMetadataQueryVariables = Exact<{
  format: Scalars['String']['input'];
}>;


export type GetStationMetadataQuery = { __typename?: 'Query', getStationMetadata: string };

export type GetUserListQueryVariables = Exact<{ [key: string]: never; }>;


export type GetUserListQuery = { __typename?: 'Query', getSysUsers: Array<{ __typename?: 'sysUser', userId: string, username: string, createdAt: number, updatedAt: number, lastLogin: number, userIp: string, userAgent: string, admin: boolean }> };

export type ImportGlobalConfigMutationVariables = Exact<{
  data: Scalars['String']['input'];
}>;


export type ImportGlobalConfigMutation = { __typename?: 'Mutation', importGlobalConfig: boolean };

export type IsCurrentUserAdminQueryVariables = Exact<{ [key: string]: never; }>;


export type IsCurrentUserAdminQuery = { __typename?: 'Query', getCurrentUser: { __typename?: 'sysUser', admin: boolean } };

export type PurgeSeisRecordsMutationVariables = Exact<{ [key: string]: never; }>;


export type PurgeSeisRecordsMutation = { __typename?: 'Mutation', purgeSeisRecords: boolean };

export type PurgeMiniSeedFilesMutationVariables = Exact<{ [key: string]: never; }>;


export type PurgeMiniSeedFilesMutation = { __typename?: 'Mutation', purgeMiniSeedFiles: boolean };

export type PurgeHelicorderFilesMutationVariables = Exact<{ [key: string]: never; }>;


export type PurgeHelicorderFilesMutation = { __typename?: 'Mutation', purgeHelicorderFiles: boolean };

export type RestoreStationConfigMutationVariables = Exact<{ [key: string]: never; }>;


export type RestoreStationConfigMutation = { __typename?: 'Mutation', restoreStationConfig: boolean };

export type RestoreServiceConfigMutationVariables = Exact<{
  serviceId?: InputMaybe<Scalars['String']['input']>;
}>;


export type RestoreServiceConfigMutation = { __typename?: 'Mutation', restoreServiceConfig: boolean };

export type RemoveUserMutationVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type RemoveUserMutation = { __typename?: 'Mutation', removeSysUser: boolean };

export type RestartServiceMutationVariables = Exact<{
  serviceId: Scalars['String']['input'];
}>;


export type RestartServiceMutation = { __typename?: 'Mutation', restartService: boolean };

export type StartServiceMutationVariables = Exact<{
  serviceId: Scalars['String']['input'];
}>;


export type StartServiceMutation = { __typename?: 'Mutation', startService: boolean };

export type StopServiceMutationVariables = Exact<{
  serviceId: Scalars['String']['input'];
}>;


export type StopServiceMutation = { __typename?: 'Mutation', stopService: boolean };

export type UpdateServiceConfigMutationVariables = Exact<{
  serviceId: Scalars['String']['input'];
  key: Scalars['String']['input'];
  val: Scalars['Any']['input'];
}>;


export type UpdateServiceConfigMutation = { __typename?: 'Mutation', updateServiceConfig: boolean };

export type UpdateStationConfigMutationVariables = Exact<{
  key: Scalars['String']['input'];
  value: Scalars['Any']['input'];
}>;


export type UpdateStationConfigMutation = { __typename?: 'Mutation', updateStationConfig: boolean };

export type UpdateUserMutationVariables = Exact<{
  userId: Scalars['String']['input'];
  username: Scalars['String']['input'];
  password?: InputMaybe<Scalars['String']['input']>;
  admin: Scalars['Boolean']['input'];
}>;


export type UpdateUserMutation = { __typename?: 'Mutation', updateSysUser: boolean };


export const GetSoftwareVersionDocument = gql`
    query getSoftwareVersion {
  getSoftwareVersion
}
    `;

/**
 * __useGetSoftwareVersionQuery__
 *
 * To run a query within a React component, call `useGetSoftwareVersionQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSoftwareVersionQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSoftwareVersionQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetSoftwareVersionQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>(GetSoftwareVersionDocument, options);
      }
export function useGetSoftwareVersionLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>(GetSoftwareVersionDocument, options);
        }
export function useGetSoftwareVersionSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>(GetSoftwareVersionDocument, options);
        }
export type GetSoftwareVersionQueryHookResult = ReturnType<typeof useGetSoftwareVersionQuery>;
export type GetSoftwareVersionLazyQueryHookResult = ReturnType<typeof useGetSoftwareVersionLazyQuery>;
export type GetSoftwareVersionSuspenseQueryHookResult = ReturnType<typeof useGetSoftwareVersionSuspenseQuery>;
export type GetSoftwareVersionQueryResult = ApolloReactCommon.QueryResult<GetSoftwareVersionQuery, GetSoftwareVersionQueryVariables>;
export const IsGenuineProductDocument = gql`
    query isGenuineProduct {
  isGenuineProduct
}
    `;

/**
 * __useIsGenuineProductQuery__
 *
 * To run a query within a React component, call `useIsGenuineProductQuery` and pass it any options that fit your needs.
 * When your component renders, `useIsGenuineProductQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useIsGenuineProductQuery({
 *   variables: {
 *   },
 * });
 */
export function useIsGenuineProductQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<IsGenuineProductQuery, IsGenuineProductQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<IsGenuineProductQuery, IsGenuineProductQueryVariables>(IsGenuineProductDocument, options);
      }
export function useIsGenuineProductLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<IsGenuineProductQuery, IsGenuineProductQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<IsGenuineProductQuery, IsGenuineProductQueryVariables>(IsGenuineProductDocument, options);
        }
export function useIsGenuineProductSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<IsGenuineProductQuery, IsGenuineProductQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<IsGenuineProductQuery, IsGenuineProductQueryVariables>(IsGenuineProductDocument, options);
        }
export type IsGenuineProductQueryHookResult = ReturnType<typeof useIsGenuineProductQuery>;
export type IsGenuineProductLazyQueryHookResult = ReturnType<typeof useIsGenuineProductLazyQuery>;
export type IsGenuineProductSuspenseQueryHookResult = ReturnType<typeof useIsGenuineProductSuspenseQuery>;
export type IsGenuineProductQueryResult = ApolloReactCommon.QueryResult<IsGenuineProductQuery, IsGenuineProductQueryVariables>;
export const GetFileListDataDocument = gql`
    query getFileListData {
  getMiniSeedFiles {
    namespace
    size
    filePath
    fileName
    modifiedAt
  }
  getHelicorderFiles {
    namespace
    size
    filePath
    fileName
    modifiedAt
  }
}
    `;

/**
 * __useGetFileListDataQuery__
 *
 * To run a query within a React component, call `useGetFileListDataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetFileListDataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetFileListDataQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetFileListDataQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetFileListDataQuery, GetFileListDataQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetFileListDataQuery, GetFileListDataQueryVariables>(GetFileListDataDocument, options);
      }
export function useGetFileListDataLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetFileListDataQuery, GetFileListDataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetFileListDataQuery, GetFileListDataQueryVariables>(GetFileListDataDocument, options);
        }
export function useGetFileListDataSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetFileListDataQuery, GetFileListDataQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetFileListDataQuery, GetFileListDataQueryVariables>(GetFileListDataDocument, options);
        }
export type GetFileListDataQueryHookResult = ReturnType<typeof useGetFileListDataQuery>;
export type GetFileListDataLazyQueryHookResult = ReturnType<typeof useGetFileListDataLazyQuery>;
export type GetFileListDataSuspenseQueryHookResult = ReturnType<typeof useGetFileListDataSuspenseQuery>;
export type GetFileListDataQueryResult = ApolloReactCommon.QueryResult<GetFileListDataQuery, GetFileListDataQueryVariables>;
export const GetSeismicEventSourceListDocument = gql`
    query getSeismicEventSourceList {
  getCurrentTime
  getEventSource {
    id
    country
    locales
    defaultLocale
  }
}
    `;

/**
 * __useGetSeismicEventSourceListQuery__
 *
 * To run a query within a React component, call `useGetSeismicEventSourceListQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSeismicEventSourceListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSeismicEventSourceListQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetSeismicEventSourceListQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>(GetSeismicEventSourceListDocument, options);
      }
export function useGetSeismicEventSourceListLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>(GetSeismicEventSourceListDocument, options);
        }
export function useGetSeismicEventSourceListSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>(GetSeismicEventSourceListDocument, options);
        }
export type GetSeismicEventSourceListQueryHookResult = ReturnType<typeof useGetSeismicEventSourceListQuery>;
export type GetSeismicEventSourceListLazyQueryHookResult = ReturnType<typeof useGetSeismicEventSourceListLazyQuery>;
export type GetSeismicEventSourceListSuspenseQueryHookResult = ReturnType<typeof useGetSeismicEventSourceListSuspenseQuery>;
export type GetSeismicEventSourceListQueryResult = ApolloReactCommon.QueryResult<GetSeismicEventSourceListQuery, GetSeismicEventSourceListQueryVariables>;
export const GetSeismicEventBySourceDocument = gql`
    query getSeismicEventBySource($sourceId: String!) {
  getEventsBySource(code: $sourceId) {
    verfied
    timestamp
    eventId
    region
    depth
    latitude
    longitude
    magnitude
    distance
    estimation
  }
}
    `;

/**
 * __useGetSeismicEventBySourceQuery__
 *
 * To run a query within a React component, call `useGetSeismicEventBySourceQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSeismicEventBySourceQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSeismicEventBySourceQuery({
 *   variables: {
 *      sourceId: // value for 'sourceId'
 *   },
 * });
 */
export function useGetSeismicEventBySourceQuery(baseOptions: ApolloReactHooks.QueryHookOptions<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables> & ({ variables: GetSeismicEventBySourceQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>(GetSeismicEventBySourceDocument, options);
      }
export function useGetSeismicEventBySourceLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>(GetSeismicEventBySourceDocument, options);
        }
export function useGetSeismicEventBySourceSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>(GetSeismicEventBySourceDocument, options);
        }
export type GetSeismicEventBySourceQueryHookResult = ReturnType<typeof useGetSeismicEventBySourceQuery>;
export type GetSeismicEventBySourceLazyQueryHookResult = ReturnType<typeof useGetSeismicEventBySourceLazyQuery>;
export type GetSeismicEventBySourceSuspenseQueryHookResult = ReturnType<typeof useGetSeismicEventBySourceSuspenseQuery>;
export type GetSeismicEventBySourceQueryResult = ApolloReactCommon.QueryResult<GetSeismicEventBySourceQuery, GetSeismicEventBySourceQueryVariables>;
export const GetSeismicRecordsDocument = gql`
    query getSeismicRecords($startTime: Int64!, $endTime: Int64!) {
  getSeisRecordsByTime(startTime: $startTime, endTime: $endTime) {
    timestamp
    sampleRate
    channelData {
      channelCode
      channelId
      data
    }
  }
}
    `;

/**
 * __useGetSeismicRecordsQuery__
 *
 * To run a query within a React component, call `useGetSeismicRecordsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSeismicRecordsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSeismicRecordsQuery({
 *   variables: {
 *      startTime: // value for 'startTime'
 *      endTime: // value for 'endTime'
 *   },
 * });
 */
export function useGetSeismicRecordsQuery(baseOptions: ApolloReactHooks.QueryHookOptions<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables> & ({ variables: GetSeismicRecordsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>(GetSeismicRecordsDocument, options);
      }
export function useGetSeismicRecordsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>(GetSeismicRecordsDocument, options);
        }
export function useGetSeismicRecordsSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>(GetSeismicRecordsDocument, options);
        }
export type GetSeismicRecordsQueryHookResult = ReturnType<typeof useGetSeismicRecordsQuery>;
export type GetSeismicRecordsLazyQueryHookResult = ReturnType<typeof useGetSeismicRecordsLazyQuery>;
export type GetSeismicRecordsSuspenseQueryHookResult = ReturnType<typeof useGetSeismicRecordsSuspenseQuery>;
export type GetSeismicRecordsQueryResult = ApolloReactCommon.QueryResult<GetSeismicRecordsQuery, GetSeismicRecordsQueryVariables>;
export const GetHomeDataDocument = gql`
    query getHomeData {
  getCurrentTime
  getStationConfig
  getDeviceId
  getDeviceConfig {
    packetInterval
    sampleRate
    channelCodes
    gnssEnabled
    model
    protocol
  }
  getDeviceInfo {
    latitude
    longitude
    elevation
    temperature
  }
  getDeviceStatus {
    startedAt
    updatedAt
    frames
    errors
    messages
  }
  getServiceStatus {
    serviceId
    isRunning
    name
    description
  }
  getSystemStatus {
    cpu
    memory
    disk
    uptime
  }
}
    `;

/**
 * __useGetHomeDataQuery__
 *
 * To run a query within a React component, call `useGetHomeDataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetHomeDataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetHomeDataQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetHomeDataQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetHomeDataQuery, GetHomeDataQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetHomeDataQuery, GetHomeDataQueryVariables>(GetHomeDataDocument, options);
      }
export function useGetHomeDataLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetHomeDataQuery, GetHomeDataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetHomeDataQuery, GetHomeDataQueryVariables>(GetHomeDataDocument, options);
        }
export function useGetHomeDataSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetHomeDataQuery, GetHomeDataQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetHomeDataQuery, GetHomeDataQueryVariables>(GetHomeDataDocument, options);
        }
export type GetHomeDataQueryHookResult = ReturnType<typeof useGetHomeDataQuery>;
export type GetHomeDataLazyQueryHookResult = ReturnType<typeof useGetHomeDataLazyQuery>;
export type GetHomeDataSuspenseQueryHookResult = ReturnType<typeof useGetHomeDataSuspenseQuery>;
export type GetHomeDataQueryResult = ApolloReactCommon.QueryResult<GetHomeDataQuery, GetHomeDataQueryVariables>;
export const CreateUserDocument = gql`
    mutation createUser($username: String!, $password: String!, $admin: Boolean!) {
  createSysUser(username: $username, password: $password, admin: $admin)
}
    `;
export type CreateUserMutationFn = ApolloReactCommon.MutationFunction<CreateUserMutation, CreateUserMutationVariables>;

/**
 * __useCreateUserMutation__
 *
 * To run a mutation, you first call `useCreateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createUserMutation, { data, loading, error }] = useCreateUserMutation({
 *   variables: {
 *      username: // value for 'username'
 *      password: // value for 'password'
 *      admin: // value for 'admin'
 *   },
 * });
 */
export function useCreateUserMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<CreateUserMutation, CreateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<CreateUserMutation, CreateUserMutationVariables>(CreateUserDocument, options);
      }
export type CreateUserMutationHookResult = ReturnType<typeof useCreateUserMutation>;
export type CreateUserMutationResult = ApolloReactCommon.MutationResult<CreateUserMutation>;
export type CreateUserMutationOptions = ApolloReactCommon.BaseMutationOptions<CreateUserMutation, CreateUserMutationVariables>;
export const ExportGlobalConfigDocument = gql`
    query exportGlobalConfig {
  exportGlobalConfig
}
    `;

/**
 * __useExportGlobalConfigQuery__
 *
 * To run a query within a React component, call `useExportGlobalConfigQuery` and pass it any options that fit your needs.
 * When your component renders, `useExportGlobalConfigQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useExportGlobalConfigQuery({
 *   variables: {
 *   },
 * });
 */
export function useExportGlobalConfigQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>(ExportGlobalConfigDocument, options);
      }
export function useExportGlobalConfigLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>(ExportGlobalConfigDocument, options);
        }
export function useExportGlobalConfigSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>(ExportGlobalConfigDocument, options);
        }
export type ExportGlobalConfigQueryHookResult = ReturnType<typeof useExportGlobalConfigQuery>;
export type ExportGlobalConfigLazyQueryHookResult = ReturnType<typeof useExportGlobalConfigLazyQuery>;
export type ExportGlobalConfigSuspenseQueryHookResult = ReturnType<typeof useExportGlobalConfigSuspenseQuery>;
export type ExportGlobalConfigQueryResult = ApolloReactCommon.QueryResult<ExportGlobalConfigQuery, ExportGlobalConfigQueryVariables>;
export const GetApplicationLogsDocument = gql`
    query getApplicationLogs {
  getApplicationLogs
}
    `;

/**
 * __useGetApplicationLogsQuery__
 *
 * To run a query within a React component, call `useGetApplicationLogsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetApplicationLogsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetApplicationLogsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetApplicationLogsQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>(GetApplicationLogsDocument, options);
      }
export function useGetApplicationLogsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>(GetApplicationLogsDocument, options);
        }
export function useGetApplicationLogsSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>(GetApplicationLogsDocument, options);
        }
export type GetApplicationLogsQueryHookResult = ReturnType<typeof useGetApplicationLogsQuery>;
export type GetApplicationLogsLazyQueryHookResult = ReturnType<typeof useGetApplicationLogsLazyQuery>;
export type GetApplicationLogsSuspenseQueryHookResult = ReturnType<typeof useGetApplicationLogsSuspenseQuery>;
export type GetApplicationLogsQueryResult = ApolloReactCommon.QueryResult<GetApplicationLogsQuery, GetApplicationLogsQueryVariables>;
export const GetServiceDataDocument = gql`
    query getServiceData {
  getServiceConfigConstraint {
    serviceName
    serviceId
    constraints {
      key
      name
      description
      configType
      isRequired
      currentValue
      options
    }
  }
  getServiceStatus {
    description
    serviceId
    name
    restarts
    startedAt
    stoppedAt
    updatedAt
    isRunning
  }
}
    `;

/**
 * __useGetServiceDataQuery__
 *
 * To run a query within a React component, call `useGetServiceDataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetServiceDataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetServiceDataQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetServiceDataQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetServiceDataQuery, GetServiceDataQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetServiceDataQuery, GetServiceDataQueryVariables>(GetServiceDataDocument, options);
      }
export function useGetServiceDataLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetServiceDataQuery, GetServiceDataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetServiceDataQuery, GetServiceDataQueryVariables>(GetServiceDataDocument, options);
        }
export function useGetServiceDataSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetServiceDataQuery, GetServiceDataQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetServiceDataQuery, GetServiceDataQueryVariables>(GetServiceDataDocument, options);
        }
export type GetServiceDataQueryHookResult = ReturnType<typeof useGetServiceDataQuery>;
export type GetServiceDataLazyQueryHookResult = ReturnType<typeof useGetServiceDataLazyQuery>;
export type GetServiceDataSuspenseQueryHookResult = ReturnType<typeof useGetServiceDataSuspenseQuery>;
export type GetServiceDataQueryResult = ApolloReactCommon.QueryResult<GetServiceDataQuery, GetServiceDataQueryVariables>;
export const GetStationConfigDocument = gql`
    query getStationConfig {
  getStationConfigConstraint {
    namespace
    name
    key
    description
    configType
    isRequired
    currentValue
    options
  }
}
    `;

/**
 * __useGetStationConfigQuery__
 *
 * To run a query within a React component, call `useGetStationConfigQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetStationConfigQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetStationConfigQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetStationConfigQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetStationConfigQuery, GetStationConfigQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetStationConfigQuery, GetStationConfigQueryVariables>(GetStationConfigDocument, options);
      }
export function useGetStationConfigLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetStationConfigQuery, GetStationConfigQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetStationConfigQuery, GetStationConfigQueryVariables>(GetStationConfigDocument, options);
        }
export function useGetStationConfigSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetStationConfigQuery, GetStationConfigQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetStationConfigQuery, GetStationConfigQueryVariables>(GetStationConfigDocument, options);
        }
export type GetStationConfigQueryHookResult = ReturnType<typeof useGetStationConfigQuery>;
export type GetStationConfigLazyQueryHookResult = ReturnType<typeof useGetStationConfigLazyQuery>;
export type GetStationConfigSuspenseQueryHookResult = ReturnType<typeof useGetStationConfigSuspenseQuery>;
export type GetStationConfigQueryResult = ApolloReactCommon.QueryResult<GetStationConfigQuery, GetStationConfigQueryVariables>;
export const GetStationMetadataDocument = gql`
    query getStationMetadata($format: String!) {
  getStationMetadata(format: $format)
}
    `;

/**
 * __useGetStationMetadataQuery__
 *
 * To run a query within a React component, call `useGetStationMetadataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetStationMetadataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetStationMetadataQuery({
 *   variables: {
 *      format: // value for 'format'
 *   },
 * });
 */
export function useGetStationMetadataQuery(baseOptions: ApolloReactHooks.QueryHookOptions<GetStationMetadataQuery, GetStationMetadataQueryVariables> & ({ variables: GetStationMetadataQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetStationMetadataQuery, GetStationMetadataQueryVariables>(GetStationMetadataDocument, options);
      }
export function useGetStationMetadataLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetStationMetadataQuery, GetStationMetadataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetStationMetadataQuery, GetStationMetadataQueryVariables>(GetStationMetadataDocument, options);
        }
export function useGetStationMetadataSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetStationMetadataQuery, GetStationMetadataQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetStationMetadataQuery, GetStationMetadataQueryVariables>(GetStationMetadataDocument, options);
        }
export type GetStationMetadataQueryHookResult = ReturnType<typeof useGetStationMetadataQuery>;
export type GetStationMetadataLazyQueryHookResult = ReturnType<typeof useGetStationMetadataLazyQuery>;
export type GetStationMetadataSuspenseQueryHookResult = ReturnType<typeof useGetStationMetadataSuspenseQuery>;
export type GetStationMetadataQueryResult = ApolloReactCommon.QueryResult<GetStationMetadataQuery, GetStationMetadataQueryVariables>;
export const GetUserListDocument = gql`
    query getUserList {
  getSysUsers {
    userId
    username
    createdAt
    updatedAt
    lastLogin
    userIp
    userAgent
    admin
  }
}
    `;

/**
 * __useGetUserListQuery__
 *
 * To run a query within a React component, call `useGetUserListQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserListQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetUserListQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetUserListQuery, GetUserListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<GetUserListQuery, GetUserListQueryVariables>(GetUserListDocument, options);
      }
export function useGetUserListLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetUserListQuery, GetUserListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<GetUserListQuery, GetUserListQueryVariables>(GetUserListDocument, options);
        }
export function useGetUserListSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<GetUserListQuery, GetUserListQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<GetUserListQuery, GetUserListQueryVariables>(GetUserListDocument, options);
        }
export type GetUserListQueryHookResult = ReturnType<typeof useGetUserListQuery>;
export type GetUserListLazyQueryHookResult = ReturnType<typeof useGetUserListLazyQuery>;
export type GetUserListSuspenseQueryHookResult = ReturnType<typeof useGetUserListSuspenseQuery>;
export type GetUserListQueryResult = ApolloReactCommon.QueryResult<GetUserListQuery, GetUserListQueryVariables>;
export const ImportGlobalConfigDocument = gql`
    mutation importGlobalConfig($data: String!) {
  importGlobalConfig(data: $data)
}
    `;
export type ImportGlobalConfigMutationFn = ApolloReactCommon.MutationFunction<ImportGlobalConfigMutation, ImportGlobalConfigMutationVariables>;

/**
 * __useImportGlobalConfigMutation__
 *
 * To run a mutation, you first call `useImportGlobalConfigMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useImportGlobalConfigMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [importGlobalConfigMutation, { data, loading, error }] = useImportGlobalConfigMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useImportGlobalConfigMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<ImportGlobalConfigMutation, ImportGlobalConfigMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<ImportGlobalConfigMutation, ImportGlobalConfigMutationVariables>(ImportGlobalConfigDocument, options);
      }
export type ImportGlobalConfigMutationHookResult = ReturnType<typeof useImportGlobalConfigMutation>;
export type ImportGlobalConfigMutationResult = ApolloReactCommon.MutationResult<ImportGlobalConfigMutation>;
export type ImportGlobalConfigMutationOptions = ApolloReactCommon.BaseMutationOptions<ImportGlobalConfigMutation, ImportGlobalConfigMutationVariables>;
export const IsCurrentUserAdminDocument = gql`
    query isCurrentUserAdmin {
  getCurrentUser {
    admin
  }
}
    `;

/**
 * __useIsCurrentUserAdminQuery__
 *
 * To run a query within a React component, call `useIsCurrentUserAdminQuery` and pass it any options that fit your needs.
 * When your component renders, `useIsCurrentUserAdminQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useIsCurrentUserAdminQuery({
 *   variables: {
 *   },
 * });
 */
export function useIsCurrentUserAdminQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useQuery<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>(IsCurrentUserAdminDocument, options);
      }
export function useIsCurrentUserAdminLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useLazyQuery<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>(IsCurrentUserAdminDocument, options);
        }
export function useIsCurrentUserAdminSuspenseQuery(baseOptions?: ApolloReactHooks.SkipToken | ApolloReactHooks.SuspenseQueryHookOptions<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>) {
          const options = baseOptions === ApolloReactHooks.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return ApolloReactHooks.useSuspenseQuery<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>(IsCurrentUserAdminDocument, options);
        }
export type IsCurrentUserAdminQueryHookResult = ReturnType<typeof useIsCurrentUserAdminQuery>;
export type IsCurrentUserAdminLazyQueryHookResult = ReturnType<typeof useIsCurrentUserAdminLazyQuery>;
export type IsCurrentUserAdminSuspenseQueryHookResult = ReturnType<typeof useIsCurrentUserAdminSuspenseQuery>;
export type IsCurrentUserAdminQueryResult = ApolloReactCommon.QueryResult<IsCurrentUserAdminQuery, IsCurrentUserAdminQueryVariables>;
export const PurgeSeisRecordsDocument = gql`
    mutation purgeSeisRecords {
  purgeSeisRecords
}
    `;
export type PurgeSeisRecordsMutationFn = ApolloReactCommon.MutationFunction<PurgeSeisRecordsMutation, PurgeSeisRecordsMutationVariables>;

/**
 * __usePurgeSeisRecordsMutation__
 *
 * To run a mutation, you first call `usePurgeSeisRecordsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePurgeSeisRecordsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [purgeSeisRecordsMutation, { data, loading, error }] = usePurgeSeisRecordsMutation({
 *   variables: {
 *   },
 * });
 */
export function usePurgeSeisRecordsMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PurgeSeisRecordsMutation, PurgeSeisRecordsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<PurgeSeisRecordsMutation, PurgeSeisRecordsMutationVariables>(PurgeSeisRecordsDocument, options);
      }
export type PurgeSeisRecordsMutationHookResult = ReturnType<typeof usePurgeSeisRecordsMutation>;
export type PurgeSeisRecordsMutationResult = ApolloReactCommon.MutationResult<PurgeSeisRecordsMutation>;
export type PurgeSeisRecordsMutationOptions = ApolloReactCommon.BaseMutationOptions<PurgeSeisRecordsMutation, PurgeSeisRecordsMutationVariables>;
export const PurgeMiniSeedFilesDocument = gql`
    mutation purgeMiniSeedFiles {
  purgeMiniSeedFiles
}
    `;
export type PurgeMiniSeedFilesMutationFn = ApolloReactCommon.MutationFunction<PurgeMiniSeedFilesMutation, PurgeMiniSeedFilesMutationVariables>;

/**
 * __usePurgeMiniSeedFilesMutation__
 *
 * To run a mutation, you first call `usePurgeMiniSeedFilesMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePurgeMiniSeedFilesMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [purgeMiniSeedFilesMutation, { data, loading, error }] = usePurgeMiniSeedFilesMutation({
 *   variables: {
 *   },
 * });
 */
export function usePurgeMiniSeedFilesMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PurgeMiniSeedFilesMutation, PurgeMiniSeedFilesMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<PurgeMiniSeedFilesMutation, PurgeMiniSeedFilesMutationVariables>(PurgeMiniSeedFilesDocument, options);
      }
export type PurgeMiniSeedFilesMutationHookResult = ReturnType<typeof usePurgeMiniSeedFilesMutation>;
export type PurgeMiniSeedFilesMutationResult = ApolloReactCommon.MutationResult<PurgeMiniSeedFilesMutation>;
export type PurgeMiniSeedFilesMutationOptions = ApolloReactCommon.BaseMutationOptions<PurgeMiniSeedFilesMutation, PurgeMiniSeedFilesMutationVariables>;
export const PurgeHelicorderFilesDocument = gql`
    mutation purgeHelicorderFiles {
  purgeHelicorderFiles
}
    `;
export type PurgeHelicorderFilesMutationFn = ApolloReactCommon.MutationFunction<PurgeHelicorderFilesMutation, PurgeHelicorderFilesMutationVariables>;

/**
 * __usePurgeHelicorderFilesMutation__
 *
 * To run a mutation, you first call `usePurgeHelicorderFilesMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePurgeHelicorderFilesMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [purgeHelicorderFilesMutation, { data, loading, error }] = usePurgeHelicorderFilesMutation({
 *   variables: {
 *   },
 * });
 */
export function usePurgeHelicorderFilesMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PurgeHelicorderFilesMutation, PurgeHelicorderFilesMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<PurgeHelicorderFilesMutation, PurgeHelicorderFilesMutationVariables>(PurgeHelicorderFilesDocument, options);
      }
export type PurgeHelicorderFilesMutationHookResult = ReturnType<typeof usePurgeHelicorderFilesMutation>;
export type PurgeHelicorderFilesMutationResult = ApolloReactCommon.MutationResult<PurgeHelicorderFilesMutation>;
export type PurgeHelicorderFilesMutationOptions = ApolloReactCommon.BaseMutationOptions<PurgeHelicorderFilesMutation, PurgeHelicorderFilesMutationVariables>;
export const RestoreStationConfigDocument = gql`
    mutation restoreStationConfig {
  restoreStationConfig
}
    `;
export type RestoreStationConfigMutationFn = ApolloReactCommon.MutationFunction<RestoreStationConfigMutation, RestoreStationConfigMutationVariables>;

/**
 * __useRestoreStationConfigMutation__
 *
 * To run a mutation, you first call `useRestoreStationConfigMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestoreStationConfigMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restoreStationConfigMutation, { data, loading, error }] = useRestoreStationConfigMutation({
 *   variables: {
 *   },
 * });
 */
export function useRestoreStationConfigMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<RestoreStationConfigMutation, RestoreStationConfigMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<RestoreStationConfigMutation, RestoreStationConfigMutationVariables>(RestoreStationConfigDocument, options);
      }
export type RestoreStationConfigMutationHookResult = ReturnType<typeof useRestoreStationConfigMutation>;
export type RestoreStationConfigMutationResult = ApolloReactCommon.MutationResult<RestoreStationConfigMutation>;
export type RestoreStationConfigMutationOptions = ApolloReactCommon.BaseMutationOptions<RestoreStationConfigMutation, RestoreStationConfigMutationVariables>;
export const RestoreServiceConfigDocument = gql`
    mutation restoreServiceConfig($serviceId: String) {
  restoreServiceConfig(serviceId: $serviceId)
}
    `;
export type RestoreServiceConfigMutationFn = ApolloReactCommon.MutationFunction<RestoreServiceConfigMutation, RestoreServiceConfigMutationVariables>;

/**
 * __useRestoreServiceConfigMutation__
 *
 * To run a mutation, you first call `useRestoreServiceConfigMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestoreServiceConfigMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restoreServiceConfigMutation, { data, loading, error }] = useRestoreServiceConfigMutation({
 *   variables: {
 *      serviceId: // value for 'serviceId'
 *   },
 * });
 */
export function useRestoreServiceConfigMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<RestoreServiceConfigMutation, RestoreServiceConfigMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<RestoreServiceConfigMutation, RestoreServiceConfigMutationVariables>(RestoreServiceConfigDocument, options);
      }
export type RestoreServiceConfigMutationHookResult = ReturnType<typeof useRestoreServiceConfigMutation>;
export type RestoreServiceConfigMutationResult = ApolloReactCommon.MutationResult<RestoreServiceConfigMutation>;
export type RestoreServiceConfigMutationOptions = ApolloReactCommon.BaseMutationOptions<RestoreServiceConfigMutation, RestoreServiceConfigMutationVariables>;
export const RemoveUserDocument = gql`
    mutation removeUser($userId: String!) {
  removeSysUser(userId: $userId)
}
    `;
export type RemoveUserMutationFn = ApolloReactCommon.MutationFunction<RemoveUserMutation, RemoveUserMutationVariables>;

/**
 * __useRemoveUserMutation__
 *
 * To run a mutation, you first call `useRemoveUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeUserMutation, { data, loading, error }] = useRemoveUserMutation({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useRemoveUserMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<RemoveUserMutation, RemoveUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<RemoveUserMutation, RemoveUserMutationVariables>(RemoveUserDocument, options);
      }
export type RemoveUserMutationHookResult = ReturnType<typeof useRemoveUserMutation>;
export type RemoveUserMutationResult = ApolloReactCommon.MutationResult<RemoveUserMutation>;
export type RemoveUserMutationOptions = ApolloReactCommon.BaseMutationOptions<RemoveUserMutation, RemoveUserMutationVariables>;
export const RestartServiceDocument = gql`
    mutation restartService($serviceId: String!) {
  restartService(serviceId: $serviceId)
}
    `;
export type RestartServiceMutationFn = ApolloReactCommon.MutationFunction<RestartServiceMutation, RestartServiceMutationVariables>;

/**
 * __useRestartServiceMutation__
 *
 * To run a mutation, you first call `useRestartServiceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestartServiceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restartServiceMutation, { data, loading, error }] = useRestartServiceMutation({
 *   variables: {
 *      serviceId: // value for 'serviceId'
 *   },
 * });
 */
export function useRestartServiceMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<RestartServiceMutation, RestartServiceMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<RestartServiceMutation, RestartServiceMutationVariables>(RestartServiceDocument, options);
      }
export type RestartServiceMutationHookResult = ReturnType<typeof useRestartServiceMutation>;
export type RestartServiceMutationResult = ApolloReactCommon.MutationResult<RestartServiceMutation>;
export type RestartServiceMutationOptions = ApolloReactCommon.BaseMutationOptions<RestartServiceMutation, RestartServiceMutationVariables>;
export const StartServiceDocument = gql`
    mutation startService($serviceId: String!) {
  startService(serviceId: $serviceId)
}
    `;
export type StartServiceMutationFn = ApolloReactCommon.MutationFunction<StartServiceMutation, StartServiceMutationVariables>;

/**
 * __useStartServiceMutation__
 *
 * To run a mutation, you first call `useStartServiceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useStartServiceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [startServiceMutation, { data, loading, error }] = useStartServiceMutation({
 *   variables: {
 *      serviceId: // value for 'serviceId'
 *   },
 * });
 */
export function useStartServiceMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<StartServiceMutation, StartServiceMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<StartServiceMutation, StartServiceMutationVariables>(StartServiceDocument, options);
      }
export type StartServiceMutationHookResult = ReturnType<typeof useStartServiceMutation>;
export type StartServiceMutationResult = ApolloReactCommon.MutationResult<StartServiceMutation>;
export type StartServiceMutationOptions = ApolloReactCommon.BaseMutationOptions<StartServiceMutation, StartServiceMutationVariables>;
export const StopServiceDocument = gql`
    mutation stopService($serviceId: String!) {
  stopService(serviceId: $serviceId)
}
    `;
export type StopServiceMutationFn = ApolloReactCommon.MutationFunction<StopServiceMutation, StopServiceMutationVariables>;

/**
 * __useStopServiceMutation__
 *
 * To run a mutation, you first call `useStopServiceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useStopServiceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [stopServiceMutation, { data, loading, error }] = useStopServiceMutation({
 *   variables: {
 *      serviceId: // value for 'serviceId'
 *   },
 * });
 */
export function useStopServiceMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<StopServiceMutation, StopServiceMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<StopServiceMutation, StopServiceMutationVariables>(StopServiceDocument, options);
      }
export type StopServiceMutationHookResult = ReturnType<typeof useStopServiceMutation>;
export type StopServiceMutationResult = ApolloReactCommon.MutationResult<StopServiceMutation>;
export type StopServiceMutationOptions = ApolloReactCommon.BaseMutationOptions<StopServiceMutation, StopServiceMutationVariables>;
export const UpdateServiceConfigDocument = gql`
    mutation updateServiceConfig($serviceId: String!, $key: String!, $val: Any!) {
  updateServiceConfig(serviceId: $serviceId, key: $key, val: $val)
}
    `;
export type UpdateServiceConfigMutationFn = ApolloReactCommon.MutationFunction<UpdateServiceConfigMutation, UpdateServiceConfigMutationVariables>;

/**
 * __useUpdateServiceConfigMutation__
 *
 * To run a mutation, you first call `useUpdateServiceConfigMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateServiceConfigMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateServiceConfigMutation, { data, loading, error }] = useUpdateServiceConfigMutation({
 *   variables: {
 *      serviceId: // value for 'serviceId'
 *      key: // value for 'key'
 *      val: // value for 'val'
 *   },
 * });
 */
export function useUpdateServiceConfigMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<UpdateServiceConfigMutation, UpdateServiceConfigMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<UpdateServiceConfigMutation, UpdateServiceConfigMutationVariables>(UpdateServiceConfigDocument, options);
      }
export type UpdateServiceConfigMutationHookResult = ReturnType<typeof useUpdateServiceConfigMutation>;
export type UpdateServiceConfigMutationResult = ApolloReactCommon.MutationResult<UpdateServiceConfigMutation>;
export type UpdateServiceConfigMutationOptions = ApolloReactCommon.BaseMutationOptions<UpdateServiceConfigMutation, UpdateServiceConfigMutationVariables>;
export const UpdateStationConfigDocument = gql`
    mutation updateStationConfig($key: String!, $value: Any!) {
  updateStationConfig(key: $key, value: $value)
}
    `;
export type UpdateStationConfigMutationFn = ApolloReactCommon.MutationFunction<UpdateStationConfigMutation, UpdateStationConfigMutationVariables>;

/**
 * __useUpdateStationConfigMutation__
 *
 * To run a mutation, you first call `useUpdateStationConfigMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateStationConfigMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateStationConfigMutation, { data, loading, error }] = useUpdateStationConfigMutation({
 *   variables: {
 *      key: // value for 'key'
 *      value: // value for 'value'
 *   },
 * });
 */
export function useUpdateStationConfigMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<UpdateStationConfigMutation, UpdateStationConfigMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<UpdateStationConfigMutation, UpdateStationConfigMutationVariables>(UpdateStationConfigDocument, options);
      }
export type UpdateStationConfigMutationHookResult = ReturnType<typeof useUpdateStationConfigMutation>;
export type UpdateStationConfigMutationResult = ApolloReactCommon.MutationResult<UpdateStationConfigMutation>;
export type UpdateStationConfigMutationOptions = ApolloReactCommon.BaseMutationOptions<UpdateStationConfigMutation, UpdateStationConfigMutationVariables>;
export const UpdateUserDocument = gql`
    mutation updateUser($userId: String!, $username: String!, $password: String, $admin: Boolean!) {
  updateSysUser(
    userId: $userId
    username: $username
    password: $password
    admin: $admin
  )
}
    `;
export type UpdateUserMutationFn = ApolloReactCommon.MutationFunction<UpdateUserMutation, UpdateUserMutationVariables>;

/**
 * __useUpdateUserMutation__
 *
 * To run a mutation, you first call `useUpdateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateUserMutation, { data, loading, error }] = useUpdateUserMutation({
 *   variables: {
 *      userId: // value for 'userId'
 *      username: // value for 'username'
 *      password: // value for 'password'
 *      admin: // value for 'admin'
 *   },
 * });
 */
export function useUpdateUserMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<UpdateUserMutation, UpdateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return ApolloReactHooks.useMutation<UpdateUserMutation, UpdateUserMutationVariables>(UpdateUserDocument, options);
      }
export type UpdateUserMutationHookResult = ReturnType<typeof useUpdateUserMutation>;
export type UpdateUserMutationResult = ApolloReactCommon.MutationResult<UpdateUserMutation>;
export type UpdateUserMutationOptions = ApolloReactCommon.BaseMutationOptions<UpdateUserMutation, UpdateUserMutationVariables>;