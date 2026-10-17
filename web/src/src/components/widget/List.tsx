import { ReactNode } from 'react';

interface IList {
    readonly data: Array<{
        readonly id: string;
        readonly primary: ReactNode | ReactNode[];
        readonly secondary: ReactNode | ReactNode[];
    }>;
    readonly onClick?: (id: string) => void;
    readonly className?: string;
}

export const List = ({ data, onClick, className }: IList) => {
    return (
        <ul className={`list rounded-box app-list ${className}`}>
            {data.map(({ primary, secondary, id }, index) => (
                <li
                    key={`${index}-${id}`}
                    onClick={() => onClick?.(id)}
                    className={`app-list-item ${onClick ? 'app-pointer' : ''}`}
                >
                    <div className="app-list-row">
                        <div className="app-list-primary">{primary}</div>
                        <div>{secondary}</div>
                    </div>
                </li>
            ))}
        </ul>
    );
};
