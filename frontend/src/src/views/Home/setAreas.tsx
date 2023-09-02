import { HomeArea } from "./Index";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import getQueueArray from "../../helpers/array/getQueueArray";
import setObjectByPath from "../../helpers/utils/setObjectByPath";

const setAreas = (
    obj: HomeArea[],
    res: ApiResponse,
    length: number
): HomeArea[] => {
    const tags = ["cpu", "memory"];

    for (let i of tags) {
        const percent = res.data[i].percent as number;
        const newArr = [Date.now(), percent];
        const srcArr = obj.find((item) => item.tag === i)?.chart.series?.data;

        const result = getQueueArray(srcArr, newArr, length);
        setObjectByPath(
            obj,
            `[tag:${i}]>area>text`,
            `当前占用率：${percent.toFixed(2)}%`
        );
        setObjectByPath(obj, `[tag:${i}]>chart>series>data`, result);
    }

    return obj;
};

export default setAreas;
