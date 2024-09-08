import { TextField } from "@mui/material";
import { ChangeEvent, InputHTMLAttributes } from "react";

interface InputProps {
	readonly label: string;
	readonly disabled?: boolean;
	readonly className?: string;
	readonly defaultValue: string | number;
	readonly numberLimit?: { max: number; min: number };
	readonly type: InputHTMLAttributes<unknown>["type"];
	readonly onValueChange?: (value: string | number) => void;
}

export const Input = (props: InputProps) => {
	const { label, disabled, className, defaultValue, numberLimit, type, onValueChange } = props;

	const handleOnChange = ({ target }: ChangeEvent<HTMLTextAreaElement>) => {
		if (!onValueChange) {
			return;
		}
		const { value } = target;
		if (type === "number") {
			const numberValue = Number(value);
			if (isNaN(numberValue)) {
				onValueChange(defaultValue);
				return;
			}
			if (numberLimit) {
				const { max, min } = numberLimit;
				if (numberValue > max || numberValue < min) {
					onValueChange(defaultValue);
					return;
				}
			}
			onValueChange(numberValue);
		} else {
			onValueChange(value);
		}
	};

	return (
		<TextField
			size="small"
			type={type}
			label={label}
			disabled={disabled}
			onChange={handleOnChange}
			defaultValue={defaultValue}
			className={`w-full ${className ?? ""}`}
			InputLabelProps={{ shrink: true }}
		/>
	);
};
