import { Component } from "react";

interface ButtonProps {
    readonly className?: string;
    readonly label: string;
    readonly onClick?: () => void;
}

export default class Button extends Component<ButtonProps> {
    constructor(props: ButtonProps) {
        super(props);
        this.state = {
            isBusy: false,
        };
    }

    render() {
        const { className, label, onClick } = this.props;
        return (
            <button
                className={`w-full text-white font-medium text-sm shadow-lg rounded-lg py-2 ${className}`}
                onClick={onClick}
            >
                {label}
            </button>
        );
    }
}
