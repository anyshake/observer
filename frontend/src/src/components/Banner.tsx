import { Component } from "react";
import SuccessIcon from "../assets/icons/rss-solid.svg";
import WarningIcon from "../assets/icons/link-solid.svg";
import ErrorIcon from "../assets/icons/link-slash-solid.svg";

export interface BannerProps {
    readonly type: "success" | "warning" | "error";
    readonly label: string;
    readonly text: string;
}

export default class Banner extends Component<BannerProps> {
    render() {
        const { type, label, text } = this.props;

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
                        <span>{label}</span>
                    </div>

                    <span className="pl-3 text-md font-medium">
                        {text.split("\n").map((item, index) => (
                            <p key={index}>{item}</p>
                        ))}
                    </span>
                </div>
            </div>
        );
    }
}
