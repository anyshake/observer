/**
 * @date 2023-03-21
 * @author Seunghun Lee - 李承訓
 * @name arrAverage - 数组平均值
 * @description 计算数组的平均值
 * @param {Array} arr - 待计算数组
 * @param {Number} round - 保留小数位数
 * @returns {Number} 平均值
 */
const arrAverage = (arr, round) => {
    if (!arr) {
        return 0;
    }

    const result = arr.reduce((a, b) => a + b) / arr.length;
    return round ? parseFloat(result.toFixed(round)) : result;
};

export default arrAverage;
