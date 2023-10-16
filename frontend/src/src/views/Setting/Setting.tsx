import { Component } from "react";
import toast, { Toaster } from "react-hot-toast";
import Content from "../../components/Content";
import Footer from "../../components/Footer";
import Header from "../../components/Header";
import Navbar from "../../components/Navbar";
import Scroller from "../../components/Scroller";
import Sidebar from "../../components/Sidebar";
import View from "../../components/View";
import Container from "../../components/Container";
import Card from "../../components/Card";
import Text from "../../components/Text";
import Button from "../../components/Button";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import GLOBAL_CONFIG, { fallbackScale } from "../../config/global";
import mapStateToProps from "../../helpers/utils/mapStateToProps";
import { update as updateScale } from "../../store/scale";
import { update as updateRetention } from "../../store/retention";
import { update as updateDuration } from "../../store/duration";
import { connect } from "react-redux";
import { ReduxStoreProps } from "../../config/store";
import { WithTranslation, withTranslation } from "react-i18next";
import { I18nTranslation } from "../../config/i18n";
import FormDialog, { FormDialogProps } from "../../components/FormDialog";

interface SettingCard {
    readonly label: I18nTranslation;
    readonly button: I18nTranslation;
    readonly content: I18nTranslation;
    readonly className: string;
    readonly onClick: () => void;
}

interface SettingState {
    readonly cards: SettingCard[];
    readonly select: SelectDialogProps;
    readonly input: FormDialogProps;
}

class Setting extends Component<
    ReduxStoreProps & WithTranslation,
    SettingState
