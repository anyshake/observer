import {
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogTitle,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	TextField
} from "@mui/material";
import { HTMLInputTypeAttribute, useEffect, useState } from "react";

export interface FormProps {
	readonly open: boolean;
	readonly title?: string;
	readonly content?: string;
	readonly cancelText?: string;
	readonly submitText?: string;
	readonly placeholder?: string;
	readonly defaultValue?: string;
	readonly onClose?: () => void;
	readonly onSubmit?: (value: string) => void;
	readonly inputType?: HTMLInputTypeAttribute | "select" | "textarea";
	readonly selectOptions?: { value: string; label: string }[];
}

export const Form = (props: FormProps) => {
	const {
		open,
		title,
		content,
		cancelText,
		submitText,
		placeholder,
		defaultValue,
		inputType,
		onSubmit,
		onClose,
		selectOptions
	} = props;

	const [inputValue, setInputValue] = useState("");
	const [selectValue, setSelectValue] = useState("");

	useEffect(() => {
		setSelectValue(selectOptions?.[0].value ?? "");
	}, [selectOptions]);

	return (
		<Dialog onClose={onClose} open={open}>
			<DialogTitle>{title}</DialogTitle>
			<DialogContent>
				{content && <DialogContentText sx={{ py: 2 }}>{content}</DialogContentText>}
				<TextField
					autoFocus
					fullWidth
					sx={{ my: 2 }}
					type={inputType}
					label={placeholder}
					defaultValue={defaultValue}
					style={{
						display: inputType !== "select" ? "block" : "none"
					}}
					onChange={({ target }) => {
						setInputValue(target.value);
					}}
					multiline={inputType === "textarea"}
				/>
				<FormControl
					sx={{ my: 2 }}
					style={{
						display: inputType === "select" ? "block" : "none"
					}}
				>
					<InputLabel>{placeholder}</InputLabel>
					<Select
						fullWidth
						label={placeholder}
						onChange={({ target }) => {
							setSelectValue(target.value);
						}}
						defaultValue={selectOptions?.[0].value ?? ""}
					>
						{selectOptions?.map(({ value, label }) => (
							<MenuItem key={value} value={value}>
								{label}
							</MenuItem>
						))}
					</Select>
				</FormControl>
			</DialogContent>
			<DialogActions>
				{cancelText && <Button onClick={onClose}>{cancelText}</Button>}
				<Button
					onClick={() => {
						onSubmit?.(
							inputType === "select"
								? selectValue
								: inputValue.length
									? inputValue
									: defaultValue ?? ""
						);
					}}
				>
					{submitText}
				</Button>
			</DialogActions>
		</Dialog>
	);
};
