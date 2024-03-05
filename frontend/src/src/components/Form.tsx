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
    TextField,
} from "@mui/material";
import { HTMLInputTypeAttribute, useEffect, useRef, useState } from "react";

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
    readonly inputType?: HTMLInputTypeAttribute | "select";
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
        selectOptions,
    } = props;

    const inputRef = useRef<HTMLInputElement>(null);
    const [selectValue, setSelectValue] = useState("");

    useEffect(() => {
        setSelectValue(selectOptions?.[0].value ?? "");
    }, [selectOptions]);

    return (
        <Dialog onClose={onClose} open={open}>
            <DialogTitle>{title}</DialogTitle>
            <DialogContent>
                {content && <DialogContentText>{content}</DialogContentText>}
                {inputType !== "select" ? (
                    <TextField
                        autoFocus
                        fullWidth
                        ref={inputRef}
                        className="mt-8"
                        type={inputType}
                        label={placeholder}
                        defaultValue={defaultValue}
                    />
                ) : (
                    <FormControl sx={{ my: 2 }} fullWidth>
                        <InputLabel id="select">{placeholder}</InputLabel>
                        <Select
                            labelId="select"
                            label={placeholder}
                            defaultValue={selectOptions?.[0].value}
                            onChange={({ target }) =>
                                setSelectValue(target?.value)
                            }
                        >
                            {selectOptions?.map(({ value, label }) => (
                                <MenuItem key={value} value={value}>
                                    {label}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                )}
            </DialogContent>
            <DialogActions>
                {cancelText && <Button onClick={onClose}>{cancelText}</Button>}
                <Button
                    onClick={() => {
                        if (onSubmit) {
                            if (inputType === "select") {
                                onSubmit(selectValue);
                            } else {
                                onSubmit(inputRef.current?.value ?? "");
                            }
                        }
                    }}
                >
                    {submitText}
                </Button>
            </DialogActions>
        </Dialog>
    );
};
