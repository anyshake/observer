import { createTheme, ThemeProvider } from '@mui/material/styles';
import {
    LocalizationProvider,
    MobileDateTimePicker,
    renderTimeViewClock
} from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFnsV3';
import * as XDatePickers from '@mui/x-date-pickers/locales';
import * as DateFnsLang from 'date-fns/locale';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { localeConfig } from '../config/locale';
import { getTimeString } from '../helpers/utils/getTimeString';

export interface TimePickerProps {
    readonly value?: number;
    readonly defaultValue?: number;
    readonly placeholder?: string;
    readonly onChange: (value: number) => void;
    readonly className?: string;
    readonly currentLocale: keyof typeof localeConfig.resources;
}

export const TimePicker = ({
    onChange,
    value,
    placeholder,
    defaultValue,
    currentLocale,
    className
}: TimePickerProps) => {
    const [open, setOpen] = useState(false);

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
        const componentLocale = currentLocale.replace(/[^a-z0-9]/gi, '');
        setLocale4Component(themeRecords[componentLocale] ? componentLocale : 'enUS');
    }, [currentLocale, themeRecords]);

    const [locale4Adapter, setLocale4Adapter] = useState('en-US');
    useEffect(() => {
        const componentLocale = currentLocale.replace(/[^a-z0-9]/gi, '');
        setLocale4Adapter(adapterLocaleRecords[componentLocale] ? componentLocale : 'en-US');
    }, [currentLocale, adapterLocaleRecords]);

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
        () =>
            createTheme(
                {
                    palette: {
                        primary: { main: '#8b3dff' },
                        secondary: { main: '#7B1FA2' },
                        background: { default: '#F3E5F5' }
                    }
                },
                themeRecords[locale4Component]
            ),
        [locale4Component, themeRecords]
    );

    return (
        <div>
            <input
                readOnly
                type="text"
                placeholder={placeholder}
                className={`cursor-pointer ${className}`}
                onClick={() => setOpen(true)}
                onFocus={({ currentTarget }) => currentTarget.blur()}
                value={getTimeString(internalValue ?? 0)}
            />

            <div hidden={true}>
                <ThemeProvider theme={theme}>
                    <LocalizationProvider
                        dateAdapter={AdapterDateFns}
                        adapterLocale={adapterLocaleRecords[locale4Adapter]}
                    >
                        <MobileDateTimePicker
                            open={open}
                            onClose={() => setOpen(false)}
                            onChange={handleDateChange}
                            format="yyyy-MM-dd HH:mm:ss"
                            className="w-full"
                            timezone="system"
                            views={['year', 'month', 'day', 'hours', 'minutes', 'seconds']}
                            viewRenderers={{
                                hours: renderTimeViewClock,
                                minutes: renderTimeViewClock,
                                seconds: renderTimeViewClock
                            }}
                            slotProps={{
                                field: { clearable: true },
                                mobilePaper: {
                                    sx: {
                                        padding: '12px',
                                        borderRadius: '12px',
                                        overflow: 'hidden'
                                    }
                                }
                            }}
                            defaultValue={new Date(defaultValue ?? 0)}
                            value={new Date(internalValue ?? 0)}
                            ampm={false}
                        />
                    </LocalizationProvider>
                </ThemeProvider>
            </div>
        </div>
    );
};
