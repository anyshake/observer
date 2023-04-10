import React, { Component } from "react";

export default class Card extends Component {
    constructor(props) {
        super(props);
        this.state = {
            cardList: this.props.cardList || [
                {
                    title: "接收消息量",
                    value: "0",
                    unit: "条",
                    icon: (
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="object-scale-down transition duration-500"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5m.75-9l3-3 2.148 2.148A12.061 12.061 0 0116.5 7.605"
                            />
                        </svg>
                    ),
                },
            ],
        };
    }

    render() {
        return (
            <div className="flex flex-wrap my-5 -mx-2">
                {this.state.cardList.map((item, index) => (
                    <div key={index} className="w-full md:w-1/2 lg:w-1/3 p-2">
                        <div className="flex items-center flex-row w-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-md p-3">
                            <div className="flex text-indigo-500 items-center bg-white p-2 rounded-md flex-none w-8 h-8 md:w-12 md:h-12 ">
                                {item.icon}
                            </div>
                            <div className="flex flex-col justify-around flex-grow ml-5 text-white">
                                <div className="text-sm whitespace-nowrap">
                                    {item.title}
                                </div>
                                <div className="text-md">
                                    {item.value}
                                    <span className="text-lg px-2">
                                        {item.unit}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        );
    }
}
