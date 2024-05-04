import { apiConfig } from "../../config/api";

export type SocketUpdates = typeof apiConfig.endpoints.socket.model.response.common;

export const getSocketUpdates = async (
	res: SocketUpdates,
	...fn: ((res: SocketUpdates) => void)[]
) => {
	fn.forEach((f) => {
		f(res);
	});
};
