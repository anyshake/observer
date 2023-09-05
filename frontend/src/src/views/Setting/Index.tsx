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
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import Button from "../../components/Button";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import GLOBAL_CONFIG, { fallbackScale } from "../../config/global";
import mapStateToProps from "../../helpers/utils/mapStateToProps";
import { update as updateScale } from "../../store/scale";
import { connect } from "react-redux";
import { ReduxStoreProps } from "../../config/store";
import { WithTranslation, withTranslation } from "react-i18next";

interface SettingState {
    readonly select: SelectDialogProps;
    readonly scale: IntensityStandardProperty;
}

class Setting extends Component<
    ReduxStoreProps & WithTranslation,
    SettingState
> {
    constructor(props: ReduxStoreProps & WithTranslation) {
        super(props);
        this.state = {
            select: {
                open: false,
                title: { id: "views.setting.selects.choose_scale.title" },
                values: GLOBAL_CONFIG.app_settings.scales.map((item) => [
                    item.property().name,
                    item.property().value,
                ]),
            },
            scale: fallbackScale.property(),
        };
    }

    componentDidMount(): void {
        // Get scale from Redux store
        const { scale } = this.props.scale;
        this.setState({ scale });
    }

    // Clear localStorage
    handlePurgeCache = (): void => {
        const { t } = this.props;
        localStorage.clear();
        toast.success(t("views.setting.toasts.cache_purged"));
        setTimeout(() => window.location.reload(), 1000);
    };

    // Open select dialog
    handleSelectScale = (): void => {
        this.setState((state) => ({
            select: {
                ...state.select,
                open: true,
            },
        }));
    };

    // Hanlder for changing scale standard
    handleScaleChange = (value: string): void => {
        // Match value with scale property
        const { scales } = GLOBAL_CONFIG.app_settings;
        const scaleStandard = scales
            .find((item) => item.property().value === value)
            ?.property();
        this.setState({
            scale: scaleStandard || fallbackScale.property(),
            select: { ...this.state.select, open: false },
        });
        // Apply scale option to Redux store
        const { t, updateScale } = this.props;
        updateScale && updateScale(scaleStandard || fallbackScale.property());
        toast.success(
            t("views.setting.toasts.scale_changed", { scale: value })
        );
    };

    render() {
        const { scale, select } = this.state;
        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container layout="grid">
                        <Card
                            className="h-[200px]"
                            label={{
                                id: "views.setting.cards.select_scale",
                            }}
                        >
                            {`当前震度标准为${scale.name}。`}
                            <Text>
                                震度标准是用来衡量地震震度的标准，不同的标准会导致不同的震度值。
                            </Text>
                            <Button
                                className="bg-lime-700 hover:bg-lime-800"
                                onClick={this.handleSelectScale}
                                label={{
                                    id: "views.setting.buttons.select_scale",
                                }}
                            />
                        </Card>

                        <Card
                            className="h-[200px]"
                            label={{
                                id: "views.setting.cards.purge_cache",
                            }}
                        >
                            <Text>应用出现问题时，可尝试重置应用偏好。</Text>
                            <Text>
                                执行重置后，浏览器中的偏好将被清理，不会对后端服务器产生影响。
                            </Text>
                            <Button
                                className="bg-rose-700 hover:bg-rose-800"
                                onClick={this.handlePurgeCache}
                                label={{
                                    id: "views.setting.buttons.purge_cache",
                                }}
                            />
                        </Card>
                    </Container>
                </Content>

                <Scroller />
                <Footer />

                <SelectDialog {...select} onSelect={this.handleScaleChange} />
                <Toaster position="top-center" />
            </View>
        );
    }
}

export default connect(mapStateToProps, {
    updateScale,
})(withTranslation()(Setting));
