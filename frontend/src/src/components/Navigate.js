import { Component } from "react";
import redirectRouter from "../helpers/router/redirectRouter";

export default class Navigate extends Component {
    constructor(props) {
        super(props);
        this.state = {
            dest: this.props.dest || "/",
            replace: this.props.replace || false,
        };
    }

    componentDidMount() {
        redirectRouter({ dest: this.state.dest, replace: this.state.replace });
    }

    render() {
        return null;
    }
}
