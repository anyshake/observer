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

const languageConfig: Record<string, { theme: object; adapterLocale: Locale }> =
    {
        "zh-CN": { theme: XDatePickers.zhCN, adapterLocale: DateFnsLang.zhCN },
        "zh-TW": { theme: XDatePickers.zhHK, adapterLocale: DateFnsLang.zhTW },
        "en-US": { theme: XDatePickers.enUS, adapterLocale: DateFnsLang.enUS },
        "ja-JP": { theme: XDatePickers.jaJP, adapterLocale: DateFnsLang.ja },
        "ko-KR": { theme: XDatePickers.koKR, adapterLocale: DateFnsLang.ko },
    };

export const TimePicker = (props: TimePickerProps) => {
    const { label, onChange, value, defaultValue, currentLocale } = props;

    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    const theme = createTheme({}, languageConfig[currentLocale].theme);

    return (
        <ThemeProvider theme={theme}>
            <LocalizationProvider
                dateAdapter={AdapterDateFns}
                adapterLocale={languageConfig[currentLocale].adapterLocale}
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
                        let val = value?.valueOf() ?? 0;
                        if (isNaN(val)) {
                            val = 0;
                        }
                        onChange(val);
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
