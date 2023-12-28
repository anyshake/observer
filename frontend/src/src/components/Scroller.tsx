import { Component } from "react";
import ScrollIcon from "../assets/icons/arrow-up-solid.svg";

export interface ScrollerState {
    readonly showButton: boolean;
}

export default class Scroller extends Component<{}, ScrollerState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            showButton: false,
        };
    }

    componentDidMount() {
        document.addEventListener("scroll", this.toggleButton);
    }

    componentWillUnmount() {
        document.removeEventListener("scroll", this.toggleButton);
    }

    toggleButton = (): void => {
        if (window.scrollY > 100) {
            this.setState({ showButton: true });
        } else {
            this.setState({ showButton: false });
        }
    };

    scrollToTop = (): void => {
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    render() {
        const { showButton } = this.state;
        return (
            <button
                onClick={this.scrollToTop}
                className={`fixed bg-purple-500 hover:bg-purple-600 duration-300 w-10 h-10 rounded-full bottom-16 right-3 flex justify-center items-center ${
                    showButton ? "" : "hidden"
                }`}
            >
                <img className="w-4 h-4" src={ScrollIcon} alt="" />
            </button>
        );
    }
}
