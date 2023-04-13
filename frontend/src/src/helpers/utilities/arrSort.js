/**
 * @date 2023-04-12
 * @author Seunghun Lee - 李承訓
 * @name arrSort - 数组排序
 * @description 数组排序
 * @param {Array} arr - 待计算数组
 * @param {String} key - 排序的键
 * @param {String} order - 排序方式
 * @returns {Number} 排序后的数组
 */
const arrSort = (arr, key, order = "asc") => {
    if (!arr) {
        return [];
    }

    return arr.sort((a, b) => {
        if (order === "asc") {
            return a[key] - b[key];
        }

        return b[key] - a[key];
    });
};

export default arrSort;
