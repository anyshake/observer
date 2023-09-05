import { Component, ReactNode } from "react";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

interface CardProps extends WithTranslation {
    readonly className?: string;
    readonly label: I18nTranslation;
    readonly sublabel?: I18nTranslation;
    readonly children: ReactNode | ReactNode[];
}

class Card extends Component<CardProps> {
    render() {
        const { t, className, label, sublabel, children } = this.props;
        const childrenArr = Array.isArray(children) ? children : [children];

        return (
            <div className="w-full h-full text-gray-800">
                <div className="flex flex-col shadow-lg rounded-lg">
                    <div className="px-4 py-3 font-bold">
                        {sublabel && (
                            <h6 className="text-gray-500 text-xs">
                                {t(sublabel.id, sublabel.format)}
                            </h6>
                        )}
                        <h2 className="text-xl">{t(label.id, label.format)}</h2>
                    </div>

                    <div
                        className={`p-4 m-2 flex flex-col justify-center gap-4 ${className}`}
                    >
                        {childrenArr.map((item, index) => (
                            <div key={index}>{item}</div>
                        ))}
                    </div>
                </div>
            </div>
        );
    }
}

export default withTranslation()(Card);
