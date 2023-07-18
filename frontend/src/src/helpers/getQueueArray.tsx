const isDeepArray = (arr: []) => arr.some((item: any) => item instanceof Array);

const getQueueArray = (srcArr: any[], newData: any, length: number): any[] => {
    if (isDeepArray(newData)) {
        for (let i = 0; i < newData.length; i++) {
            srcArr.push(newData[i]);
        }
    } else {
        srcArr.push(newData);
    }

    if (srcArr.length > length) {
        srcArr.splice(0, srcArr.length - length);
    }

    return [...srcArr];
};

export default getQueueArray;
