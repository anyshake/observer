import { Component, ReactNode } from "react";
import getRouterTitle from "../helpers/router/getRouterTitle";
import GLOBAL_CONFIG from "../config/global";

export interface ViewProps {
    readonly className?: string;
    readonly children?: ReactNode;
}

export default class View extends Component<ViewProps> {
    componentDidMount(): void {
        const subtitle = getRouterTitle();
        const { title } = GLOBAL_CONFIG.app_settings;
        document.title = `${subtitle} | ${title}`;
    }

    render() {
        const { className, children } = this.props;
        return <div className={className}>{children}</div>;
    }
}
