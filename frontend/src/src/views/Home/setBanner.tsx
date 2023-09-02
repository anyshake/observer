import { BannerProps } from "../../components/Banner";
import { ApiResponse } from "../../helpers/request/restfulApiByTag";

const setBanner = (res?: ApiResponse): BannerProps => {
    const { error, message } = res || {};
    const { uuid, station, uptime, os } = res?.data || {};

    let label = "连接失败";
    let text = "无法连接到服务器，请尝试刷新页面或更换网络";
    if (!error) {
        label = `连接成功：${station}`;
        text = `${message}\n
                服务器在线时长 ${uptime} 秒\n
                服务器采用架构 ${os.arch}/${os.os}\n
                UUID：${uuid}\n`;
    }

    return {
        type: error ? "error" : "success",
        label,
        text,
    };
};

export default setBanner;
