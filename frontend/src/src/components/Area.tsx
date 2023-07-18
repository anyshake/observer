import { Component, ReactNode } from "react";

export interface AreaProps {
    readonly label: string;
    readonly text: string;
    readonly children?: ReactNode;
}

export default class Area extends Component<AreaProps> {
    render() {
        const { children, label, text } = this.props;
        return (
            <div className="mb-4 bg-white flex flex-col rounded-xl text-gray-700 shadow-lg">
                <div className="mx-4 rounded-lg overflow-hidden shadow-lg">
                    {children}
                </div>

                <div className="p-4">
                    <h6 className="text-md font-bold text-gray-800">{label}</h6>
                    <span className="text-md">
                        {text.split("\n").map((item, index) => (
                            <p key={index}>{item}</p>
                        ))}
                    </span>
                </div>
            </div>
        );
    }
}
