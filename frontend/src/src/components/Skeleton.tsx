import { useEffect, useState } from "react";

export const Skeleton = () => {
    const DEFAULT_SKELETON_ROWS = 2;
    const [skeletonRows, setSkeletonRows] = useState(DEFAULT_SKELETON_ROWS);

    useEffect(() => {
        const rows = Math.floor(0.55 * (window.innerHeight / 100));
        setSkeletonRows(rows > 0 ? rows : DEFAULT_SKELETON_ROWS);
    }, []);

    return (
        <div className="p-8 m-auto space-y-3 w-[calc(90%)] animate-pulse animate-duration-700">
            {[...new Array(skeletonRows)].map((_, index) => (
                <div key={index} className="space-y-3">
                    <div className="h-2 bg-gray-300 rounded-full w-32 mb-4" />
                    <div className="h-2 bg-gray-300 rounded-full" />
                    <div className="h-2 bg-gray-300 rounded-full" />
                    <div className="h-2 bg-gray-300 rounded-full" />
                    <div className="h-2 bg-gray-300 rounded-full" />
                    <div className="h-2 bg-gray-300 rounded-full" />
                </div>
            ))}
        </div>
    );
};
