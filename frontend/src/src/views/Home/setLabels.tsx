import { LabelProps } from "../../components/Label";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import setObjectByPath from "../../helpers/utils/setObjectByPath";

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
        setObjectByPath(obj, `[tag:${i}]>value`, status[i]);
    }

    return obj;
};

export default setLabels;
