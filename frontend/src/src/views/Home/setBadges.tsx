import { BadgesProps } from "../../components/Badges";
import { ApiResponse } from "../../helpers/requestByTag";
import setObject from "../../helpers/setObject";

const setCard = (obj: BadgesProps, res: ApiResponse): BadgesProps => {
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
        setObject(obj, `list[tag:${i}]>value`, status[i]);
    }

    return obj;
};

export default setCard;
