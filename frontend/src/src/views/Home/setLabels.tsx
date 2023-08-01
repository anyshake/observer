import { LabelProps } from "../../components/Label";
import { ApiResponse } from "../../helpers/requestByTag";
import setObject from "../../helpers/setObject";

const setLabels = (obj: LabelProps[], res: ApiResponse): LabelProps[] => {
    const tags = [
        "messages",
        "errors",
        "pushed",
        "failures",
        "queued",
        "offset",
    ];
    for (let i of tags) {
        const { status } = res.data;
        setObject(obj, `[tag:${i}]>value`, status[i]);
    }

    return obj;
};

export default setLabels;
