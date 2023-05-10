/**
 * @date 2023-05-08
 * @author Seunghun Lee - 李承訓
 * @name getIntensity - 取得震度等级
 * @description 传入加速度，返回 JMA 震度等级（0-7）
 * @param {Number} acceleration - 加速度
 * @returns {String} 震度等级
 */
const getIntensity = (acceleration) => {
    const intensity = Math.round(
        Math.round(2 * Math.log(Math.abs(acceleration)) + 0.94),
        2
    );

    switch (true) {
        case intensity < 0.5:
            return "0";
        case intensity < 1.5:
            return "1";
        case intensity < 2.5:
            return "2";
        case intensity < 3.5:
            return "3";
        case intensity < 4.5:
            return "4";
        case intensity < 5.0:
            return "5-";
        case intensity < 5.5:
            return "5+";
        case intensity < 6.0:
            return "6-";
        case intensity < 6.5:
            return "6+";
        default:
            return "7";
    }
};

export default getIntensity;
