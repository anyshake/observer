import { useState } from "react";
import folderIcon from "../assets/icons/folder-open-regular.svg";

export interface TableColumn {
    key: string;
    label: string;
}

export type TableData = {
    [key: TableColumn["key"]]: string | number;
};

export interface TableAction {
    icon: string;
    label: string;
    onClick?: (data: TableData) => void;
}

export interface TableProps {
    data?: TableData[];
    columns: TableColumn[];
    actions: TableAction[];
    placeholder: string;
    rowsLimit?: number;
    loadMoreText?: string;
}

export const Table = (props: TableProps) => {
    const { columns, actions, data, placeholder, rowsLimit, loadMoreText } =
        props;

    const [rowsLength, setRowsLength] = useState(rowsLimit ?? -1);

    const loadMoreRows = () =>
        setRowsLength((currentLength) => currentLength + (rowsLimit ?? 0));

    return (
        <div className="flex flex-col overflow-x-auto">
            {!!data?.length ? (
                <div className="space-y-2 pb-4">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead>
                            <tr>
                                {columns.map(({ label }) => (
                                    <th
                                        key={label}
                                        scope="col"
                                        className="px-6 py-3 whitespace-nowrap text-left text-xs font-medium text-gray-500"
                                    >
                                        {label}
                                    </th>
                                ))}
                                {actions.map(({ label }) => (
                                    <th
                                        key={label}
                                        scope="col"
                                        className="px-6 py-3 whitespace-nowrap text-left text-xs font-medium text-gray-500"
                                    >
                                        {label}
                                    </th>
                                ))}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200 text-gray-700">
                            {data
                                .slice(
                                    0,
                                    rowsLength === -1 ? data.length : rowsLength
                                )
                                .map((item, index) => (
                                    <tr
                                        key={index}
                                        className="hover:bg-gray-100"
                                    >
                                        {Object.keys(item).map(
                                            (_item, _index) => (
                                                <td
                                                    key={_item}
                                                    className="px-6 py-4 whitespace-nowrap text-sm font-medium"
                                                >
                                                    {columns
                                                        .filter(
                                                            ({ key }) =>
                                                                key ===
                                                                columns[_index]
                                                                    .key
                                                        )
                                                        .map(
                                                            ({ key }) =>
                                                                item[key]
                                                        ) || _item}
                                                </td>
                                            )
                                        )}
                                        {actions.map(
                                            ({ icon, label, onClick }) => (
                                                <td
                                                    key={label}
                                                    className="px-6 py-4 whitespace-nowrap text-sm font-medium"
                                                    onClick={() =>
                                                        onClick && onClick(item)
                                                    }
                                                >
                                                    <img
                                                        className="w-5 h-5 cursor-pointer transition-all duration-200 hover:scale-125"
                                                        src={icon}
                                                        alt={label}
                                                    />
                                                </td>
                                            )
                                        )}
                                    </tr>
                                ))}
                        </tbody>
                    </table>
                    {rowsLength !== -1 && data.length > rowsLength && (
                        <div
                            className="text-sm text-center text-gray-500 cursor-pointer hover:text-gray-700"
                            onClick={loadMoreRows}
                        >
                            {loadMoreText ?? "Load more..."}
                        </div>
                    )}
                </div>
            ) : (
                <div className="flex justify-center items-center h-40 text-gray-500 space-x-2">
                    <img
                        className="size-5 md:size-6 lg:size-8"
                        src={folderIcon}
                        alt=""
                    />
                    <h1 className="text-lg md:text-xl lg:text-2xl font-medium">
                        {placeholder}
                    </h1>
                </div>
            )}
        </div>
    );
};
