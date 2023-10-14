import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import Typography from "@mui/material/Typography";
import { Component } from "react";

export interface ProgressProps {
    readonly precision?: number;
    readonly label?: string;
    value: number;
}

export default class Progress extends Component<ProgressProps> {
    render() {
        const { value, label, precision } = this.props;
        return (
            <Box sx={{ display: "flex", alignItems: "center" }}>
                <Box sx={{ width: "100%", my: 1, mx: 2 }}>
                    <LinearProgress
                        className="rounded-lg"
                        variant="determinate"
                        color="secondary"
                        value={value}
                    />
                </Box>
                <Box sx={{ minWidth: 100 }}>
                    <Typography
                        className="overflow-scroll"
                        color="text.secondary"
                        variant="body2"
                    >{`[${value.toFixed(
                        precision || 2
                    )}%] ${label}`}</Typography>
                </Box>
            </Box>
        );
    }
}
