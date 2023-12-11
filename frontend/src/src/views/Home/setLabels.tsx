import { LabelProps } from "../../components/Label";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";
import setObjectByPath from "../../helpers/utils/setObjectByPath";

const setLabels = (obj: LabelProps[], res: ApiResponse): LabelProps[] => {
    for (let i of obj) {
        const { status } = res.data;
        setObjectByPath(obj, `[tag:${i.tag}]>value`, status[i.tag || ""]);
    }

    return obj;
};

export default setLabels;
