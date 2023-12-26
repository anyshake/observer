import { Component, ReactNode } from "react";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";
import CollapseIcon from "../assets/icons/square-caret-up-solid.svg";

export enum CollapseMode {
    COLLAPSE_DISABLE, // Disable collapsing
    COLLAPSE_SHOW, // Enable collapsing and show content
    COLLAPSE_HIDE, // Enable collapsing and hide content
}

export interface AreaProps {
    readonly label: I18nTranslation;
    readonly text?: I18nTranslation;
    readonly children?: ReactNode;
    readonly collapse?: CollapseMode;
}

export interface AreaState {
    readonly collapsed: boolean;
}

class Area extends Component<AreaProps & WithTranslation, AreaState> {
    constructor(props: AreaProps & WithTranslation) {
        super(props);
        this.state = {
            collapsed: false,
        };
    }

    componentDidMount(): void {
        const collapse = this.props.collapse || CollapseMode.COLLAPSE_DISABLE;
        this.setState({ collapsed: collapse === CollapseMode.COLLAPSE_HIDE });
    }

    render() {
        const { collapsed } = this.state;
        const { t, children, label, text } = this.props;
        const collapse = this.props.collapse || CollapseMode.COLLAPSE_DISABLE;
        const collapse_is_enabled = collapse !== CollapseMode.COLLAPSE_DISABLE;

        return (
            <div className="mb-4 flex flex-col rounded-xl text-gray-700 shadow-lg">
                <div className="mx-4 rounded-lg overflow-hidden shadow-lg">
                    {children}
                </div>

                <div className="p-4">
                    <h6
                        className={`text-md font-bold text-gray-800 flex ${
                            collapse_is_enabled && "cursor-pointer select-none"
                        }`}
                        onClick={() =>
                            collapse_is_enabled &&
                            this.setState({ collapsed: !collapsed })
                        }
                    >
                        {collapse_is_enabled && (
                            <img
                                className={`mx-1 ${collapsed && "rotate-180"}`}
                                src={CollapseIcon}
                                alt=""
                            />
                        )}
                        {t(label.id, label.format)}
                    </h6>
                    {text && !collapsed && (
                        <span className="text-md">
                            {t(text.id, text.format)
                                .split("\n")
                                .map((item: string, key: number) => (
                                    <p key={key}>
                                        {item}
                                        <br />
                                    </p>
                                ))}
                        </span>
                    )}
                </div>
            </div>
        );
    }
}

export default withTranslation()(Area);
