import { ReactNode } from "react";

interface PanelProps<T = ReactNode> {
    readonly embedded?: boolean;
    readonly className?: string;
    readonly sublabel?: string;
    readonly label: string;
    readonly children: T;
}

export const Panel = (props: PanelProps) => {
    const { embedded, className, label, sublabel, children } = props;

    return (
        <div className="w-full text-gray-800">
            <div className="flex flex-col shadow-lg rounded-lg">
                <div className="px-4 py-3 font-bold">
                    {sublabel && (
                        <h6 className="text-gray-500 text-xs">{sublabel}</h6>
                    )}
                    <h2 className={embedded ? "text-md" : "text-lg"}>
                        {label}
                    </h2>
                </div>
                <div
                    className={`p-4 m-2 flex flex-col justify-center gap-4 ${
                        className ?? ""
                    }`}
                >
                    {children}
                </div>
            </div>
        </div>
    );
};
