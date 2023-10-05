import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import Typography from "@mui/material/Typography";
import { Component } from "react";

export interface ProgressProps {
    value: number;
}

export default class Progress extends Component<ProgressProps> {
    render() {
        const { value } = this.props;
        return (
            <Box sx={{ display: "flex", alignItems: "center" }}>
                <Box sx={{ width: "100%", mr: 2 }}>
                    <LinearProgress
                        className="rounded-lg"
                        variant="determinate"
                        color="secondary"
                        value={value}
                    />
                </Box>
                <Box sx={{ minWidth: 35 }}>
                    <Typography
                        variant="body2"
                        color="text.secondary"
                    >{`${Math.round(value)}%`}</Typography>
                </Box>
            </Box>
        );
    }
}
