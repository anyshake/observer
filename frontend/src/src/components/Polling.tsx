import { Component, ReactNode } from "react";
import ReactPolling from "react-polling";
import { ApiResponse } from "../helpers/request/restfulApiByTag";

export interface PollingProps {
    readonly tag: string;
    readonly timer: number;
    readonly retry?: number;
    readonly onError?: () => void;
    readonly onData: (res: ApiResponse) => void;
    readonly onFetch: (tag: string) => Promise<ApiResponse>;
    readonly children?: ReactNode | ReactNode[];
}

export default class Polling extends Component<PollingProps> {
    render() {
        const { tag, timer, onData, onError, onFetch, children, retry } =
            this.props;
        const childrenArr = Array.isArray(children) ? children : [children];

        return (
            <ReactPolling
                url={tag}
                interval={timer}
                promise={onFetch}
                retryCount={retry}
                onFailure={onError}
                render={() => childrenArr}
                onSuccess={(res: ApiResponse) => {
                    onData && onData(res);
                    return true;
                }}
            />
        );
    }
}
