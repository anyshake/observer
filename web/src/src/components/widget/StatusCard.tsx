import { ReactNode } from 'react';

import { Card } from './Card';

interface IStatusCard {
    readonly iconPath: string;
    readonly title: string;
    readonly fields: Array<{ label: string; value: string | number }>;
    readonly children?: ReactNode | ReactNode[];
}

export const StatusCard = ({ title, iconPath, fields, children }: IStatusCard) => {
    return (
        <Card title={title} iconPath={iconPath}>
            {fields.map(({ label, value }, index) => (
                <div className="mt-2 flex justify-between" key={`${index}-${label}`}>
                    <span className="pr-4 font-medium whitespace-nowrap">{label}</span>
                    <div className="scrollbar-hide max-w-[180px] overflow-x-scroll">
                        <span className="bg-base-200 rounded px-2 font-mono text-sm whitespace-nowrap text-[#9221ed]">
                            {value}
                        </span>
                    </div>
                </div>
            ))}
            {children && <div className="mt-4">{children}</div>}
        </Card>
    );
};
