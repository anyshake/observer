import { useEffect, useState } from "react";
import { Container } from "../../components/Container";
import { Button } from "../../components/Button";
import { Panel } from "../../components/Panel";
import { useTranslation } from "react-i18next";
import { sendUserAlert } from "../../helpers/interact/sendUserAlert";
import { sendUserConfirm } from "../../helpers/interact/sendUserConfirm";
import { useSelector, useDispatch } from "react-redux";
import { ReduxStoreProps } from "../../config/store";
import { globalConfig } from "../../config/global";
import { Form, FormProps } from "../../components/Form";
import { onUpdate as updateRetention } from "../../stores/retention";
import { Select, SelectProps } from "../../components/Select";
import { onUpdate as updateScale } from "../../stores/scale";
import { apiConfig } from "../../config/api";
import { requestRestApi } from "../../helpers/request/requestRestApi";
import { Code } from "../../components/Code";

const Settings = () => {
    const { t } = useTranslation();
    const dispatch = useDispatch();

    const { retention, duration, scales } = globalConfig;
    const handleReload = () => {
        setTimeout(() => {
            window.location.reload();
        }, 2500);
    };

    const [form, setForm] = useState<
        FormProps & { values?: Record<string, string | number> }
    >({ open: false, inputType: "number" });

    const handleCloseForm = () => {
        setForm({ ...form, open: false });
    };

    const handleRetentionChange = (newValue?: string) => {
        if (!!newValue?.length) {
            const value = parseInt(newValue);
            const { maximum, minimum } = retention;
            if (isNaN(value) || value < minimum || value > maximum) {
                return;
            }
            sendUserAlert(
                t("views.setting.toasts.retention_set", { current: value })
            );
            dispatch(updateRetention(value));
            handleCloseForm();
            handleReload();
        } else {
            setForm({
                ...form,
                open: true,
                values: { ...retention },
                onSubmit: handleRetentionChange,
                cancelText: "views.setting.forms.waveform_retention.cancel",
                submitText: "views.setting.forms.waveform_retention.submit",
                title: "views.setting.forms.waveform_retention.title",
                content: "views.setting.forms.waveform_retention.content",
                placeholder:
                    "views.setting.forms.waveform_retention.placeholder",
            });
        }
    };

    const handleDurationChange = (newValue?: string) => {
        if (!!newValue?.length) {
            const value = parseInt(newValue);
            const { maximum, minimum } = retention;
            if (isNaN(value) || value < minimum || value > maximum) {
                return;
            }
            sendUserAlert(
                t("views.setting.toasts.duration_set", { current: value })
            );
            dispatch(updateRetention(value));
            handleCloseForm();
            handleReload();
        } else {
            setForm({
                ...form,
                open: true,
                values: { ...duration },
                onSubmit: handleRetentionChange,
                cancelText: "views.setting.forms.query_duration.cancel",
                submitText: "views.setting.forms.query_duration.submit",
                title: "views.setting.forms.query_duration.title",
                content: "views.setting.forms.query_duration.content",
                placeholder: "views.setting.forms.query_duration.placeholder",
            });
        }
    };

    const [select, setSelect] = useState<SelectProps>({ open: false });

    const handleCloseSelect = () => {
        setSelect({ ...select, open: false });
    };

    const handleSelectScale = (newValue?: string) => {
        if (!!newValue?.length) {
            const isNewValueInScales = scales.some(
                (item) => item.property().value === newValue
            );
            if (isNewValueInScales) {
                const newScaleName = scales
                    .find((item) => item.property().value === newValue)
                    ?.property().name;
                sendUserAlert(
                    t("views.setting.toasts.scale_changed", {
                        scale: newScaleName,
                    })
                );
                dispatch(updateScale(newValue));
                handleCloseSelect();
                handleReload();
            }
        } else {
            setSelect({
                ...select,
                open: true,
                onSelect: handleSelectScale,
                title: "views.setting.selects.choose_scale.title",
                options: scales.map(({ property }) => {
                    const { name, value } = property();
                    return [name, value];
                }),
            });
        }
    };

    const handlePurgeCache = () => {
        sendUserConfirm(t("views.setting.toasts.confirm_purge"), {
            title: t("views.setting.toasts.confirm_title"),
            confirmText: t("views.setting.toasts.confirm_button"),
            cancelText: t("views.setting.toasts.cancel_button"),
            onConfirmed: () => {
                sendUserAlert(t("views.setting.toasts.cache_purged"));
                localStorage.clear();
                handleReload();
            },
        });
    };

    const { retention: currentRetention } = useSelector(
        ({ retention }: ReduxStoreProps) => retention
    );
    const { duration: currentDuration } = useSelector(
        ({ duration }: ReduxStoreProps) => duration
    );
    const { scale: currentScale } = useSelector(
        ({ scale }: ReduxStoreProps) => scale
    );
    const scaleName =
        scales
            .find((item) => item.property().value === currentScale)
            ?.property().name || "Unknown";

    const [panels] = useState([
        {
            label: "views.setting.panels.waveform_retention",
            content: "views.setting.contents.waveform_retention",
            button: "views.setting.buttons.waveform_retention",
            className: "bg-teal-700 hover:bg-teal-800",
            onClick: handleRetentionChange,
            values: { current: currentRetention, ...retention },
        },
        {
            label: "views.setting.panels.query_duration",
            content: "views.setting.contents.query_duration",
            button: "views.setting.buttons.query_duration",
            className: "bg-lime-700 hover:bg-lime-800",
            onClick: handleDurationChange,
            values: { current: currentDuration, ...duration },
        },
        {
            label: "views.setting.panels.select_scale",
            button: "views.setting.buttons.select_scale",
            className: "bg-sky-700 hover:bg-sky-800",
            content: "views.setting.contents.select_scale",
            onClick: handleSelectScale,
            values: { scale: scaleName },
        },
        {
            label: "views.setting.panels.purge_cache",
            content: "views.setting.contents.purge_cache",
            button: "views.setting.buttons.purge_cache",
            className: "bg-pink-700 hover:bg-pink-800",
            onClick: handlePurgeCache,
        },
    ]);

    const [stationInventory, setStationInventory] = useState<string>();

    const getStationInventory = async () => {
        const { backend, endpoints } = apiConfig;
        const payload = { format: "json" };
        const res = await requestRestApi({
            backend,
            payload,
            timeout: 30,
            endpoint: endpoints.inventory,
        });
        if (!res?.data) {
            return;
        }
        setStationInventory(res.data);
    };

    useEffect(() => {
        void getStationInventory();
    }, []);

    return (
        <>
            <Container className="gap-4 grid md:grid-cols-2">
                {panels.map(
                    ({
                        label,
                        content,
                        button,
                        className,
                        onClick,
                        values,
                    }) => (
                        <Panel key={label} className="" label={t(label)}>
                            {t(content, { ...values })
                                .split("\n")
                                .map((item) => (
                                    <div key={item}>{item}</div>
                                ))}
                            <Button
                                label={t(button)}
                                onClick={onClick}
                                className={className}
                            />
                        </Panel>
                    )
                )}
                <Form
                    {...form}
                    onClose={handleCloseForm}
                    title={t(form.title ?? "")}
                    cancelText={t(form.cancelText ?? "")}
                    submitText={t(form.submitText ?? "")}
                    placeholder={t(form.placeholder ?? "")}
                    content={t(form.content ?? "", { ...form.values })}
                />
                <Select
                    {...select}
                    onClose={handleCloseSelect}
                    title={t(select.title ?? "")}
                />
            </Container>
            {!!stationInventory?.length && (
                <Panel label={t("views.setting.panels.station_inventory")}>
                    <Code language="xml" fileName="inventory.xml">
                        {stationInventory}
                    </Code>
                </Panel>
            )}
        </>
    );
};

export default Settings;
