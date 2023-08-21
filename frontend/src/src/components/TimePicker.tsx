import { Component } from "react";
import {
    DateTimePicker,
    LocalizationProvider,
    renderTimeViewClock,
} from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import { zhCN } from "@mui/x-date-pickers/locales";
import zh from "date-fns/locale/zh-CN";

export interface TimePickerProps {
    readonly label: string;
    readonly value?: number;
    readonly defaultValue?: number;
    readonly onChange: (value: number) => void;
}

export default class TimePicker extends Component<TimePickerProps> {
    render() {
        const { label, onChange, value, defaultValue } = this.props;
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        const theme = createTheme({}, zhCN);

        return (
            <ThemeProvider theme={theme}>
                <LocalizationProvider
                    dateAdapter={AdapterDateFns}
                    adapterLocale={zh}
                >
                    <DateTimePicker
                        className="w-full"
                        timezone="system"
                        views={[
                            "year",
                            "month",
                            "day",
                            "hours",
                            "minutes",
                            "seconds",
                        ]}
                        viewRenderers={{
                            hours: renderTimeViewClock,
                            minutes: renderTimeViewClock,
                            seconds: renderTimeViewClock,
                        }}
                        onChange={(value) => {
                            const val = value?.valueOf();
                            onChange(val as number);
                        }}
                        label={`${label}（时区 ${timezone}）`}
                        defaultValue={defaultValue}
                        value={value}
                        ampm={false}
                    />
                </LocalizationProvider>
            </ThemeProvider>
        );
    }
}
