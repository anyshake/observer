import axios from "axios";

/**
 * @date 2023-04-10
 * @author Seunghun Lee - 李承訓
 * @name createRequest -  请求创建器
 * @description 创建请求到后端 API
 * @param {String} url - 请求地址
 * @param {String} method - 请求方法
 * @param {Object} headers - 请求头
 * @param {Object} data - 请求数据
 * @returns {Promise} Promise - 返回 Promise
 */
const createRequest = ({ url, method, headers, data }) => {
    return axios.request({
        headers: headers,
        method: method,
        data: data,
        url: url,
    });
};

export default createRequest;
