import React, { Component } from "react";

export default class Notification extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <div className={this.props.className}>
                <div className="flex flex-col gap-y-2 font-bold">
                    <div className="flex flex-row gap-2 font-bold text-lg">
                        {this.props.icon}
                        <span>{this.props.title}</span>
                    </div>
                    <span className="pl-3 text-md font-medium">
                        {this.props.text}
                    </span>
                </div>
            </div>
        );
    }
}
