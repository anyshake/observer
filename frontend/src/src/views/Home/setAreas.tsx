import { HomeArea } from ".";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import getQueueArray from "../../helpers/array/getQueueArray";
import setObjectByPath from "../../helpers/utils/setObjectByPath";

const setAreas = (
    obj: HomeArea[],
    res: ApiResponse,
    length: number
): HomeArea[] => {
    for (let i of obj) {
        // Get percent by tag
        const percent = res.data[i.tag].percent as number;
        // Create new array, get source array
        const newArr = [Date.now(), percent];
        const srcArr = obj.find((item) => item.tag === i.tag)?.chart.series
            ?.data;
        // Merge new array with source array
        const result = getQueueArray(srcArr, newArr, length);
        setObjectByPath(obj, `[tag:${i.tag}]>area>text`, {
            id: `views.home.areas.${i.tag}.text`,
            format: { usage: percent.toFixed(2) },
        });
        setObjectByPath(obj, `[tag:${i.tag}]>chart>series>data`, result);
    }

    return obj;
};

export default setAreas;
