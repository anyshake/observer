/**
 * @date 2023-04-11
 * @author Seunghun Lee - 李承訓
 * @name arrSmoothing - 数组滤波
 * @description 使用中值法对数组进行滤波
 * @param {Array} arr - 待滤波数组
 * @param {Number} window - 滤波窗口大小
 * @returns {Array} 滤波后的数组
 */
const arrSmoothing = (arr, window) => {
    const result = [];
    for (let i = 0; i < arr.length; i++) {
        const start = i - window;
        const end = i + window;
        const temp = [];
        for (let j = start; j <= end; j++) {
            if (j >= 0 && j < arr.length) {
                temp.push(arr[j]);
            }
        }
        temp.sort((a, b) => a - b);
        result.push(temp[Math.floor(temp.length / 2)]);
    }

    return result;
};

export default arrSmoothing;
