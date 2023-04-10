import store from "../redux/store";
import arrFilter from "../utilities/arrFilter";

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name hookAddEventListener - 实现 addEventListener
 * @description 实现自定义 addEventListener，事件注册前先保存函数引用到 React-Redux，方便后续移除
 * @returns {VoidFunction} 无返回值
 */
const hookAddEventListener = () => {
    // 用户需要为 options 参数附上一个唯一 ID，以便后续移除事件
    EventTarget.prototype.appAddEventListener = ({
        type,
        listener,
        options,
    }) => {
        EventTarget.prototype.addEventListener.call(
            this,
            type,
            listener,
            options
        );
        store.dispatch({
            type: "ADD_EVENT_LISTENER",
            payload: {
                listener: listener,
                id: options?.id || "default",
            },
        });
    };
};

/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name hookRemoveEventListener - 拷贝 removeEventListener 原型
 * @description 实现自定义 addEventListener，先移除 React-Redux 保存的函数引用，再移除事件
 * @returns {VoidFunction} 无返回值
 */
const hookRemoveEventListener = () => {
    // 透过唯一 ID，将 React-Redux 保存的函数引用移除
    EventTarget.prototype.appRemoveEventListener = ({ type, id, options }) => {
        arrFilter({
            arr: store.getState().eventListener,
            filter: "id",
            keyword: id,
            callback: (item) => {
                EventTarget.prototype.removeEventListener.call(
                    this,
                    type,
                    item.listener,
                    options
                );
                store.dispatch({
                    type: "REMOVE_EVENT_LISTENER",
                    payload: {
                        id: id,
                    },
                });
            },
        });
    };
};

export { hookAddEventListener, hookRemoveEventListener };
