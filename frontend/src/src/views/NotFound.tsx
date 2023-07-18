import { Component } from "react";
import Error from "../components/Error";
import View from '../components/View';

export default class NotFound extends Component {
    render() {
        return (
            <View className="w-full min-h-screen flex flex-col items-center justify-center">
                <Error label={"找不到页面"} />
            </View>
        );
    }
}
