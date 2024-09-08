import axios, {
	AxiosError,
	AxiosProgressEvent,
	AxiosResponse,
	InternalAxiosRequestConfig
} from "axios";
import { saveAs } from "file-saver";

import { Endpoint } from "../../config/api";
import { getProtocol } from "../app/getProtocol";

interface Options<APIRequest, APICommonResponse, APIErrorResponse> {
	readonly throwError?: boolean;
	readonly backend: string;
	readonly timeout: number;
	readonly payload?: APIRequest;
	readonly header?: Record<string, string>;
	readonly abortController?: AbortController;
	readonly endpoint: Endpoint<APIRequest, APICommonResponse, APIErrorResponse>;
	readonly blobOptions?: {
		readonly fileName: string;
		readonly onDownload?: (progressEvent: AxiosProgressEvent) => void;
	};
}

export const requestRestApi = async <APIRequest, APICommonResponse, APIErrorResponse>({
	header,
	payload,
	backend,
	endpoint,
	throwError,
	blobOptions,
	abortController,
	timeout = 100
}: Options<APIRequest, APICommonResponse, APIErrorResponse>): Promise<
	APICommonResponse | APIErrorResponse
> => {
	const _axios = axios.create({
		timeout: timeout * 1000
	});
	_axios.interceptors.request.use((config: InternalAxiosRequestConfig) => {
		if (!blobOptions) {
			config.headers.Accept = "application/json";
		}
		return config;
	});
	_axios.interceptors.response.use(
		(res: AxiosResponse) => res,
		(err: AxiosError) => Promise.reject(err)
	);

	const { response } = endpoint.model;
	try {
		if (endpoint.type === "socket") {
			throw new Error("websocket protocol is not supported");
		}

		const protocol = getProtocol(true);
		let reqPath = `${protocol}//${backend}${endpoint.path}`;
		const query = new URLSearchParams();
		if (endpoint.method === "get" && !!payload) {
			Object.entries(payload).forEach(([key, value]) => {
				query.set(key, value as string);
			});
			reqPath += `?${query.toString()}`;
		}

		const { data, headers } = await _axios.request({
			url: reqPath,
			headers: header,
			method: endpoint.method,
			signal: abortController?.signal,
			onDownloadProgress: blobOptions?.onDownload,
			responseType: blobOptions ? "blob" : "json",
			data: endpoint.method === "post" ? payload : {}
		});
		if (blobOptions) {
			const { "content-disposition": contentDisposition } = headers;
			let fileName = blobOptions.fileName;
			if (contentDisposition) {
				fileName = contentDisposition
					.split(";")
					.find((item: string) => item.includes("filename="))
					?.split("=")[1];
			}

			saveAs(data, !fileName.length ? "stream" : fileName);
			return response.common;
		}

		return { ...response.common, ...data };
	} catch (e) {
		const result = response.error ?? response.common;
		return throwError ? Promise.reject(e) : result;
	}
};
