import { Component } from "react";

export interface TextProps {
    readonly className?: string;
    readonly children: string;
}

export default class Text extends Component<TextProps> {
    render() {
        const { className, children } = this.props;
        return <div className={className}>{children}</div>;
    }
}
