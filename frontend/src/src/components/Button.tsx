import { Component } from "react";
import { I18nTranslation } from "../config/i18n";
import { WithTranslation, withTranslation } from "react-i18next";

interface ButtonProps extends WithTranslation {
    readonly label: I18nTranslation;
    readonly className?: string;
    readonly onClick?: () => void;
}

class Button extends Component<ButtonProps> {
    constructor(props: ButtonProps) {
        super(props);
        this.state = {
            isBusy: false,
        };
    }

    render() {
        const { t, className, label, onClick } = this.props;
        return (
            <button
                className={`w-full text-white font-medium text-sm shadow-lg rounded-lg py-2 ${
                    className ?? ""
                }`}
                onClick={onClick}
            >
                {t(label.id, label.format)}
            </button>
        );
    }
}

export default withTranslation()(Button);
