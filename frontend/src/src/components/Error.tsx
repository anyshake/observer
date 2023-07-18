import { Component } from "react";
import ErrorIcon from "../assets/icons/circle-exclamation-solid.svg";

export interface ErrorProps {
    readonly label: string;
}

export default class Error extends Component<ErrorProps> {
    render() {
        const { label } = this.props;
        return (
            <>
                <div className="animate-bounce">
                    <img className="py-2 w-20 h-20" src={ErrorIcon} alt="" />
                </div>

                <h2 className="py-2 text-2xl font-bold text-gray-600">
                    {label}
                </h2>
            </>
        );
    }
}
