import { Component } from "react";
import {
    DateTimePicker,
    LocalizationProvider,
    renderTimeViewClock,
} from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import * as XDatePickers from "@mui/x-date-pickers/locales";
import * as DateFnsLang from "date-fns/locale";
import getLanguage from "../helpers/i18n/getLanguage";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../config/i18n";

export interface TimePickerProps {
    readonly value?: number;
    readonly defaultValue?: number;
    readonly label: I18nTranslation;
    readonly onChange: (value: number) => void;
}

const languageConfig: Record<string, { theme: any; adapterLocale: Locale }> = {
    "zh-CN": { theme: XDatePickers.zhCN, adapterLocale: DateFnsLang.zhCN },
    "zh-TW": { theme: XDatePickers.zhHK, adapterLocale: DateFnsLang.zhTW },
    "en-US": { theme: XDatePickers.enUS, adapterLocale: DateFnsLang.enUS },
    "ja-JP": { theme: XDatePickers.jaJP, adapterLocale: DateFnsLang.ja },
    "ko-KR": { theme: XDatePickers.koKR, adapterLocale: DateFnsLang.ko },
};

class TimePicker extends Component<TimePickerProps & WithTranslation> {
    render() {
        const { t, label, onChange, value, defaultValue } = this.props;
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

        const i18n = getLanguage();
        const theme = createTheme({}, languageConfig[i18n].theme);

        return (
            <ThemeProvider theme={theme}>
                <LocalizationProvider
                    dateAdapter={AdapterDateFns}
                    adapterLocale={languageConfig[i18n].adapterLocale}
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
                            const val = value?.valueOf();
                            onChange(val as number);
                        }}
                        slotProps={{
                            field: { clearable: true },
                        }}
                        label={`${t(label.id, label.format)} - ${timezone}`}
                        defaultValue={defaultValue}
                        value={value}
                        ampm={false}
                    />
                </LocalizationProvider>
            </ThemeProvider>
        );
    }
}

export default withTranslation()(TimePicker);
