import stationRequestModel0 from "../models/request/station/0.json";
import stationCommonResponseModel0 from "../models/response/common/station/0.json";
import stationErrorResponseModel from "../models/response/error/station.json";
import historyRequestModel0 from "../models/request/history/0.json";
import historyCommonResponseModel0 from "../models/response/common/history/0.json";
import historyErrorResponseModel from "../models/response/error/history.json";
import traceRequestModel0 from "../models/request/trace/0.json";
import traceCommonResponseModel0 from "../models/response/common/trace/0.json";
import traceCommonResponseModel1 from "../models/response/common/trace/1.json";
import traceErrorResponseModel from "../models/response/error/trace.json";
import mseedRequestModel0 from "../models/request/mseed/0.json";
import mseedRequestModel1 from "../models/request/mseed/1.json";
import mseedCommonResponseModel0 from "../models/response/common/mseed/0.json";
import mseedErrorResponseModel from "../models/response/error/mseed.json";
import socketRequestModel0 from "../models/request/socket/0.json";
import socketCommonResponseModel0 from "../models/response/common/socket/0.json";

export interface Endpoint<
    APIRequest,
    APICommonResponse,
    APIErrorResponse = null
> {
    readonly model: {
        request?: APIRequest;
        response: {
            common: APICommonResponse;
            error?: APIErrorResponse;
        };
    };
    readonly path: string;
    readonly type: "http" | "socket";
    readonly method?: "get" | "post";
}

export {
    stationRequestModel0,
    stationCommonResponseModel0,
    stationErrorResponseModel,
};

const station: Endpoint<
    typeof stationRequestModel0,
    typeof stationCommonResponseModel0,
    typeof stationErrorResponseModel
> = {
    path: "/api/v1/station",
    method: "get",
    type: "http",
    model: {
        request: { ...stationRequestModel0 },
        response: {
            common: { ...stationCommonResponseModel0 },
            error: stationErrorResponseModel,
        },
    },
};

export {
    historyRequestModel0,
    historyCommonResponseModel0,
    historyErrorResponseModel,
};

const history: Endpoint<
    typeof historyRequestModel0,
    typeof historyCommonResponseModel0,
    typeof historyErrorResponseModel
> = {
    path: "/api/v1/history",
    method: "post",
    type: "http",
    model: {
        request: { ...historyRequestModel0 },
        response: {
            common: { ...historyCommonResponseModel0 },
            error: historyErrorResponseModel,
        },
    },
};

export {
    traceRequestModel0,
    traceCommonResponseModel0,
    traceCommonResponseModel1,
    traceErrorResponseModel,
};

const trace: Endpoint<
    typeof traceRequestModel0,
    typeof traceCommonResponseModel0 | typeof traceCommonResponseModel1,
    typeof traceErrorResponseModel
> = {
    path: "/api/v1/trace",
    method: "post",
    type: "http",
    model: {
        request: { ...traceRequestModel0 },
        response: {
            common: {
                ...traceCommonResponseModel0,
                ...traceCommonResponseModel1,
            },
            error: traceErrorResponseModel,
        },
    },
};

export {
    mseedRequestModel0,
    mseedRequestModel1,
    mseedCommonResponseModel0,
    mseedErrorResponseModel,
};

const mseed: Endpoint<
    typeof mseedRequestModel0 | typeof mseedRequestModel1,
    typeof mseedCommonResponseModel0,
    typeof mseedErrorResponseModel
> = {
    path: "/api/v1/mseed",
    method: "post",
    type: "http",
    model: {
        request: { ...mseedRequestModel0, ...mseedRequestModel1 },
        response: {
            common: { ...mseedCommonResponseModel0 },
            error: mseedErrorResponseModel,
        },
    },
};

export { socketRequestModel0, socketCommonResponseModel0 };

const socket: Endpoint<
    typeof socketRequestModel0,
    typeof socketCommonResponseModel0
> = {
    path: "/api/v1/socket",
    type: "socket",
    model: {
        request: { ...socketRequestModel0 },
        response: { common: { ...socketCommonResponseModel0 } },
    },
};

export const apiConfig = {
    backend:
        process.env.NODE_ENV === "production"
            ? `${window.location.host}`
            : `${process.env.REACT_APP_BACKEND}`,
    endpoints: {
        station,
        history,
        trace,
        mseed,
        socket,
    },
};
