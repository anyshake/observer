import { ApolloLink, FetchResult, Observable, Operation } from '@apollo/client';
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { saveAs } from 'file-saver';

import { useCredentialStore } from '../../stores/credential';

interface IResponse<T = unknown> {
    code: number;
    error: boolean;
    message: string;
    data: T | null;
}

class AxiosWrapper {
    private instance: AxiosInstance;

    constructor() {
        this.instance = axios.create();
        this.instance.interceptors.request.use((config) => {
            const { token } = useCredentialStore.getState().getCredential();
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
            return config;
        });
        this.instance.interceptors.response.use(
            (response) => response,
            (error) => {
                if (error.response.status === 401) {
                    useCredentialStore.getState().clearCredential();
                }
                return Promise.reject(error);
            }
        );
    }

    createGraphQlLink(url: string) {
        return new ApolloLink(({ query, variables, operationName }: Operation) => {
            return new Observable<FetchResult>((observer) => {
                this.instance
                    .post(url, { variables, operationName, query: query.loc?.source.body })
                    .then((result) => {
                        observer.next(result.data);
                        observer.complete();
                    })
                    .catch((error) => {
                        observer.error(error);
                    });
            });
        });
    }

    async request<T = unknown>(
        config: AxiosRequestConfig & { ignoreErrors?: boolean }
    ): Promise<IResponse<T>> {
        const { ignoreErrors = false, ...axiosConfig } = config;

        try {
            const response = await this.instance.request<IResponse<T>>(axiosConfig);
            return response.data;
        } catch (error) {
            if (!ignoreErrors) {
                throw error;
            }

            if (axios.isAxiosError(error) && error.response) {
                return {
                    data: null,
                    code: error.response.status ?? 0,
                    error: true,
                    message: error.response.data?.message ?? error.message
                };
            }

            return {
                data: null,
                code: -1,
                error: true,
                message: 'no response received from server'
            };
        }
    }

    async getBlob(config: AxiosRequestConfig & { ignoreErrors?: boolean }): Promise<Blob | null> {
        const { ignoreErrors = false, ...axiosConfig } = config;

        try {
            const response = await this.instance.request<Blob>({
                ...axiosConfig,
                responseType: 'blob'
            });

            return response.data;
        } catch (error) {
            if (!ignoreErrors) {
                throw axios.isAxiosError(error) && error.response
                    ? JSON.parse(await error.response.data.text()).message
                    : error;
            }

            if (axios.isAxiosError(error) && error.response) {
                return null;
            }

            return null;
        }
    }

    async saveAs(
        config: AxiosRequestConfig & { ignoreErrors?: boolean }
    ): Promise<IResponse<void> | void> {
        const { ignoreErrors = false, ...axiosConfig } = config;

        try {
            const response = await this.instance.request<Blob>({
                ...axiosConfig,
                responseType: 'blob'
            });

            const { 'content-disposition': disposition } = response.headers;
            let filename = 'octet-stream';

            if (disposition) {
                filename = disposition
                    .split(';')
                    .find((item: string) => item.includes('filename='))
                    ?.split('=')[1];
            }

            const blob = new Blob([response.data]);
            saveAs(blob, filename);
        } catch (error) {
            if (!ignoreErrors) {
                throw axios.isAxiosError(error) && error.response
                    ? JSON.parse(await error.response.data.text()).message
                    : error;
            }

            if (axios.isAxiosError(error) && error.response) {
                return {
                    data: null,
                    code: error.response.status ?? 0,
                    error: true,
                    message: error.response.data?.message ?? error.message
                };
            }

            return {
                data: null,
                code: -1,
                error: true,
                message: 'no response received from server'
            };
        }
    }
}

export const ApiClient = new AxiosWrapper();
