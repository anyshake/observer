import { ForwardedRef, ReactNode, forwardRef } from "react";
import { Toaster } from "react-hot-toast";

export interface ContainerProps<T = ReactNode> {
    readonly children: T;
    readonly main?: boolean;
    readonly toaster?: boolean;
    readonly className?: string;
}

export const Container = forwardRef(
    (props: ContainerProps, ref: ForwardedRef<HTMLDivElement>) => {
        const { main, className, toaster, children } = props;

        return (
            <div
                className={
                    main
                        ? "bg-gray-50 min-h-screen ml-10 p-20 px-4 flex flex-col space-y-3"
                        : className ?? ""
                }
                ref={ref}
            >
                {children}
                {toaster && <Toaster />}
            </div>
        );
    }
);
