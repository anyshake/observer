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
import { IntensityScaleStandard } from "../../helpers/getIntensity";
import getLocalStorage from "../../helpers/getLocalStorage";
import Button from "../../components/Button";
import SelectDialog, { SelectDialogProps } from "../../components/SelectDialog";
import setLocalStorage from "../../helpers/setLocalStorage";

interface State {
    readonly scale: IntensityScaleStandard;
    readonly select: SelectDialogProps;
}

export default class Setting extends Component<{}, State> {
    constructor(props: {}) {
        super(props);
        this.state = {
            scale: "JMA",
            select: {
                title: "选择烈度标准",
                open: false,
                values: [
                    ["日本気象庁震度", "JMA"],
                    ["中国地震局地震烈度", "CSIS"],
                    ["修訂麥加利地震烈度表", "MMI"],
                    ["台湾中央氣象局新震度", "CWB"],
                ],
            },
        };
    }

    componentDidMount(): void {
        const scale = getLocalStorage("scale", "JMA") as IntensityScaleStandard;
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
        setLocalStorage("scale", value);
        this.setState({
            scale: value as IntensityScaleStandard,
            select: { ...this.state.select, open: false },
        });
        toast.success(`已切换烈度标准为 ${value}`);
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
                        <Card className="h-[200px]" label="烈度标准">
                            {`当前烈度标准为 ${scale}`}
                            <Text>
                                烈度标准是用来衡量地震烈度的标准，不同的标准会导致不同的烈度值。
                            </Text>
                            <Button
                                className="bg-lime-700 hover:bg-lime-800"
                                onClick={this.handleSelectScale}
                                label="点击选择"
                            />
                        </Card>

                        <Card className="h-[200px]" label="重置应用">
                            <Text>
                                当应用出现问题时，可以尝试在此重置应用偏好。
                            </Text>
                            <Text>
                                在重置应用后，应用将会恢复到初始状态，但是不会对后端数据造成影响。
                            </Text>
                            <Button
                                className="bg-rose-700 hover:bg-rose-800"
                                onClick={this.handlePurgePref}
                                label="点击清理"
                            />
                        </Card>
                    </Container>

                    <Card label="开源许可">
                        <Text className="font-medium">
                            版权所有 (C) Project ES
                        </Text>
                        <Text>
                            许可证持有人特此授予任何获得本软件和相关文档文件（以下简称「软件」）副本的人免费许可，以无限制地处理本软件的权利，包括但不限于使用、复制、修改、合并、发布、分发、再许可和
                            /
                            或销售本软件的副本，并允许提供本软件的人员这样做，但须满足以下条件：
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
                        <Text>
                            MIT
                            开源许可协议是一种宽松的许可协议，允许您自由使用、修改和分发您的软件，只要您包含原始许可和版权声明。
                        </Text>
                        <Text>
                            请注意，MIT
                            许可协议不提供任何担保，并且作者或版权持有者不对软件的使用产生的任何索赔或损害负责。如果您决定使用本许可协议，请确保在您的项目中包含适当的许可和版权声明，以保护您的权益。谢谢您选择开源，并祝您的项目也取得成功！
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
