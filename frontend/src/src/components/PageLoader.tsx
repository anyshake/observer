import { Component } from "react";
import SpinnerIcon from "../assets/icons/spinner-solid.svg";

export default class PageLoader extends Component {
    render() {
        return (
            <div className="w-full min-h-screen flex flex-col items-center justify-center">
                <div className="animate-spin">
                    <img className="py-2 w-20 h-20" src={SpinnerIcon} alt="" />
                </div>

                <h2 className="py-2 text-2xl font-bold text-gray-600">
                    Loading...
                </h2>
            </div>
        );
    }
}
