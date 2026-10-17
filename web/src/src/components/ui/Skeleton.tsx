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
        <div className="app-skeleton">
            {[...new Array(skeletonRows)].map((_, index) => (
                <div key={index} className="app-stack">
                    <div className="app-skeleton-title"></div>

                    <div className="app-skeleton-lines">
                        <div className="app-skeleton-line-full"></div>
                        <div className="app-skeleton-line-long"></div>
                        <div className="app-skeleton-line-medium"></div>
                        <div className="app-skeleton-line-short"></div>
                        <div className="app-skeleton-line-shortest"></div>
                    </div>

                    {index < skeletonRows - 1 && <div className="app-skeleton-divider"></div>}
                </div>
            ))}
        </div>
    );
};
