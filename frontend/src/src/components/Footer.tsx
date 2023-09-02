import { Component } from "react";
import GLOBAL_CONFIG from "../config/global";

export interface FooterProps {
    readonly text?: string;
}

export default class Footer extends Component<FooterProps> {
    render() {
        const { text } = this.props;
        const { author, description } = GLOBAL_CONFIG.app_settings;

        return (
            <footer className="w-full bg-gray-200 text-gray-500 flex flex-col px-6 py-2 sm:flex-row justify-between">
                <span className="text-xs text-center ml-8 md:ml-12">
                    {!!text ? text : description}
                </span>

                <span className="text-sm text-center justify-center">
                    {`© ${new Date().getFullYear()} ${author}`}
                </span>
            </footer>
        );
    }
}
