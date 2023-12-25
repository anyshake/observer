import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import { HomeMap } from ".";

const setMap = (obj: HomeMap, res: ApiResponse): HomeMap => {
    const { location } = res.data;
    const { longitude, latitude, elevation } = location;
    return {
        ...obj,
        area: {
            ...obj.area,
            text: {
                id: "views.home.map.area.text",
                format: {
                    elevation: elevation.toFixed(2),
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
