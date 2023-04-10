/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name redirectRouter - 重定向路由
 * @description 应用内无刷新重定向路由
 * @param {String} dest - 目标路由位址
 * @param {Boolean} replace - 是否替换当前路由
 * @returns {VoidFunction} 无返回值
 */
const redirectRouter = ({ dest, replace }) => {
    if (replace) {
        window.location.hash.includes("#/")
            ? window.location.replace("#" + dest)
            : window.location.replace(dest);
    } else {
        window.location.hash.includes("#/")
            ? (window.location.href = "#" + dest)
            : (window.location.href = dest);
    }
};

export default redirectRouter;
