import { Component } from "react";

export interface LabelProps {
    readonly className?: string;
    readonly tag?: string;
    readonly label: string;
    readonly value: string;
    readonly unit: string;
    readonly icon?: string;
    readonly color?: boolean;
}

export default class Label extends Component<LabelProps> {
    render() {
        const { className, icon, label, value, unit, color } = this.props;
        return (
            <div className={`w-full p-2 ${className}`}>
                <div
                    className={`flex flex-row bg-gradient-to-r rounded-md p-4 shadow-xl ${
                        color
                            ? `from-indigo-500 via-purple-500 to-pink-500`
                            : `bg-slate-50 hover:bg-slate-100`
                    }`}
                >
                    {icon && (
                        <img
                            className="bg-white p-2 rounded-md w-8 h-8 md:w-12 md:h-12 self-center"
                            src={icon}
                            alt=""
                        />
                    )}

                    <div
                        className={`flex flex-col flex-grow ${icon && `ml-5`}`}
                    >
                        <div
                            className={`text-sm whitespace-nowrap ${
                                color ? `text-gray-50` : `text-gray-600`
                            }`}
                        >
                            {label}
                        </div>
                        <div
                            className={`text-md font-medium ${
                                color ? `text-gray-100` : `text-gray-800`
                            }`}
                        >
                            {value}
                            <span className="text-lg px-2">{unit}</span>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
