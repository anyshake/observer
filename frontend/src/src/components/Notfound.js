import React, { Component } from "react";
import Navigate from "./Navigate";

export default class NotFound extends Component {
    constructor(props) {
        super(props);
        this.state = {
            countDown: 5,
            countDownTimer: null,
        };
    }

    componentDidMount() {
        this.setState({
            countDownTimer: setInterval(() => {
                this.setState({
                    countDown: this.state.countDown - 1,
                });
                if (this.state.countDown === 1) {
                    clearInterval(this.state.countDownTimer);
                }
            }, 1000),
        });
    }

    render() {
        return (
            <div className="w-full h-screen bg-gradient-to-br text-white from-purple-600 to-blue-600 flex flex-col items-center justify-center">
                <svg
                    className="py-2 w-16 h-16 animate-bounce fill-gray-200"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        fillRule="evenodd"
                        d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                        clipRule="evenodd"
                    />
                </svg>
                <h2 className="py-2 text-center text-xl font-bold">
                    404 - 查无此页
                </h2>
                {this.state.countDown === 0 ? (
                    <Navigate dest="/" replace={true} />
                ) : (
                    <i className="px-14 text-center">
                        剩余 {this.state.countDown} 秒后返回首页
                    </i>
                )}
            </div>
        );
    }
}
