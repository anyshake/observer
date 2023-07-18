import {
    Dialog,
    DialogTitle,
    ListItem,
    ListItemButton,
    ListItemText,
} from "@mui/material";
import { Component } from "react";

export interface SelectDialogProps {
    readonly title: string;
    readonly open: boolean;
    readonly values: string[][];
    readonly onClose?: (value: string) => void;
    readonly onSelect?: (value: string) => void;
}

export default class SelectDialog extends Component<SelectDialogProps> {
    render() {
        const { title, open, values, onClose, onSelect } = this.props;
        return (
            <Dialog onClose={onClose} open={open}>
                <DialogTitle>{title}</DialogTitle>
                {values.map((item, index) => (
                    <ListItem key={index} disableGutters>
                        <ListItemButton
                            onClick={() => onSelect && onSelect(item[1])}
                        >
                            <ListItemText
                                primary={item[0]}
                                secondary={item[3] ? item[3] : item[1]}
                            />
                        </ListItemButton>
                    </ListItem>
                ))}
            </Dialog>
        );
    }
}
