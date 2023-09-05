import getSortedArray from "../../helpers/array/getSortedArray";
import getCounts from "../../helpers/seismic/getCounts";
import setObjectByPath from "../../helpers/utils/setObjectByPath";
import { HistoryArea } from "./Index";

const setAreas = (obj: HistoryArea[], data: any): HistoryArea[] => {
    const tags = ["ehz", "ehe", "ehn"];

    // Sort data by timestamp
    const sortedData: any = getSortedArray(data, "ts", "asc");
    for (let i of tags) {
        // Previous timestamp
        let prevTs = 0;
        // Store chart data format
        const resultArr = [];
        // Go through each sorted sample
        for (let j of sortedData) {
            // Get data sample rate
            const channelData = j[i];
            const sampleRate = channelData.length;
            // Get time difference and time span
            const timeDiff = prevTs !== 0 ? prevTs - j.ts : 1000;
            const timeSpan = timeDiff / sampleRate;
            // Get counts (offset data to center around 0)
            const counts = getCounts(channelData);
            // Append result array
            for (let k = 0; k < counts.length; k++) {
                resultArr.push([j.ts + k * timeSpan, counts[k]]);
            }

            // Set previous timestamp
            prevTs = j.ts;
        }

        // Set velocity chart data
        setObjectByPath(obj, `[tag:${i}]>chart>series>data`, resultArr);
    }

    return obj;
};

export default setAreas;
