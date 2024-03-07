import { apiConfig } from "../../config/api";
import { requestRestApi } from "../../helpers/request/requestRestApi";

export type ExportsUpdates =
    | typeof apiConfig.endpoints.mseed.model.response.common
    | typeof apiConfig.endpoints.mseed.model.response.error;

export const getExportsUpdates = async (
    ...fn: ((res: ExportsUpdates) => void)[]
) => {
    const { endpoints, backend } = apiConfig;
    const data = await requestRestApi({
        backend,
        timeout: 30,
        endpoint: endpoints.mseed,
        payload: { action: "show", name: "" },
    });
    fn.forEach((f) => {
        f(data);
    });
};
