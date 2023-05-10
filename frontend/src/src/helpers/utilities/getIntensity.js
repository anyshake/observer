/**
 * @date 2023-05-08
 * @author Seunghun Lee - 李承訓
 * @name getIntensity - 取得震度等级
 * @description 传入加速度，返回震度等级（0-7）
 * @param {Number} acceleration - 加速度
 * @returns {Number} 震度等级
 */
const getIntensity = (acceleration) => {
    switch (true) {
        case acceleration < 0.1:
            return 0;
        case acceleration < 0.2:
            return 1;
        case acceleration < 0.4:
            return 2;
        case acceleration < 0.8:
            return 3;
        case acceleration < 1.6:
            return 4;
        case acceleration < 3.2:
            return 5;
        case acceleration < 6.4:
            return 6;
        default:
            return 7;
    }
};

export default getIntensity;
