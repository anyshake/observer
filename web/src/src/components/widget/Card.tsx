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
        <div className="card bg-base-100 p-5 text-gray-700 shadow-md h-full">
            <h2 className="flex items-center space-x-2 text-lg font-bold">
                <Icon className="flex-shrink-0" path={iconPath} size={1} />
                <span>{title}</span>
            </h2>
            <hr className="my-4 text-gray-200" />
            <div className={className}>{children}</div>
        </div>
    );
};
