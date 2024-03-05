import {
    Button,
    Dialog,
    AppBar,
    Toolbar,
    Typography,
    List,
    ListItem,
    ListItemText,
    Divider,
    ListItemButton,
} from "@mui/material";

export interface SelectProps {
    readonly open: boolean;
    readonly title?: string;
    readonly options?: (string)[][];
    readonly onClose?: () => void;
    readonly onSelect?: (value: string) => void;
}

export const Select = (props: SelectProps) => {
    const { title, open, options, onClose, onSelect } = props;

    return (
        <Dialog fullWidth onClose={onClose} open={open}>
            <AppBar className="bg-violet-500" sx={{ position: "relative" }}>
                <Toolbar>
                    <Typography sx={{ mt: 1, ml: 1, flex: 1 }} variant="h6">
                        {title}
                    </Typography>
                    <Button autoFocus color="inherit" onClick={onClose}>
                        X
                    </Button>
                </Toolbar>
            </AppBar>

            <List>
                {!!options?.length &&
                    options.map(
                        (item) =>
                            item.length > 1 && (
                                <div key={item[1]}>
                                    <ListItem>
                                        <ListItemButton
                                            onClick={() =>
                                                onSelect && onSelect(item[1])
                                            }
                                        >
                                            <ListItemText
                                                primary={item[0]}
                                                secondary={item[
                                                    item.length === 3 ? 2 : 1
                                                ]
                                                    .split("\n")
                                                    .map((item) => (
                                                        <span key={item}>
                                                            {item}
                                                            <br />
                                                        </span>
                                                    ))}
                                            />
                                        </ListItemButton>
                                    </ListItem>
                                    <Divider />
                                </div>
                            )
                    )}
            </List>
        </Dialog>
    );
};
