import { Component } from "react";
import GLOBAL_CONFIG from "../config/global";
import { WithTranslation, withTranslation } from "react-i18next";

class Footer extends Component<WithTranslation> {
    render() {
        const { t } = this.props;
        const { author, description } = GLOBAL_CONFIG.app_settings;

        return (
            <footer className="w-full bg-gray-200 text-gray-500 flex flex-col px-6 py-2 sm:flex-row justify-between">
                <span className="text-xs text-center ml-8 md:ml-12">
                    {t(description)}
                </span>

                <span className="text-sm text-center justify-center">
                    {`© ${new Date().getFullYear()} ${t(author)}`}
                </span>
            </footer>
        );
    }
}

export default withTranslation()(Footer);
