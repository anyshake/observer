export interface LabelProps {
	readonly className?: string;
	readonly icon?: string;
	readonly color?: boolean;
	readonly value: string;
	readonly unit?: string;
	readonly label: string;
}

export const Label = (props: LabelProps) => {
	const { className, icon, label, value, unit, color } = props;

	return (
		<div className={`w-full p-2 ${className ?? ""}`}>
			<div
				className={`flex flex-row bg-gradient-to-r rounded-md p-4 shadow-xl ${
					color
						? `from-indigo-500 via-purple-500 to-pink-500`
						: `bg-gray-50 hover:bg-gray-100 transition-all`
				}`}
			>
				{icon && (
					<img
						className="bg-white p-2 rounded-md w-8 h-8 md:w-12 md:h-12 self-center"
						src={icon}
						alt=""
					/>
				)}

				<div className={`flex flex-col flex-grow ${icon ? `ml-5` : ""}`}>
					<div
						className={`text-sm whitespace-nowrap ${
							color ? `text-gray-50` : `text-gray-600`
						}`}
					>
						{label}
					</div>
					<div
						className={`text-md font-medium flex-nowrap ${
							color ? `text-gray-100` : `text-gray-800`
						}`}
					>
						{`${value} ${unit ?? ""}`}
					</div>
				</div>
			</div>
		</div>
	);
};
