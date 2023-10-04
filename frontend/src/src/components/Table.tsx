import { Component } from "react";
import FolderIcon from "../assets/icons/folder-open-regular.svg";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

export interface TableColumn {
    key: string;
    label: I18nTranslation;
}

export interface TableAction {
    icon: string;
    label: I18nTranslation;
    onClick?: (data: TableData) => void;
}

export type TableData = {
    [key: TableColumn["key"]]: string;
};

export interface TableProps {
    data: TableData[];
    columns: TableColumn[];
    actions: TableAction[];
    placeholder: I18nTranslation;
}

class Table extends Component<TableProps & WithTranslation> {
    render() {
        const { t, columns, actions, data, placeholder } = this.props;
        return (
            <div className="flex flex-col">
                <div className="-m-1.5 overflow-x-auto">
                    <div className="p-1.5 min-w-full inline-block align-middle">
                        <div className="overflow-hidden">
                            {!!data.length ? (
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead>
                                        <tr>
                                            {columns.map(({ label }, index) => (
                                                <th
                                                    key={index}
                                                    scope="col"
                                                    className="px-6 py-3 whitespace-nowrap text-left text-xs font-medium text-gray-500"
                                                >
                                                    {t(label.id, label.format)}
                                                </th>
                                            ))}
                                            {actions.map(({ label }, index) => (
                                                <th
                                                    key={index}
                                                    scope="col"
                                                    className="px-6 py-3 whitespace-nowrap text-left text-xs font-medium text-gray-500"
                                                >
                                                    {t(label.id, label.format)}
                                                </th>
                                            ))}
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200 text-gray-700">
                                        {data.map((item, index) => (
                                            <tr
                                                key={index}
                                                className="hover:bg-gray-100"
                                            >
                                                {Object.keys(item).map(
                                                    (_item, _index) => (
                                                        <td
                                                            key={_index}
                                                            className="px-6 py-4 whitespace-nowrap text-sm font-medium"
                                                        >
                                                            {columns
                                                                .filter(
                                                                    ({ key }) =>
                                                                        key ===
                                                                        columns[
                                                                            _index
                                                                        ].key
                                                                )
                                                                .map(
                                                                    ({ key }) =>
                                                                        item[
                                                                            key
                                                                        ]
                                                                ) || _item}
                                                        </td>
                                                    )
                                                )}
                                                {actions.map(
                                                    (
                                                        {
                                                            icon,
                                                            label,
                                                            onClick,
                                                        },
                                                        _index
                                                    ) => (
                                                        <td
                                                            key={_index}
                                                            className="px-6 py-4 whitespace-nowrap text-sm font-medium"
                                                            onClick={() =>
                                                                onClick &&
                                                                onClick(item)
                                                            }
                                                        >
                                                            <img
                                                                className="w-5 h-5 cursor-pointer transition-all duration-200 hover:scale-125"
                                                                src={icon}
                                                                alt={t(
                                                                    label.id,
                                                                    label.format
                                                                )}
                                                            />
                                                        </td>
                                                    )
                                                )}
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            ) : (
                                <div className="flex justify-center items-center h-40 text-gray-500 space-x-2">
                                    <img
                                        src={FolderIcon}
                                        alt="Folder Icon"
                                        className="w-8 h-8"
                                    />
                                    <h1 className="text-2xl font-medium">
                                        {t(placeholder.id, placeholder.format)}
                                    </h1>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default withTranslation()(Table);
