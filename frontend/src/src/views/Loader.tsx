import { Component } from "react";
import Spinner from "../components/Spinner";
import View from "../components/View";

export default class Loader extends Component {
    render() {
        return (
            <View className="w-full min-h-screen flex flex-col  items-center justify-center">
                <Spinner label={"正在加载中"} />
            </View>
        );
    }
}
