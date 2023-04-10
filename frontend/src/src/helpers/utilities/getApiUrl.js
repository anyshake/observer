/**
 * @date 2023-04-11
 * @author Seunghun Lee - 李承訓
 * @name getApiUrl - 取得 API URL
 * @description 传入 API 的 host、port、api、version、tls、type，返回 API 的 URL
 * @param {String} host - API 的 host
 * @param {String} port - API 的 port
 * @param {String} api - API 的 api 名称
 * @param {String} version - API 的版本号
 * @param {Boolean} tls - API 的 tls 状态
 * @param {String} type - API 的 type，http 或 websocket
 * @returns {String} API 的 URL
 */
const getApiUrl = ({ host, port, api, version, tls, type }) => {
    const baseUrl = `${host}:${port}/api/${version}/${api}`;
    switch (type) {
        case `http`:
            if (tls) {
                return `https://${baseUrl}`;
            }
            return `http://${baseUrl}`;

        case `websocket`:
            if (tls) {
                return `wss://${baseUrl}`;
            }
            return `ws://${baseUrl}`;

        default:
            return null;
    }
};

export default getApiUrl;
