import { apiConfig } from "../../config/api";
import { requestRestApi } from "../../helpers/request/requestRestApi";

export type StationUpdates =
	| typeof apiConfig.endpoints.station.model.response.common
	| typeof apiConfig.endpoints.station.model.response.error;

export const getStationUpdates = async (...fn: ((res: StationUpdates) => void)[]) => {
	const { endpoints, backend } = apiConfig;
	const data = await requestRestApi({
		backend,
		timeout: 30,
		endpoint: endpoints.station
	});
	fn.forEach((f) => {
		f(data);
	});
};
