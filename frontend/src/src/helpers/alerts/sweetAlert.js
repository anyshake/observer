import Swal from "sweetalert2";

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name successAlert - 成功提示框
 * @description SweetAlert2 成功提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @returns {Promise} Promise - 返回 Promise
 */
const successAlert = ({ title, html }) => {
    return Swal.fire({
        title: title,
        html: html,
        icon: "success",
        confirmButtonText: "确认",
        allowOutsideClick: false,
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name errorAlert - 错误提示框
 * @description SweetAlert2 错误提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @returns {Promise} Promise - 返回 Promise
 */
const errorAlert = ({ title, html }) => {
    return Swal.fire({
        title: title,
        html: html,
        icon: "error",
        confirmButtonText: "确认",
        allowOutsideClick: false,
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name warningAlert - 警告提示框
 * @description SweetAlert2 警告提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @returns {Promise} Promise - 返回 Promise
 */
const warningAlert = ({ title, html }) => {
    return Swal.fire({
        title: title,
        html: html,
        icon: "warning",
        confirmButtonText: "确认",
        showCancelButton: true,
        cancelButtonText: "取消",
        allowOutsideClick: false,
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name infoAlert - 资讯提示框
 * @description SweetAlert2 资讯提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @returns {Promise} Promise - 返回 Promise
 */
const infoAlert = ({ title, html }) => {
    return Swal.fire({
        title: title,
        html: html,
        icon: "info",
        confirmButtonText: "确认",
        allowOutsideClick: false,
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name confirmAlert - 确认提示框
 * @description SweetAlert2 确认提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @param {String} confirmButtonText - 确认按钮文字
 * @param {String} cancelButtonText - 取消按钮文字
 * @param {Function} callback - 回呼函数
 * @returns {Promise} Promise - 返回 Promise
 */
const confirmAlert = ({
    title,
    html,
    confirmButtonText,
    cancelButtonText,
    callback,
}) => {
    return Swal.fire({
        title: title,
        html: html,
        icon: "warning",
        showCancelButton: cancelButtonText ? true : false,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        allowOutsideClick: false,
        cancelButtonText: cancelButtonText,
        confirmButtonText: confirmButtonText,
    }).then((result) => {
        if (result.value) {
            callback && callback();
        }
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name timerAlert - 计时提示框
 * @description SweetAlert2 计时提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @param {Number} timer - 计时器时长
 * @param {Boolean} loading 显示加载动画
 * @param {Function} callback - 回呼函数
 * @returns {Promise} Promise - 返回 Promise
 */
const timerAlert = async ({ title, html, timer, loading, callback }) => {
    let timerInterval = null;
    return await Swal.fire({
        title: title,
        html: html,
        timer: timer,
        timerProgressBar: true,
        allowOutsideClick: false,
        showConfirmButton: false,
        didOpen: () => {
            loading && Swal.showLoading();
        },
        willClose: () => {
            clearInterval(timerInterval);
        },
    }).then((result) => {
        if (result.dismiss === Swal.DismissReason.timer) {
            callback && callback();
        }
        return result;
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name toastAlert - 悬浮提示框
 * @description SweetAlert2 悬浮提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @param {String} icon - 提示框图标
 * @param {Number} timer - 计时器时长
 * @returns {Promise} Promise - 返回 Promise
 */
const toastAlert = ({ title, html, icon, timer }) => {
    const Toast = Swal.mixin({
        toast: true,
        position: "top-end",
        timer: timer,
        timerProgressBar: true,
        showConfirmButton: false,
        didOpen: (toast) => {
            toast.addEventListener("mouseenter", Swal.stopTimer);
            toast.addEventListener("mouseleave", Swal.resumeTimer);
        },
    });
    return Toast.fire({
        icon: icon,
        title: title,
        html: html,
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name selectAlert - 选择提示框
 * @description SweetAlert2 选择提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @param {Object} inputOptions - 提示框选项
 * @param {Function} callback - 回呼函数
 * @returns {Promise} Promise - 返回 Promise
 */
const selectAlert = ({ title, html, inputOptions, callback }) => {
    return Swal.fire({
        title: title,
        html: html,
        input: "select",
        inputOptions: inputOptions,
        inputPlaceholder: "请选择",
        showCancelButton: false,
        allowOutsideClick: false,
        inputValidator: (value) => {
            if (!value) {
                return "请选择有效的选项";
            }
        },
    }).then((result) => {
        if (result.value) {
            callback && callback(result.value);
        }
    });
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name inputAlert - 输入框提示框
 * @description SweetAlert2 输入框提示框
 * @param {String} title - 提示框标题
 * @param {String} html - 提示框内联 HTML 内容
 * @param {String} input - 提示框输入类型
 * @param {Function} callback - 回呼函数
 * @returns {Promise} Promise - 返回 Promise
 */
const inputAlert = ({ title, html, input, callback }) => {
    return Swal.fire({
        title: title,
        html: html,
        input: input,
        inputAttributes: {
            autocapitalize: "off",
        },
        confirmButtonText: "确认",
        showCancelButton: false,
        allowOutsideClick: false,
        inputValidator: (value) => {
            if (!value) {
                return "请输入有效的值";
            }
        },
    }).then((result) => {
        if (result.value) {
            callback && callback(result.value);
        }
    });
};

export {
    successAlert,
    errorAlert,
    warningAlert,
    infoAlert,
    confirmAlert,
    timerAlert,
    toastAlert,
    selectAlert,
    inputAlert,
};
