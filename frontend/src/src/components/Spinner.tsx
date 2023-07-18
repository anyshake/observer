import { Component } from "react";
import SpinnerIcon from "../assets/icons/spinner-solid.svg";

export interface SpinnerProps {
    readonly label: string;
}

export default class Spinner extends Component<SpinnerProps> {
    render() {
        const { label } = this.props;
        return (
            <>
                <div className="animate-spin">
                    <img className="py-2 w-20 h-20" src={SpinnerIcon} alt="" />
                </div>

                <h2 className="py-2 text-2xl font-bold text-gray-600">
                    {label}
                </h2>
            </>
        );
    }
}
