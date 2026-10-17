import Icon from '@mdi/react';
import { ReactNode } from 'react';

interface ICard {
    readonly iconPath: string;
    readonly title: string;
    readonly className?: string;
    readonly children: ReactNode | ReactNode[];
}

export const Card = ({ title, className, iconPath, children }: ICard) => {
    return (
        <div className="card app-panel app-card">
            <h2 className="app-section-title">
                <Icon className="app-icon" path={iconPath} size={1} />
                <span>{title}</span>
            </h2>
            <hr className="app-card-divider" />
            <div className={className}>{children}</div>
        </div>
    );
};
