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
                <div className="app-status-row" key={`${index}-${label}`}>
                    <span className="app-status-label">{label}</span>
                    <div className="app-status-value-scroll">
                        <span className="app-status-value">{value}</span>
                    </div>
                </div>
            ))}
            {children && <div className="app-section-footer">{children}</div>}
        </Card>
    );
};
