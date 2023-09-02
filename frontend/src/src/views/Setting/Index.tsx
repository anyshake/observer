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
import getLocalStorage from "../../helpers/storage/getLocalStorage";
import Button from "../../components/Button";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import setLocalStorage from "../../helpers/storage/setLocalStorage";
import GLOBAL_CONFIG, { fallbackScale } from "../../config/global";

interface SettingState {
    readonly select: SelectDialogProps;
    readonly scale: IntensityStandardProperty;
}

export default class Setting extends Component<{}, SettingState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            select: {
                open: false,
                title: "选择震度标准",
                values: GLOBAL_CONFIG.app_settings.scales.map((item) => [
                    item.property().name,
                    item.property().value,
                ]),
            },
            scale: fallbackScale.property(),
        };
    }

    componentDidMount(): void {
        const scale = getLocalStorage(
            "scale",
            fallbackScale.property(),
            true
        ) as IntensityStandardProperty;
        this.setState({ scale });
    }

    handlePurgePref = (): void => {
        localStorage.clear();
        toast.success("已清理应用偏好，页面即将刷新");
        setTimeout(() => window.location.reload(), 1000);
    };

    handleSelectScale = (): void => {
        this.setState((state) => ({
            select: {
                ...state.select,
                open: true,
            },
        }));
    };

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

        setLocalStorage("scale", scaleStandard, true);
        toast.success(`已切换震度标准为 ${value}`);
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
                        <Card className="h-[200px]" label="震度标准">
                            {`当前震度标准为${scale.name}。`}
                            <Text>
                                震度标准是用来衡量地震震度的标准，不同的标准会导致不同的震度值。
                            </Text>
                            <Button
                                className="bg-lime-700 hover:bg-lime-800"
                                onClick={this.handleSelectScale}
                                label="点击选择"
                            />
                        </Card>

                        <Card className="h-[200px]" label="重置应用">
                            <Text>应用出现问题时，可尝试重置应用偏好。</Text>
                            <Text>
                                执行重置后，浏览器中的偏好将被清理，不会对后端服务器产生影响。
                            </Text>
                            <Button
                                className="bg-rose-700 hover:bg-rose-800"
                                onClick={this.handlePurgePref}
                                label="点击清理"
                            />
                        </Card>
                    </Container>

                    <Card label="关于软件">
                        <Text>
                            感谢使用本软件！本项目隶属于 Project ES
                            项目中的一部分，由 github.com/bclswl0827
                            主导开发，并遵循 MIT 协议开源。
                        </Text>
                        <Text>
                            许可证持有人特此授予任何获得本软件和相关文档文件副本的人免费许可，以无限制地处理本软件的权利，包括但不限于使用、复制、修改、合并、发布、分发、再许可和（或）销售本软件的副本，并允许提供本软件的人员这样做，但须满足以下条件：
                        </Text>
                        <Text className="ml-4">
                            1. Project ES
                            的开发者和版权持有者在此明确指出，本软件不附带任何明示或暗示的担保或条件，包括但不限于适销性、特定用途适用性和非侵权性。任何情况下，作者或版权持有者都不对任何索赔、损害或其他责任负责，无论是在合同诉讼、侵权行为还是其他方面的诉讼中产生的，与本软件或本软件的使用或其他处理有关。
                        </Text>
                        <Text className="ml-4">
                            2. Project ES
                            的开发者和版权持有者保留拒绝对其软件提供技术支持的权利。用户可以通过本软件自行修改和使用，但不得寻求
                            Project ES 的开发者和版权持有者的技术支持。
                        </Text>
                        <Text className="font-medium">
                            版权所有 (C) Project ES
                        </Text>
                    </Card>
                </Content>

                <Scroller />
                <Footer />

                <SelectDialog {...select} onSelect={this.handleScaleChange} />
                <Toaster position="top-center" />
            </View>
        );
    }
}
