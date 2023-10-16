import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    TextField,
} from "@mui/material";
import { Component, HTMLInputTypeAttribute, createRef } from "react";
import { I18nTranslation } from "../config/i18n";
import { WithTranslation, withTranslation } from "react-i18next";

export interface FormDialogProps {
    readonly tag?: string;
    readonly open: boolean;
    readonly defaultValue?: string;
    readonly title: I18nTranslation;
    readonly content?: I18nTranslation;
    readonly cancelText?: I18nTranslation;
    readonly submitText?: I18nTranslation;
    readonly placeholder?: I18nTranslation;
    readonly inputType: HTMLInputTypeAttribute;
    readonly onSubmit?: (value: string) => void;
    readonly onClose?: () => void;
}

interface FormDialogState {
    readonly input: React.RefObject<HTMLInputElement>;
}

class FormDialog extends Component<
    FormDialogProps & WithTranslation,
    FormDialogState
> {
    constructor(props: FormDialogProps & WithTranslation) {
        super(props);
        this.state = {
            input: createRef(),
        };
    }

    render() {
        const {
            t,
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
        } = this.props;
        const { input } = this.state;

        return (
            <Dialog open={open}>
                <DialogTitle>{t(title.id, title.format)}</DialogTitle>
                <DialogContent>
                    {content && (
                        <DialogContentText>
                            {t(content.id, content.format)}
                        </DialogContentText>
                    )}
                    <TextField
                        autoFocus
                        fullWidth
                        inputRef={input}
                        className="mt-8"
                        type={inputType}
                        defaultValue={defaultValue}
                        label={t(placeholder?.id ?? "", placeholder?.format)}
                    />
                </DialogContent>
                <DialogActions>
                    {cancelText && (
                        <Button onClick={onClose}>
                            {t(cancelText.id, cancelText.format)}
                        </Button>
                    )}
                    <Button
                        onClick={() =>
                            onSubmit && onSubmit(input.current?.value ?? "")
                        }
                    >
                        {submitText
                            ? t(submitText.id, submitText.format)
                            : "Submit"}
                    </Button>
                </DialogActions>
            </Dialog>
        );
    }
}

export default withTranslation()(FormDialog);
