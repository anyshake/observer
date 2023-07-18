import { Component } from "react";

export interface BadgeProps {
    readonly tag?: string;
    readonly label: string;
    readonly value: string;
    readonly unit: string;
    readonly icon: string;
}

export interface BadgesProps {
    readonly list: BadgeProps[];
}

export default class Badges extends Component<BadgesProps> {
    render() {
        const { list } = this.props;
        return (
            <div className="flex flex-wrap">
                {list.map(({ icon, label, value, unit }, index) => (
                    <div key={index} className="w-full md:w-1/2 lg:w-1/3 p-2">
                        <div className="flex flex-row bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-md p-4 shadow-xl">
                            <img
                                className="bg-white p-2 rounded-md w-8 h-8 md:w-12 md:h-12 self-center"
                                src={icon}
                                alt=""
                            />

                            <div className="flex flex-col flex-grow ml-5 text-white">
                                <div className="text-sm whitespace-nowrap">
                                    {label}
                                </div>

                                <div className="text-md">
                                    {value}
                                    <span className="text-lg px-2">{unit}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        );
    }
}
