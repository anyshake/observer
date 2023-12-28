import { Component, ReactNode } from "react";

export interface ContainerProps {
    readonly className?: string;
    readonly layout: "none" | "grid" | "flex";
    readonly children?: ReactNode;
}

export default class Container extends Component<ContainerProps> {
    render() {
        const { layout, className, children } = this.props;

        let layoutClassName = "mt-5";
        switch (layout) {
            case "flex":
                layoutClassName = "mt-5 flex flex-wrap";
                break;
            case "grid":
                layoutClassName = "mt-5 gap-4 grid grid-cols-1 md:grid-cols-2";
                break;
        }

        return (
            <div className={`${layoutClassName} ${className ?? ""}`}>
                {children}
            </div>
        );
    }
}
