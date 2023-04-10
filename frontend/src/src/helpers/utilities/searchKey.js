/**
 * @date 2023-04-11
 * @author Seunghun Lee - 李承訓
 * @name searchKey - 搜寻物件中的 key
 * @description 搜寻物件中的 key，返回 key 对应的值
 * @param {Object} obj
 * @param {String} key
 * @returns {String} key 对应的值
 */
const searchKey = (obj, key) => {
    for (let k in obj) {
        if (k === key) {
            return obj[k];
        } else if (obj[k] instanceof Object) {
            const result = searchKey(obj[k], key);
            if (result) {
                return result;
            }
        }
    }
};

export default searchKey;
