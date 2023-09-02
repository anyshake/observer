import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import { HomeMap } from "./Index";

const setMap = (obj: HomeMap, res: ApiResponse): HomeMap => {
    const { location } = res.data;
    return {
        ...obj,
        area: {
            ...obj.area,
            text: `测站经度：${location.longitude} °\n
                测站纬度：${location.latitude} °\n
                测站海拔：${location.altitude} m`,
        },
        instance: {
            ...obj.instance,
            center: [location.latitude, location.longitude],
            marker: [location.latitude, location.longitude],
        },
    };
};

export default setMap;
