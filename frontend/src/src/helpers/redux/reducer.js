/**
 * @date 2023-03-20
 * @author Seunghun Lee - 李承訓
 * @name reducer - Redux 触发函数
 * @description 根据 action.type 触发不同的动作
 * @param {Object} state - Redux 状态物件
 * @param {Object} action - Redux 触发动作
 * @returns {Object} 返回新的状态物件
 */
const reducer = (state, action) => {
    if (!state) {
        return {
            bluethDevice: null,
            focusDuration: 0,
            eventListener: [],
            userPreference: null,
        };
    }

    switch (action.type) {
        case "SET_BLUETH_DEVICE":
            return {
                ...state,
                bluethDevice: action.payload,
            };
        case "SET_FOCUS_DURATION":
            return {
                ...state,
                focusDuration: action.payload,
            };
        case "ADD_EVENT_LISTENER":
            if (
                !state.eventListener.some(
                    (item) => item.id === action.payload.id
                ) &&
                action.payload.id !== "default"
            ) {
                state.eventListener.push(action.payload);
            }
            return {
                ...state,
            };
        case "REMOVE_EVENT_LISTENER":
            state.eventListener.splice(
                state.eventListener.findIndex(
                    (item) => item.id === action.payload.id
                ),
                1
            );
            return {
                ...state,
            };
        case "SET_USER_PREFERENCE":
            return {
                ...state,
                userPreference: action.payload,
            };
        default:
            return state;
    }
};

export default reducer;
