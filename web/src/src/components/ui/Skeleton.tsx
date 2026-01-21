import { useEffect, useState } from 'react';

interface ISkeleton {
    readonly height?: number;
}

export const Skeleton = ({ height }: ISkeleton) => {
    const [skeletonRows, setSkeletonRows] = useState(2);

    useEffect(() => {
        const rows = Math.floor(0.6 * ((height ?? window.innerHeight) / 100));
        setSkeletonRows(rows > 0 ? rows : 2);
    }, [height]);

    return (
        <div className="my-auto w-full animate-pulse space-y-6 overflow-x-hidden p-8">
            {[...new Array(skeletonRows)].map((_, index) => (
                <div key={index} className="space-y-4">
                    <div className="bg-base-300 h-4 w-48 rounded-full"></div>

                    <div className="flex flex-col space-y-3">
                        <div className="bg-base-300 h-3 w-full rounded-full"></div>
                        <div className="bg-base-300 h-3 w-11/12 rounded-full"></div>
                        <div className="bg-base-300 h-3 w-10/12 rounded-full"></div>
                        <div className="bg-base-300 h-3 w-9/12 rounded-full"></div>
                        <div className="bg-base-300 h-3 w-8/12 rounded-full"></div>
                    </div>

                    {index < skeletonRows - 1 && <div className="divider my-2 opacity-20"></div>}
                </div>
            ))}
        </div>
    );
};
