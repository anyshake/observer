import { Station } from "../../config/station";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";

const setStation = (res: ApiResponse): Station => {
    const { station, network, location } = res.data?.station || {};
    return { station, network, location };
};

export default setStation;
