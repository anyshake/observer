import { Component, ReactNode } from "react";

interface CardProps {
    readonly className?: string;
    readonly label: string;
    readonly sublabel?: string;
    readonly children: ReactNode | ReactNode[];
}

export default class Card extends Component<CardProps> {
    render() {
        const { className, label, sublabel, children } = this.props;
        const childrenArr = Array.isArray(children) ? children : [children];

        return (
            <div className="w-full h-full text-gray-800">
                <div className="flex flex-col shadow-lg rounded-lg">
                    <div className="px-4 py-3 font-bold">
                        <h6 className="text-gray-500 text-xs">{sublabel}</h6>
                        <h2 className="text-xl">{label}</h2>
                    </div>

                    <div className={`p-4 m-2 flex flex-col justify-center gap-4 ${className}`}>
                        {childrenArr.map((item, index) => (
                            <div key={index}>{item}</div>
                        ))}
                    </div>
                </div>
            </div>
        );
    }
}
