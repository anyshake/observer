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
import { Component } from "react";
import { I18nTranslation } from "../config/i18n";
import { WithTranslation, withTranslation } from "react-i18next";

export interface ModalDialogProps {
    readonly open: boolean;
    readonly values: string[][];
    readonly onClose?: () => void;
    readonly title: I18nTranslation;
    readonly onSelect?: (value: string) => void;
}

class ModalDialog extends Component<ModalDialogProps & WithTranslation> {
    render() {
        const { t, title, open, values, onClose, onSelect } = this.props;
        return (
            <Dialog fullWidth open={open}>
                <AppBar className="bg-purple-500" sx={{ position: "relative" }}>
                    <Toolbar>
                        <Typography sx={{ ml: 2, flex: 1 }} variant="h6">
                            {t(title.id, title.format)}
                        </Typography>
                        <Button autoFocus color="inherit" onClick={onClose}>
                            X
                        </Button>
                    </Toolbar>
                </AppBar>

                <List>
                    {values.map((item, index) => (
                        <div key={index}>
                            <ListItem>
                                <ListItemButton
                                    onClick={() =>
                                        onSelect && onSelect(item[1])
                                    }
                                >
                                    <ListItemText
                                        primary={item[0]}
                                        secondary={
                                            <>
                                                {item[2]
                                                    ? item[2]
                                                          .split("\n")
                                                          .map(
                                                              (item, index) => (
                                                                  <span
                                                                      key={
                                                                          index
                                                                      }
                                                                  >
                                                                      {item}
                                                                      <br />
                                                                  </span>
                                                              )
                                                          )
                                                    : item[1]
                                                          .split("\n")
                                                          .map(
                                                              (item, index) => (
                                                                  <span
                                                                      key={
                                                                          index
                                                                      }
                                                                  >
                                                                      {item}
                                                                      <br />
                                                                  </span>
                                                              )
                                                          )}
                                            </>
                                        }
                                    />
                                </ListItemButton>
                            </ListItem>
                            <Divider />
                        </div>
                    ))}
                </List>
            </Dialog>
        );
    }
}

export default withTranslation()(ModalDialog);
