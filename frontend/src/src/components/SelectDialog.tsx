import {
    Dialog,
    DialogTitle,
    ListItem,
    ListItemButton,
    ListItemText,
} from "@mui/material";
import { Component } from "react";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

export interface SelectDialogProps {
    readonly title: I18nTranslation;
    readonly open: boolean;
    readonly values: string[][];
    readonly onClose?: (value: string) => void;
    readonly onSelect?: (value: string) => void;
}

class SelectDialog extends Component<SelectDialogProps & WithTranslation> {
    render() {
        const { t, title, open, values, onClose, onSelect } = this.props;
        return (
            <Dialog onClose={onClose} open={open}>
                <DialogTitle>{t(title.id, title.format)}</DialogTitle>
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

export default withTranslation()(SelectDialog);
