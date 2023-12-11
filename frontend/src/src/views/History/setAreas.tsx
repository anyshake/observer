import getSortedArray from "../../helpers/array/getSortedArray";
import setObjectByPath from "../../helpers/utils/setObjectByPath";
import { HistoryArea } from ".";

const setAreas = (obj: HistoryArea[], data: any): HistoryArea[] => {
    // Sort data by timestamp
    const sortedData: any = getSortedArray(data, "ts", "number", "asc");
    for (let i of obj) {
        // Previous timestamp
        let prevTs = 0;
        // Store chart data format
        const resultArr = [];
        // Go through each sorted sample
        for (let j of sortedData) {
            // Get data sample rate
            const channelData = j[i.tag];
            const sampleRate = channelData.length;
            // Get time difference and time span
            const timeDiff = prevTs !== 0 ? prevTs - j.ts : 1000;
            const timeSpan = timeDiff / sampleRate;
            // Append result array
            for (let k = 0; k < channelData.length; k++) {
                resultArr.push([j.ts + k * timeSpan, channelData[k]]);
            }

            // Set previous timestamp
            prevTs = j.ts;
        }

        // Set velocity chart data
        setObjectByPath(obj, `[tag:${i.tag}]>chart>series>data`, resultArr);
    }

    return obj;
};

export default setAreas;
