import {
    DateTimePicker,
    LocalizationProvider,
    renderTimeViewClock,
} from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import * as XDatePickers from "@mui/x-date-pickers/locales";
import * as DateFnsLang from "date-fns/locale";

export interface TimePickerProps {
    readonly label: string;
    readonly value?: number;
    readonly defaultValue?: number;
    readonly currentLocale: string;
    readonly onChange: (value: number) => void;
}

export const TimePicker = (props: TimePickerProps) => {
    const { label, onChange, value, defaultValue, currentLocale } = props;

    const themeRecords = Object.entries(XDatePickers).reduce(
        (acc, [locale, value]) => {
            acc[locale] = value;
            return acc;
        },
        {} as Record<string, any>
    );
    const adapterLocaleRecords = Object.entries(DateFnsLang).reduce(
        (acc, [locale, value]) => {
            acc[locale] = value;
            return acc;
        },
        {} as Record<string, any>
    );

    let locale4Component = currentLocale.replaceAll(/[^a-z0-9]/gi, "");
    if (!themeRecords[locale4Component] || !adapterLocaleRecords[locale4Component]) {
        locale4Component = "enUS"
    }

    const theme = createTheme({}, themeRecords[locale4Component]);
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    return (
        <ThemeProvider theme={theme}>
            <LocalizationProvider
                dateAdapter={AdapterDateFns}
                adapterLocale={adapterLocaleRecords[locale4Component]}
            >
                <DateTimePicker
                    format="yyyy-MM-dd HH:mm:ss"
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
                        const dateObj = new Date(value ?? 0);
                        const dateTs = dateObj.getTime();
                        onChange(isNaN(dateTs) ? 0 : dateTs);
                    }}
                    slotProps={{
                        field: { clearable: true },
                    }}
                    label={`${label} - ${timezone}`}
                    defaultValue={defaultValue}
                    value={value ?? 0}
                    ampm={false}
                />
            </LocalizationProvider>
        </ThemeProvider>
    );
};
