import { Component, ReactNode } from "react";
import getRouterTitle from "../helpers/router/getRouterTitle";
import GLOBAL_CONFIG from "../config/global";
import { connect } from "react-redux";
import mapStateToProps from "../helpers/utils/mapStateToProps";
import { withTranslation, WithTranslation } from "react-i18next";

export interface ViewProps extends WithTranslation {
    readonly className?: string;
    readonly children?: ReactNode;
}

class View extends Component<ViewProps> {
    componentDidMount(): void {
        const { t } = this.props;
        const subtitle = getRouterTitle();
        const { title } = GLOBAL_CONFIG.app_settings;
        document.title = `${t(subtitle)} | ${t(title)}`;
    }

    render() {
        const { className, children } = this.props;
        return <div className={className}>{children}</div>;
    }
}

export default connect(mapStateToProps)(withTranslation()(View));
