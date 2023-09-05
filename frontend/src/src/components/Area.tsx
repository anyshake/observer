import { Component, ReactNode } from "react";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

export interface AreaProps {
    readonly label: I18nTranslation;
    readonly text?: I18nTranslation;
    readonly children?: ReactNode;
}

class Area extends Component<AreaProps & WithTranslation> {
    render() {
        const { t, children, label, text } = this.props;

        return (
            <div className="mb-4 flex flex-col rounded-xl text-gray-700 shadow-lg">
                <div className="mx-4 rounded-lg overflow-hidden shadow-lg">
                    {children}
                </div>

                <div className="p-4">
                    <h6 className="text-md font-bold text-gray-800">
                        {t(label.id, label.format)}
                    </h6>
                    {text && (
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
