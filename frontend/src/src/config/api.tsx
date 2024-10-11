import authRequestModel0 from "../models/request/auth/0.json";
import helicorderRequestModel0 from "../models/request/helicorder/0.json";
import helicorderRequestModel1 from "../models/request/helicorder/1.json";
import historyRequestModel0 from "../models/request/history/0.json";
import inventoryRequestModel0 from "../models/request/inventory/0.json";
import miniseedRequestModel0 from "../models/request/miniseed/0.json";
import miniseedRequestModel1 from "../models/request/miniseed/1.json";
import socketRequestModel0 from "../models/request/socket/0.json";
import stationRequestModel0 from "../models/request/station/0.json";
import traceRequestModel0 from "../models/request/trace/0.json";
import userRequestModel0 from "../models/request/user/0.json";
import authCommonResponseModel0 from "../models/response/common/auth/0.json";
import authCommonResponseModel1 from "../models/response/common/auth/1.json";
import authCommonResponseModel2 from "../models/response/common/auth/2.json";
import helicorderCommonResponseModel0 from "../models/response/common/helicorder/0.json";
import historyCommonResponseModel0 from "../models/response/common/history/0.json";
import inventoryCommonResponseModel0 from "../models/response/common/inventory/0.json";
import miniseedCommonResponseModel0 from "../models/response/common/miniseed/0.json";
import socketCommonResponseModel0 from "../models/response/common/socket/0.json";
import stationCommonResponseModel0 from "../models/response/common/station/0.json";
import traceCommonResponseModel0 from "../models/response/common/trace/0.json";
import traceCommonResponseModel1 from "../models/response/common/trace/1.json";
import userCommonResponseModel0 from "../models/response/common/user/0.json";
import userCommonResponseModel1 from "../models/response/common/user/1.json";
import userCommonResponseModel2 from "../models/response/common/user/2.json";
import userCommonResponseModel3 from "../models/response/common/user/3.json";
import authErrorResponseModel from "../models/response/error/auth.json";
import helicorderErrorResponseModel from "../models/response/error/helicorder.json";
import historyErrorResponseModel from "../models/response/error/history.json";
import inventoryErrorResponseModel from "../models/response/error/inventory.json";
import miniseedErrorResponseModel from "../models/response/error/miniseed.json";
import stationErrorResponseModel from "../models/response/error/station.json";
import traceErrorResponseModel from "../models/response/error/trace.json";
import userErrorResponseModel from "../models/response/error/user.json";

export interface Endpoint<APIRequest, APICommonResponse, APIErrorResponse = null> {
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

export { stationCommonResponseModel0, stationErrorResponseModel, stationRequestModel0 };

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
			error: stationErrorResponseModel
		}
	}
};

export {
	authCommonResponseModel0,
	authCommonResponseModel1,
	authCommonResponseModel2,
	authRequestModel0
};

const auth: Endpoint<
	typeof authRequestModel0,
	| typeof authCommonResponseModel0
	| typeof authCommonResponseModel1
	| typeof authCommonResponseModel2,
	typeof authErrorResponseModel
> = {
	path: "/api/v1/auth",
	method: "post",
	type: "http",
	model: {
		request: { ...authRequestModel0 },
		response: {
			common: {
				...authCommonResponseModel0,
				...authCommonResponseModel1,
				...authCommonResponseModel2
			},
			error: authErrorResponseModel
		}
	}
};

export { historyCommonResponseModel0, historyErrorResponseModel, historyRequestModel0 };

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
			error: historyErrorResponseModel
		}
	}
};

export {
	traceCommonResponseModel0,
	traceCommonResponseModel1,
	traceErrorResponseModel,
	traceRequestModel0
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
				...traceCommonResponseModel1
			},
			error: traceErrorResponseModel
		}
	}
};

export {
	miniseedCommonResponseModel0,
	miniseedErrorResponseModel,
	miniseedRequestModel0,
	miniseedRequestModel1
};

const miniseed: Endpoint<
	typeof miniseedRequestModel0 | typeof miniseedRequestModel1,
	typeof miniseedCommonResponseModel0,
	typeof miniseedErrorResponseModel
> = {
	path: "/api/v1/miniseed",
	method: "post",
	type: "http",
	model: {
		request: { ...miniseedRequestModel0, ...miniseedRequestModel1 },
		response: {
			common: { ...miniseedCommonResponseModel0 },
			error: miniseedErrorResponseModel
		}
	}
};

export {
	helicorderCommonResponseModel0,
	helicorderErrorResponseModel,
	helicorderRequestModel0,
	helicorderRequestModel1
};

const helicorder: Endpoint<
	typeof helicorderRequestModel0 | typeof helicorderRequestModel1,
	typeof helicorderCommonResponseModel0,
	typeof helicorderErrorResponseModel
> = {
	path: "/api/v1/helicorder",
	method: "post",
	type: "http",
	model: {
		request: { ...helicorderRequestModel0, ...helicorderRequestModel1 },
		response: {
			common: { ...helicorderCommonResponseModel0 },
			error: helicorderErrorResponseModel
		}
	}
};

export { socketCommonResponseModel0, socketRequestModel0 };

const socket: Endpoint<typeof socketRequestModel0, typeof socketCommonResponseModel0> = {
	path: "/api/v1/socket",
	type: "socket",
	model: {
		request: { ...socketRequestModel0 },
		response: { common: { ...socketCommonResponseModel0 } }
	}
};

export { inventoryCommonResponseModel0, inventoryErrorResponseModel, inventoryRequestModel0 };

const inventory: Endpoint<
	typeof inventoryRequestModel0,
	typeof inventoryCommonResponseModel0,
	typeof inventoryErrorResponseModel
> = {
	path: "/api/v1/inventory",
	method: "get",
	type: "http",
	model: {
		request: { ...inventoryRequestModel0 },
		response: {
			common: { ...inventoryCommonResponseModel0 },
			error: inventoryErrorResponseModel
		}
	}
};

export {
	userCommonResponseModel0,
	userCommonResponseModel1,
	userCommonResponseModel2,
	userCommonResponseModel3,
	userRequestModel0
};

const user: Endpoint<
	typeof userRequestModel0,
	| typeof userCommonResponseModel0
	| typeof userCommonResponseModel1
	| typeof userCommonResponseModel2
	| typeof userCommonResponseModel3,
	typeof userErrorResponseModel
> = {
	path: "/api/v1/user",
	method: "post",
	type: "http",
	model: {
		request: { ...userRequestModel0 },
		response: {
			common: {
				...userCommonResponseModel0,
				...userCommonResponseModel1,
				...userCommonResponseModel2,
				...userCommonResponseModel3
			},
			error: userErrorResponseModel
		}
	}
};

export const apiConfig = {
	backend:
		process.env.NODE_ENV === "production"
			? `${window.location.host}`
			: `${process.env.REACT_APP_BACKEND}`,
	endpoints: {
		auth,
		station,
		history,
		trace,
		miniseed,
		helicorder,
		socket,
		inventory,
		user
	}
};
