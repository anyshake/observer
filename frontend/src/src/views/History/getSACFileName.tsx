import { getDayOfYear } from "../../helpers/utils/getDayOfYear";
import { Station } from "../../stores/station";

export const getSACFileName = (
    startTimeMS: number,
    channelName: string,
    { network, station, location }: Station
) => {
    const startTimeObj = new Date(startTimeMS);
    return `${startTimeObj.getUTCFullYear()}.${getDayOfYear(startTimeObj)
        .toString()
        .padStart(3, "0")}.${startTimeObj
        .getUTCHours()
        .toString()
        .padStart(2, "0")}.${startTimeObj
        .getUTCMinutes()
        .toString()
        .padStart(2, "0")}.${startTimeObj
        .getUTCSeconds()
        .toString()
        .padStart(2, "0")}.${startTimeObj
        .getUTCMilliseconds()
        .toString()
        .padStart(4, "0")}.${network.slice(0, 2)}.${station.slice(
        0,
        5
    )}.${location.slice(0, 2)}.${channelName}.D.sac`;
};
