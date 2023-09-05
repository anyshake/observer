import { Component } from "react";
import SuccessIcon from "../assets/icons/rss-solid.svg";
import WarningIcon from "../assets/icons/link-solid.svg";
import ErrorIcon from "../assets/icons/link-slash-solid.svg";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

export interface BannerProps {
    readonly text: I18nTranslation;
    readonly label: I18nTranslation;
    readonly type: "success" | "warning" | "error";
}

class Banner extends Component<BannerProps & WithTranslation> {
    render() {
        const { t, type, label, text } = this.props;

        let icon = ErrorIcon;
        let colorClassName = "";
        switch (type) {
            case "success":
                icon = SuccessIcon;
                colorClassName = "from-green-400 to-blue-500";
                break;
            case "warning":
                icon = WarningIcon;
                colorClassName = "from-orange-400 to-orange-600";
                break;
            case "error":
                icon = ErrorIcon;
                colorClassName = "from-red-400 to-red-600";
                break;
        }

        return (
            <div
                className={`my-2 shadow-xl p-6 text-sm text-white rounded-lg bg-gradient-to-r ${colorClassName}`}
            >
                <div className="flex flex-col gap-y-2">
                    <div className="flex gap-2 font-bold text-lg">
                        <img className="w-6 h-6" src={icon} alt="" />
                        <span>{t(label.id, label.format)}</span>
                    </div>

                    <span className="pl-3 text-md font-medium">
                        {t(text.id, text.format)
                            .split("\n")
                            .map((item: string, key: number) => (
                                <p key={key}>
                                    {item}
                                    <br />
                                </p>
                            ))}
                    </span>
                </div>
            </div>
        );
    }
}

export default withTranslation()(Banner);
