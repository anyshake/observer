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
        <ul
            className={`list bg-base-100 rounded-box overflow-y-scroll border border-dashed border-gray-300 ${className}`}
        >
            {data.map(({ primary, secondary, id }, index) => (
                <li
                    key={`${index}-${id}`}
                    onClick={() => onClick?.(id)}
                    className={`border-base-200 border-b last:border-b-0 ${onClick ? 'cursor-pointer' : ''}`}
                >
                    <div className="hover:bg-base-200 flex items-center justify-between px-4 py-3">
                        <div className="flex items-center space-x-3 pr-6">{primary}</div>
                        <div>{secondary}</div>
                    </div>
                </li>
            ))}
        </ul>
    );
};
