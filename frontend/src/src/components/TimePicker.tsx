import { Component } from "react";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import {
    DateTimePicker,
    LocalizationProvider,
    renderTimeViewClock,
} from "@mui/x-date-pickers";
import dayjs from "dayjs";

interface TimePickerProps {
    readonly label: string;
    readonly value?: number;
    readonly defaultValue?: number;
    readonly onChange: (value: number) => void;
}

export default class TimePicker extends Component<TimePickerProps> {
    render() {
        const { label, onChange, value, defaultValue } = this.props;
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

        return (
            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <DateTimePicker
                    className="w-full"
                    timezone="default"
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
                    {...(defaultValue && {
                        defaultValue: dayjs(value),
                    })}
                    {...(value && {
                        value: dayjs(value),
                    })}
                />
            </LocalizationProvider>
        );
    }
}
