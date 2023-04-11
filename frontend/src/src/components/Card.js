import React, { Component } from "react";

export default class Card extends Component {
    constructor(props) {
        super(props);
        this.state = {
            icon: this.props.icon,
            title: this.props.title,
            unit: this.props.unit,
        };
    }

    render() {
        return (
            <div className="w-full md:w-1/2 lg:w-1/3 p-2">
                <div className="flex items-center flex-row w-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-md p-3 shadow-xl">
                    <div className="flex text-indigo-500 items-center bg-white p-2 rounded-md flex-none w-8 h-8 md:w-12 md:h-12 ">
                        {this.state.icon}
                    </div>
                    <div className="flex flex-col justify-around flex-grow ml-5 text-white">
                        <div className="text-sm whitespace-nowrap">
                            {this.state.title}
                        </div>
                        <div className="text-md">
                            {this.props.value}
                            <span className="text-lg px-2">
                                {this.state.unit}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
