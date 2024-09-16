import { FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";
import { InputHTMLAttributes } from "react";

interface InputProps {
	readonly label: string;
	readonly disabled?: boolean;
	readonly className?: string;
	readonly fullWidth?: boolean;
	readonly defaultValue?: string | number;
	readonly numberLimit?: { max: number; min: number };
	readonly type: InputHTMLAttributes<unknown>["type"] | "select";
	readonly onValueChange?: (value: string | number) => void;
	readonly selectOptions?: { value: string; label: string }[];
}

export const Input = ({
	label,
	disabled,
	className,
	defaultValue,
	fullWidth,
	numberLimit,
	type,
	onValueChange,
	selectOptions
}: InputProps) => {
	const handleChange = (value: string) => {
		if (!onValueChange) {
			return;
		}
		if (type === "number") {
			const numberValue = Number(value);
			if (isNaN(numberValue)) {
				onValueChange(defaultValue ?? "");
				return;
			}
			if (numberLimit) {
				const { max, min } = numberLimit;
				if (numberValue > max || numberValue < min) {
					onValueChange(defaultValue ?? "");
					return;
				}
			}
			onValueChange(numberValue);
		} else {
			onValueChange(value);
		}
	};

	return (
		<FormControl fullWidth={fullWidth} sx={{ minWidth: 80 }}>
			{type === "select" ? (
				<>
					<InputLabel>{label}</InputLabel>
					<Select
						size="small"
						label={label}
						disabled={disabled}
						onChange={({ target }) => handleChange(String(target.value))}
						defaultValue={selectOptions?.[0].value}
					>
						{selectOptions?.map(({ value, label }) => (
							<MenuItem key={value} value={value}>
								{label}
							</MenuItem>
						))}
					</Select>
				</>
			) : (
				<TextField
					size="small"
					type={type}
					label={label}
					disabled={disabled}
					defaultValue={defaultValue}
					className={`w-full ${className ?? ""}`}
					InputLabelProps={{ shrink: true }}
					onChange={({ target }) => handleChange(target.value)}
				/>
			)}
		</FormControl>
	);
};
