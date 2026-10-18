import { ThemeProvider } from '@mui/material/styles';
import {
    LocalizationProvider,
    MobileDatePicker,
    MobileDateTimePicker,
    renderTimeViewClock
} from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFnsV3';
import * as XDatePickers from '@mui/x-date-pickers/locales';
import * as DateFnsLang from 'date-fns/locale';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { getDateString } from '../../helpers/utils/getDateString';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { selectThemeMode, useThemeStore } from '../../stores/theme';
import { createAppTheme } from '../../styles/timePickerTheme';

export interface TimePickerProps {
    readonly id?: string;
    readonly value?: number;
    readonly defaultValue?: number;
    readonly placeholder?: string;
    readonly onChange: (value: number) => void;
    readonly className?: string;
    readonly currentLocale: string;
    readonly dateOnly?: boolean;
}

export const TimePicker = ({
    id,
    onChange,
    value,
    placeholder,
    defaultValue,
    currentLocale,
    className,
    dateOnly = false
}: TimePickerProps) => {
    const pickerRef = useRef<HTMLDivElement>(null);
    const [open, setOpen] = useState(false);
    const mode = useThemeStore(selectThemeMode);
    const safeCurrentLocale = currentLocale || 'en-US';

    const themeRecords = useMemo(
        () =>
            Object.entries(XDatePickers).reduce(
                (acc, [locale, value]) => {
                    acc[locale] = value;
                    return acc;
                },
                {} as Record<string, object>
            ),
        []
    );
    const adapterLocaleRecords = useMemo(
        () =>
            Object.entries(DateFnsLang).reduce(
                (acc, [locale, value]) => {
                    acc[locale] = value;
                    return acc;
                },
                {} as Record<string, object>
            ),
        []
    );

    const [locale4Component, setLocale4Component] = useState('enUS');
    useEffect(() => {
        const componentLocale = safeCurrentLocale.replace(/[^a-z0-9]/gi, '');
        setLocale4Component(themeRecords[componentLocale] ? componentLocale : 'enUS');
    }, [safeCurrentLocale, themeRecords]);

    const [locale4Adapter, setLocale4Adapter] = useState('enUS');
    useEffect(() => {
        const componentLocale = safeCurrentLocale.replace(/[^a-z0-9]/gi, '');
        setLocale4Adapter(adapterLocaleRecords[componentLocale] ? componentLocale : 'enUS');
    }, [safeCurrentLocale, adapterLocaleRecords]);

    const [internalValue, setInternalValue] = useState<number | null>();
    useEffect(() => {
        setInternalValue(value ?? defaultValue ?? null);
    }, [value, defaultValue]);

    const handleDateChange = useCallback(
        (newValue: Date | null) => {
            const newValueUnixMillis = newValue?.getTime() ?? 0;
            setInternalValue(newValueUnixMillis);
            onChange(newValueUnixMillis);
        },
        [onChange]
    );

    const theme = useMemo(
        () => createAppTheme(mode, themeRecords[locale4Component]),
        [mode, locale4Component, themeRecords]
    );
    const getDialogContainer = useCallback(
        () => pickerRef.current?.closest('dialog') ?? document.body,
        []
    );

    return (
        <div className="app-time-picker" ref={pickerRef}>
            <input
                id={id}
                readOnly
                type="text"
                placeholder={placeholder}
                className={`app-pointer ${className}`}
                onClick={() => setOpen(true)}
                onFocus={({ currentTarget }) => currentTarget.blur()}
                value={
                    internalValue === null || internalValue === undefined || internalValue === 0
                        ? ''
                        : dateOnly
                          ? getDateString(internalValue)
                          : getTimeString(internalValue)
                }
            />

            <div className="app-time-picker-anchor">
                <ThemeProvider theme={theme}>
                    <LocalizationProvider
                        dateAdapter={AdapterDateFns}
                        adapterLocale={adapterLocaleRecords[locale4Adapter]}
                    >
                        {dateOnly ? (
                            <MobileDatePicker
                                open={open}
                                onClose={() => setOpen(false)}
                                onChange={handleDateChange}
                                format="yyyy-MM-dd"
                                className="app-full-width"
                                timezone="system"
                                views={['year', 'month', 'day']}
                                slotProps={{
                                    dialog: { container: getDialogContainer },
                                    field: { clearable: true },
                                    mobilePaper: {
                                        className: 'app-time-picker-paper'
                                    }
                                }}
                                defaultValue={
                                    defaultValue === null || defaultValue === undefined
                                        ? null
                                        : new Date(defaultValue)
                                }
                                value={
                                    internalValue === null || internalValue === undefined
                                        ? null
                                        : new Date(internalValue)
                                }
                            />
                        ) : (
                            <MobileDateTimePicker
                                open={open}
                                onClose={() => setOpen(false)}
                                onChange={handleDateChange}
                                format="yyyy-MM-dd HH:mm:ss"
                                className="app-full-width"
                                timezone="system"
                                views={['year', 'month', 'day', 'hours', 'minutes', 'seconds']}
                                viewRenderers={{
                                    hours: renderTimeViewClock,
                                    minutes: renderTimeViewClock,
                                    seconds: renderTimeViewClock
                                }}
                                slotProps={{
                                    dialog: { container: getDialogContainer },
                                    field: { clearable: true },
                                    mobilePaper: {
                                        className: 'app-time-picker-paper'
                                    }
                                }}
                                defaultValue={defaultValue ? new Date(defaultValue) : null}
                                value={internalValue ? new Date(internalValue) : null}
                                ampm={false}
                            />
                        )}
                    </LocalizationProvider>
                </ThemeProvider>
            </div>
        </div>
    );
};
