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
				<div className="mt-4 px-6 py-2 font-bold">
					{sublabel && <h5 className="text-gray-500 text-xs">{sublabel}</h5>}
					<h2 className={embedded ? "text-md" : "text-lg"}>{label}</h2>
				</div>
				<div className={`p-4 m-2 flex flex-col justify-center gap-4 ${className ?? ""}`}>
					{children}
				</div>
			</div>
		</div>
	);
};
