import { ReactNode, useEffect, useState } from "react";
import collapseIcon from "../assets/icons/square-caret-up-solid.svg";
import advancedIcon from "../assets/icons/ellipsis-solid.svg";
import closeIcon from "../assets/icons/xmark-solid.svg";

export enum CollapseMode {
    COLLAPSE_DISABLE, // Disable collapsing
    COLLAPSE_SHOW, // Enable collapsing and show content
    COLLAPSE_HIDE, // Enable collapsing and hide content
}

export interface HolderProps<T = ReactNode> {
    readonly label: string;
    readonly text?: string;
    readonly children?: T;
    readonly advanced?: T;
    readonly collapse?: CollapseMode;
}

export const Holder = (props: HolderProps) => {
    const { label, text, children, collapse, advanced } = props;
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [advancedOpen, setAdvancedOpen] = useState<boolean>(false);

    useEffect(() => {
        const initCollapse = collapse || CollapseMode.COLLAPSE_DISABLE;
        setCollapsed(initCollapse === CollapseMode.COLLAPSE_HIDE);
    }, [collapse]);

    const currentCollapse = collapse || CollapseMode.COLLAPSE_DISABLE;
    const collapse_is_enabled =
        currentCollapse !== CollapseMode.COLLAPSE_DISABLE;

    return (
        <div className="mb-4 flex flex-col rounded-xl text-gray-700 shadow-lg">
            <div className="mx-4 rounded-lg overflow-hidden shadow-lg">
                {children}
            </div>

            <div className="p-4">
                <h6
                    className={`text-md font-bold text-gray-800 flex ${
                        collapse_is_enabled ? "cursor-pointer select-none" : ""
                    }`}
                    onClick={() =>
                        collapse_is_enabled && setCollapsed(!collapsed)
                    }
                >
                    {collapse_is_enabled && (
                        <img
                            className={`mx-1 ${collapsed ? "rotate-180" : ""}`}
                            src={collapseIcon}
                            alt=""
                        />
                    )}
                    {label}
                </h6>
                {text && !collapsed && (
                    <div className="text-md pt-2">
                        {text.split("\n").map((item) => (
                            <div key={item}>{item}</div>
                        ))}
                    </div>
                )}
                {advanced && (
                    <div
                        className={`mt-2 space-y-2 ${
                            collapsed ? "hidden" : "block"
                        }`}
                    >
                        <div
                            className="mx-1 cursor-pointer"
                            onClick={() => setAdvancedOpen(!advancedOpen)}
                        >
                            <img
                                className={`size-4 ${
                                    advancedOpen ? "hidden" : "block"
                                }`}
                                src={advancedIcon}
                                alt=""
                            />
                            <img
                                className={`size-4 ${
                                    advancedOpen ? "block" : "hidden"
                                }`}
                                src={closeIcon}
                                alt=""
                            />
                        </div>
                        <div className={advancedOpen ? "block" : "hidden"}>
                            {advanced}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