> {
    constructor(props: ReduxStoreProps & WithTranslation) {
        super(props);
        this.state = {
            cards: [
                {
                    label: { id: "views.setting.cards.waveform_retention" },
                    content: {
                        id: "views.setting.contents.waveform_retention",
                        format: {
                            retention:
                                this.props.retention.retention.toString(),
                        },
                    },
                    button: { id: "views.setting.buttons.waveform_retention" },
                    className: "bg-teal-700 hover:bg-teal-800",
                    onClick: this.handleRetentionChange,
                },
                {
                    label: { id: "views.setting.cards.query_duration" },
                    content: {
                        id: "views.setting.contents.query_duration",
                        format: {
                            duration: this.props.duration.duration.toString(),
                        },
                    },
                    button: { id: "views.setting.buttons.query_duration" },
                    className: "bg-lime-700 hover:bg-lime-800",
                    onClick: this.handleDurationChange,
                },
                {
                    label: { id: "views.setting.cards.select_scale" },
                    button: { id: "views.setting.buttons.select_scale" },
                    className: "bg-sky-700 hover:bg-sky-800",
                    content: {
                        id: "views.setting.contents.select_scale",
                        format: { scale: props.scale.scale.value },
                    },
                    onClick: this.handleSelectScale,
                },
                {
                    label: { id: "views.setting.cards.purge_cache" },
                    content: { id: "views.setting.contents.purge_cache" },
                    button: { id: "views.setting.buttons.purge_cache" },
                    className: "bg-pink-700 hover:bg-pink-800",
                    onClick: this.handlePurgeCache,
                },
            ],
            select: {
                open: false,
                title: { id: "views.setting.selects.choose_scale.title" },
                values: GLOBAL_CONFIG.app_settings.scales.map((item) => [
                    item.property().name,
                    item.property().value,
                ]),
            },
            input: {
                open: false,
                title: { id: "" },
                inputType: "number",
                submitText: { id: "" },
                onSubmit: this.handleInputSubmit,
            },
        };
    }

    // Handler for submitting input dialog
    handleInputSubmit = (value: string): void => {
        const { t } = this.props;
        const { tag } = this.state.input;
        const { updateRetention, updateDuration } = this.props;
        const roundedValue = Math.round(Number(value));

        // Update Redux store
        switch (tag) {
            case "retention":
                updateRetention && updateRetention(roundedValue);
                toast.success(
                    t("views.setting.toasts.retention_set", {
                        retention: roundedValue.toString(),
                    })
                );
                break;
            case "duration":
                updateDuration && updateDuration(roundedValue);
                toast.success(
                    t("views.setting.toasts.duration_set", {
                        duration: roundedValue.toString(),
                    })
                );
                break;
            default:
                return;
        }

        // Update state to close input dialog
        this.setState((state) => ({
            input: { ...state.input, open: false },
        }));
        // Reload page
        setTimeout(() => window.location.reload(), 1000);
    };

    // Handler for changing waveform retention time
    handleRetentionChange = (): void => {
        // Update state to open input dialog
        this.setState((state) => ({
            input: {
                ...state.input,
                open: true,
                tag: "retention",
                title: {
                    id: "views.setting.inputs.waveform_retention.title",
                },
                content: {
                    id: "views.setting.inputs.waveform_retention.content",
                },
                placeholder: {
                    id: "views.setting.inputs.waveform_retention.placeholder",
                },
                submitText: {
                    id: "views.setting.inputs.waveform_retention.submit",
                },
                defaultValue: this.props.retention.retention.toString(),
            },
        }));
    };

    // Handler for changing query duration
    handleDurationChange = (): void => {
        // Update state to open input dialog
        this.setState((state) => ({
            input: {
                ...state.input,
                open: true,
                tag: "duration",
                title: {
                    id: "views.setting.inputs.query_duration.title",
                },
                content: {
                    id: "views.setting.inputs.query_duration.content",
                },
                placeholder: {
                    id: "views.setting.inputs.query_duration.placeholder",
                },
                submitText: {
                    id: "views.setting.inputs.query_duration.submit",
                },
                defaultValue: this.props.duration.duration.toString(),
            },
        }));
    };

    // Handler for selecting scale
    handleSelectScale = (): void => {
        this.setState((state) => ({
            select: {
                ...state.select,
                open: true,
            },
        }));
    };

    // Hanlder for selecting scale standard
    handleScaleChange = (value: string): void => {
        // Match value with scale property
        const { scales } = GLOBAL_CONFIG.app_settings;
        const scaleStandard = scales
            .find((item) => item.property().value === value)
            ?.property();
        this.setState({
            select: { ...this.state.select, open: false },
        });
        // Apply scale option to Redux store
        const { t, updateScale } = this.props;
        updateScale && updateScale(scaleStandard || fallbackScale.property());
        toast.success(
            t("views.setting.toasts.scale_changed", { scale: value })
        );
        setTimeout(() => window.location.reload(), 1000);
    };

    // Handler for purging cache
    handlePurgeCache = (): void => {
        const { t } = this.props;
        localStorage.clear();
        toast.success(t("views.setting.toasts.cache_purged"));
        setTimeout(() => window.location.reload(), 1000);
    };

    render() {
        const { t } = this.props;
        const { select, cards, input } = this.state;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container layout="grid">
                        {cards.map(
                            (
                                { label, content, onClick, button, className },
                                index
                            ) => (
                                <Card
                                    key={index}
                                    className="min-h-[40vh]"
                                    label={label}
                                >
                                    {t(content.id, content.format)
                                        .split("\n")
                                        .map((item: string, index: number) => (
                                            <Text key={index} className="mb-3">
                                                {item}
                                            </Text>
                                        ))}
                                    <Button
                                        className={className}
                                        onClick={onClick}
                                        label={button}
                                    />
                                </Card>
                            )
                        )}
                    </Container>
                </Content>

                <Scroller />
                <Footer />

                <SelectDialog {...select} onSelect={this.handleScaleChange} />
                <FormDialog {...input} />
                <Toaster position="top-center" />
            </View>
        );
    }
}

export default connect(mapStateToProps, {
    updateScale,
    updateRetention,
    updateDuration,
})(withTranslation()(Setting));
