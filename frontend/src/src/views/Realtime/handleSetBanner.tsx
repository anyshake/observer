import { Dispatch, SetStateAction } from "react";
import { BannerProps } from "../../components/Banner";
import { SocketUpdates } from "./getSocketUpdates";
import { getTimeString } from "../../helpers/utils/getTimeString";
import { globalConfig } from "../../config/global";
import store from "../../config/store";

export const handleSetBanner = (
    res: SocketUpdates,
    stateFn: Dispatch<
        SetStateAction<BannerProps & { values?: Record<string, string> }>
    >
) => {
    if (!res.ts) {
        return;
    }

    // Get scale name and its instance
    const { scale: scaleId } = store.getState().scale;
    const scaleName =
        globalConfig.scales
            .find((s) => s.property().value === scaleId)
            ?.property().name ?? "Unknown";
    const time = getTimeString(res.ts);

    // Get sample rate in average
    const channels = Object.values(res).filter((v) =>
        Array.isArray(v)
    ) as number[][];
    const sampleRate = (
        channels.reduce((acc, cur) => acc + cur.length, 0) / channels.length
    ).toFixed(0);

    stateFn((prev) => ({
        ...prev,
        type: "success",
        title: "views.realtime.banner.success.label",
        content: "views.realtime.banner.success.text",
        values: { sampleRate, time, scale: scaleName },
    }));
};
