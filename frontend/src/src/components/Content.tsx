import { Component, ReactNode } from "react";

export interface ContentProps {
    readonly children?: ReactNode | ReactNode[];
}

export default class Content extends Component<ContentProps> {
    render() {
        const { children } = this.props;
        const childrenArr = Array.isArray(children) ? children : [children];

        return (
            <main className="bg-gray-50 min-h-screen ml-10 p-20 px-4 flex flex-col space-y-3">
                {childrenArr.map((item, index) => (
                    <div key={index}>{item}</div>
                ))}
            </main>
        );
    }
}
