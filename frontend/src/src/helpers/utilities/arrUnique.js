/**
 * @date 2023-03-30
 * @author Seunghun Lee - 李承訓
 * @name arrUnique - 数组去重
 * @description 数组去重，支援物件数组
 * @param {Array} arr - 数组
 * @returns {Array} 去重后的数组
 */
const arrUnique = (arr) => {
    const newArr = [];
    const obj = {};

    arr.forEach((item) => {
        if (!obj[item]) {
            newArr.push(item);
            obj[item] = true;
        }
    });

    return newArr;
};

export default arrUnique;
