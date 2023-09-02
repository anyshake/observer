import { ADC } from "../../config/adc";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";

const setADC = (res: ApiResponse): ADC => {
    const { resolution, fullscale } = res.data.adc;
    return { resolution, fullscale };
};

export default setADC;
