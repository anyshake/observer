import axios, {
    AxiosError,
    AxiosResponse,
    InternalAxiosRequestConfig,
} from "axios";
import API_CONFIG from "../config/api";
import getBackend from "./getBackend";
import getApiUrl from "./getApiUrl";
import fileDownload from "js-file-download";

export interface RequestByTag {
    readonly tag: string;
    readonly blob?: boolean;
    readonly timeout?: number;
    readonly filename?: string;
    readonly header?: { [key: string]: string };
    readonly body?: { [key: string]: string | number | undefined };
}

export interface ApiResponse {
    readonly time: string;
    readonly status: number;
    readonly error: boolean;
    readonly path: string;
    readonly message: string;
    readonly data: any;
}

const requestByTag = async ({
    tag,
    header,
    body,
    blob,
    filename,
    timeout = 10000,
}: RequestByTag): Promise<ApiResponse> => {
    const _axios = axios.create({
        timeout: timeout,
    });

    _axios.interceptors.request.use((config: InternalAxiosRequestConfig) => {
        if (!blob) {
            config.headers["Accept"] = "application/json";
        }

        return config;
    });
    _axios.interceptors.response.use(
        (res: AxiosResponse | any) => res,
        (err: AxiosError) => Promise.reject(err)
    );

    const url = getApiUrl(tag);
    const method = API_CONFIG.find((config) => config.tag === tag)?.method;

    try {
        const type = API_CONFIG.find((config) => config.tag === tag)?.type;
        if (type === "ws") {
            throw new Error("websocket protocol is not supported");
        }

        const backend = `${window.location.protocol}${getBackend()}`;
        const { data, headers } = await _axios.request({
            responseType: blob ? "blob" : "json",
            url: `${backend}${url}`,
            headers: header,
            method: method,
            data: body,
        });

        if (blob) {
            const { "content-disposition": contentDisposition } = headers;
            if (contentDisposition) {
                const filename = contentDisposition
                    .split(";")
                    .find((item: string) => item.includes("filename="))
                    ?.split("=")[1];
                if (filename) {
                    fileDownload(data, filename);
                } else {
                    fileDownload(data, "stream");
                }
            } else if (filename) {
                fileDownload(data, filename);
            } else {
                fileDownload(data, "stream");
            }

            const time = new Date().toISOString();
            return {
                time: time,
                path: url,
                data: null,
                error: false,
                status: 200,
                message: "Returned data is a blob",
            };
        }

        return data as ApiResponse;
    } catch (err: unknown) {
        const time = new Date().toISOString();
        const { message, status } = err as AxiosError;

        return {
            path: url,
            data: null,
            error: true,
            status: status || 500,
            message,
            time,
        };
    }
};

export default requestByTag;
