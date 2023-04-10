/**
 * @date 2023-04-10
 * @author Seunghun Lee - 李承訓
 * @name createSocket - WebSocket 创建器
 * @description 建立同后端通信的 WebSocket 连接
 * @param {String} url - 请求地址
 * @param {String} type - 数据类型
 * @param {Object} onMessageCallback - 连接讯息回呼函数
 * @param {Object} onCloseCallback - 连接关闭回呼函数
 * @param {Object} onErrorCallback - 连接出错回呼函数
 * @returns {WebSocket} Promise - 返回 Promise
 */
const createSocket = ({
    url,
    type,
    onMessageCallback,
    onCloseCallback,
    onErrorCallback,
}) => {
    const conn = new WebSocket(url);
    conn.onmessage = onMessageCallback;
    conn.onclose = onCloseCallback;
    conn.onerror = onErrorCallback;
    conn.binaryType = type;

    return conn;
};

export default createSocket;
