import { HomeArea } from "./Index";
import { ApiResponse } from "../../helpers/requestByTag";
import getQueueArray from "../../helpers/getQueueArray";
import setObject from "../../helpers/setObject";

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
        setObject(
            obj,
            `[tag:${i}]>area>text`,
            `当前占用率：${percent.toFixed(2)}%`
        );
        setObject(obj, `[tag:${i}]>chart>series>data`, result);
    }

    return obj;
};

export default setAreas;
