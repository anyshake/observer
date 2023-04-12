/**
 * @date 2023-03-30
 * @author Seunghun Lee - 李承訓
 * @name getTime - 取得当前时间
 * @description 取得本地 YYYY-MM-DD hh:mm:ss 格式时间
 * @param {DateConstructor} date - 时间对象
 * @returns {String} 时间字符串
 */
const getTime = (date) => {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, "0");
    const day = date.getDate().toString().padStart(2, "0");
    const hour = date.getHours().toString().padStart(2, "0");
    const minute = date.getMinutes().toString().padStart(2, "0");
    const second = date.getSeconds().toString().padStart(2, "0");

    return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
};

export default getTime;
