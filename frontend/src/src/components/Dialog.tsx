import {
	Button,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogTitle,
	FormControl
} from "@mui/material";
import MuiDialog from "@mui/material/Dialog";
import { ReactNode } from "react";

export interface DialogProps {
	readonly open: boolean;
	readonly title?: string;
	readonly content?: string;
	readonly cancelText?: string;
	readonly submitText?: string;
	readonly onClose?: () => void;
	readonly onSubmit?: () => void;
	readonly children?: ReactNode | ReactNode[];
}

export const Dialog = ({
	open,
	title,
	content,
	cancelText,
	submitText,
	onClose,
	onSubmit,
	children
}: DialogProps) => {
	return (
		<MuiDialog fullWidth={true} onClose={onClose} open={open}>
			<DialogTitle>{title}</DialogTitle>
			<DialogContent>
				{content && <DialogContentText sx={{ py: 2 }}>{content}</DialogContentText>}
				<FormControl sx={{ mt: 2, width: "100%" }}>{children}</FormControl>
			</DialogContent>
			<DialogActions>
				{cancelText && <Button onClick={onClose}>{cancelText}</Button>}
				<Button onClick={onSubmit}>{submitText}</Button>
			</DialogActions>
		</MuiDialog>
	);
};
