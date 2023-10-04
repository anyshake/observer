import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import { HomeMap } from "./Home";

const setMap = (obj: HomeMap, res: ApiResponse): HomeMap => {
    const { location } = res.data;
    const { longitude, latitude, altitude } = location;
    return {
        ...obj,
        area: {
            ...obj.area,
            text: {
                id: "views.home.map.area.text",
                format: {
                    altitude: altitude.toFixed(2),
                    latitude: latitude.toFixed(2),
                    longitude: longitude.toFixed(2),
                },
            },
        },
        instance: {
            ...obj.instance,
            center: [latitude, longitude],
            marker: [latitude, longitude],
        },
    };
};

export default setMap;
