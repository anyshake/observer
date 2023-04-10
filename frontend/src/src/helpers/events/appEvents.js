/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name registerEvent - 注册事件
 * @description 向 window 注册事件
 * @param {Array} eventArray - 事件列表
 * @param {Function} onEventCallback - 事件触发回呼函数
 * @returns {VoidFunction} 无返回值
 */
const registerEvents = ({ eventArray, onEventCallback }) => {
    eventArray.map((item) =>
        window.appAddEventListener({
            type: item.trigger,
            listener: onEventCallback,
            options: {
                id: item.id,
            },
        })
    );
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name removeEvent - 移除事件
 * @description 从 window 移除事件
 * @param {Array} eventArray - 事件列表
 * @returns {VoidFunction} 无返回值
 */
const removeEvents = (eventArray) => {
    eventArray.map((item) =>
        window.appRemoveEventListener({
            type: item.trigger,
            id: item.id,
        })
    );
};

export { registerEvents, removeEvents };
