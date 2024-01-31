import axios, {
    AxiosError,
    AxiosProgressEvent,
    AxiosResponse,
    CancelTokenSource,
    InternalAxiosRequestConfig,
    isCancel,
} from "axios";
import API_CONFIG from "../../config/api";
import getBackend from "../app/getBackend";
import getApiUrl from "../app/getApiUrl";
import GLOBAL_CONFIG from "../../config/global";
import { saveAs } from "file-saver";

export interface RESTfulApiByTag {
    readonly tag: string;
    readonly blob?: boolean;
    readonly timeout?: number;
    readonly filename?: string;
    readonly header?: { [key: string]: string };
    readonly body?: { [key: string]: string | number | undefined };
    readonly onUpload?: (progressEvent: AxiosProgressEvent) => void;
    readonly onDownload?: (progressEvent: AxiosProgressEvent) => void;
    readonly cancelToken?: CancelTokenSource;
}

export interface ApiResponse {
    readonly time: string;
    readonly status: number;
    readonly error: boolean;
    readonly path: string;
    readonly message: string;
    readonly data: any;
}

const restfulApiByTag = async ({
    tag,
    header,
    body,
    blob,
    filename,
    onUpload,
    onDownload,
    cancelToken,
    timeout = GLOBAL_CONFIG.app_settings.timeout,
}: RESTfulApiByTag): Promise<ApiResponse> => {
    const _axios = axios.create({
        timeout: timeout * 1000,
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
            data: body,
            method: method,
            headers: header,
            url: `${backend}${url}`,
            onUploadProgress: onUpload,
            onDownloadProgress: onDownload,
            cancelToken: cancelToken?.token,
            responseType: blob ? "blob" : "json",
        });

        if (blob) {
            const { "content-disposition": contentDisposition } = headers;
            if (contentDisposition) {
                filename = contentDisposition
                    .split(";")
                    .find((item: string) => item.includes("filename="))
                    ?.split("=")[1];
            }
            saveAs(data, filename ?? "stream");

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
            time,
            message,
            path: url,
            data: null,
            error: !isCancel(err),
            status: status || 500,
        };
    }
};

export default restfulApiByTag;
