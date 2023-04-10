/**
 * @date 2023-03-21
 * @author Seunghun Lee - 李承訓
 * @name arrFilter - 数组过滤器
 * @description 重新封装数组过滤器，支援对象数组过滤
 * @param {Array} arr - 对象数组
 * @param {String} filter - 过滤器
 * @param {String} keyword - 过滤关键字
 * @param {Function} callback - 回呼函数
 * @returns {VoidFunction} 无返回值
 */
const arrFilter = ({ arr, filter, keyword, callback }) => {
    const index = arr.findIndex((item) => item[filter] === keyword);
    if (index !== -1) {
        const item = arr[index];
        callback(item, index, arr);
    }
};

export default arrFilter;
