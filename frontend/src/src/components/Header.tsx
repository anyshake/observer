import { ChangeEvent, Component } from "react";
import EarthIcon from "../assets/icons/earth-americas-solid.svg";
import GLOBAL_CONFIG from "../config/global";
import mapStateToProps from "../helpers/utils/mapStateToProps";
import { connect } from "react-redux";
import I18N_CONFIG, { i18n } from "../config/i18n";
import toggleI18n from "../helpers/i18n/toggleI18n";
import getLanguage from "../helpers/i18n/getLanguage";
import { WithTranslation, withTranslation } from "react-i18next";

interface HeaderState {
    readonly i18n: string;
}

class Header extends Component<WithTranslation, HeaderState> {
    constructor(props: WithTranslation) {
        super(props);
        this.state = {
            i18n: getLanguage(),
        };
    }

    handleI18nChange = (e: ChangeEvent<HTMLSelectElement>) => {
        const { value } = e.target;
        toggleI18n(i18n, value);
        this.setState({
            i18n: value,
        });
    };

    render() {
        const { t } = this.props;
        const { i18n } = this.state;
        const { name } = GLOBAL_CONFIG.app_settings;

        return (
            <header className="fixed w-full z-10 flex justify-between bg-gray-200 items-center h-16 px-5">
                <h1 className="ml-14 text-gray-800 text-xl font-bold">
                    {t(name)}
                </h1>

                <div className="flex text-gray-500 flex-nowrap space-x-1">
                    <img className="w-4 h-4" src={EarthIcon} alt="" />
                    <select
                        className="text-xs bg-transparent focus:outline-none"
                        onChange={this.handleI18nChange}
                        value={i18n}
                    >
                        {I18N_CONFIG.list.map(({ name, value }, index) => (
                            <option key={index} value={value}>
                                {name}
                            </option>
                        ))}
                    </select>
                </div>
            </header>
        );
    }
}

export default connect(mapStateToProps)(withTranslation()(Header));
