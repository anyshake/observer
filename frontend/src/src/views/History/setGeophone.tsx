import { Geophone } from "../../config/geophone";
import { ApiResponse } from "../../helpers/requestByTag";

const setGeophone = (res: ApiResponse): Geophone => {
    const { ehz, ehe, ehn } = res.data.geophone;
    return { ehz, ehe, ehn };
};

export default setGeophone;
