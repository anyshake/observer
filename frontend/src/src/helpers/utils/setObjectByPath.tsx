const setObjectByPath = (obj: any, path: string, value: any): any => {
    let current = obj;

    const keys = path.split(">");
    for (let i = 0; i < keys.length - 1; i++) {
        const key = keys[i];

        try {
            if (key.includes("[") || key.includes("]")) {
                const arrKey: string = key.match(/^(.*?)\[/)?.[1] || "";
                const arrPath: string = key.match(/\[(.*?)\]/)?.[1] || ":";
                if (!arrPath.length) {
                    throw new Error("invalid path given");
                }

                const [tag, target] = arrPath.split(":");
                if (!arrKey.length) {
                    current = current.find((item: any) => item[tag] === target);
                } else {
                    current = current[arrKey].find(
                        (item: any) => item[tag] === target
                    );
                }
            } else {
                current = current[key];
            }
        } catch {
            return obj;
        }
    }

    const lastKey = keys[keys.length - 1];
    current[lastKey] = value;
    return obj;
};

export default setObjectByPath;
