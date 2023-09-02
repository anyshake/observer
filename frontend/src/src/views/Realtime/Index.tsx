import { Component } from "react";
import Header from "../../components/Header";
import Sidebar from "../../components/Sidebar";
import Content from "../../components/Content";
import Navbar from "../../components/Navbar";
import View from "../../components/View";
import Scroller from "../../components/Scroller";
import Banner, { BannerProps } from "../../components/Banner";
import Footer from "../../components/Footer";
import Area, { AreaProps } from "../../components/Area";
import Container from "../../components/Container";
import Chart, { ChartProps } from "../../components/Chart";
import restfulApiByTag from "../../helpers/request/restfulApiByTag";
import setADC from "./setADC";
import setGeophone from "./setGeophone";
import websocketByTag from "../../helpers/request/websocketByTag";
import toast, { Toaster } from "react-hot-toast";
import setBanner from "./setBanner";
import setAreas from "./setAreas";
import { Geophone } from "../../config/geophone";
import { ADC } from "../../config/adc";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import getLocalStorage from "../../helpers/storage/getLocalStorage";
import { fallbackScale } from "../../config/global";
import { ReduxStore, ReduxStoreProps } from "../../config/store";
import { connect } from "react-redux";
import { update as updateADC } from "../../store/adc";
import { update as updateGeophone } from "../../store/geophone";

// 180s by default
const QUENE_LENGTH = 180;

export interface RealtimeArea {
    readonly tag: string;
    readonly area: AreaProps;
    readonly chart: ChartProps;
}

interface RealtimeState {
    readonly adc: ADC;
    readonly banner: BannerProps;
    readonly areas: RealtimeArea[];
    readonly geophone: Geophone;
    readonly scale: IntensityStandardProperty;
}

class Realtime extends Component<ReduxStoreProps, RealtimeState> {
    prevTs: number;
    websocket: WebSocket | null;
    constructor(props: ReduxStoreProps) {
        super(props);
        this.state = {
            banner: {
                type: "warning",
                label: "正在连接服务器",
                text: "请稍等...",
            },
            areas: [
                {
                    tag: "ehz",
                    area: {
                        label: "EHZ 通道加速度",
                        text: "PGA：正在获取中\nPGV：正在获取中\n震度：正在获取中",
                    },
                    chart: {
                        backgroundColor: "#d97706",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHZ",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
                {
                    tag: "ehe",
                    area: {
                        label: "EHE 通道加速度",
                        text: "PGA：正在获取中\nPGV：正在获取中\n震度：正在获取中",
                    },
                    chart: {
                        backgroundColor: "#10b981",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHE",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
                {
                    tag: "ehn",
                    area: {
                        label: "EHN 通道加速度",
                        text: "PGA：正在获取中\nPGV：正在获取中\n震度：正在获取中",
                    },
                    chart: {
                        backgroundColor: "#0ea5e9",
                        lineWidth: 1,
                        height: 300,
                        series: {
                            name: "EHN",
                            type: "line",
                            color: "#f1f5f9",
                            data: [],
                        },
                    },
                },
            ],
            geophone: {
                ehz: 0.288,
                ehe: 0.288,
                ehn: 0.288,
            },
            adc: {
                fullscale: 5,
                resolution: 24,
            },
            scale: fallbackScale.property(),
        };
        this.prevTs = 0;
        this.websocket = null;
    }

    handleWebsocketOpen = (): void => {
        setTimeout(() => {
            toast.success("Websocket 连线已经建立");
        }, 500);
    };

    handleWebsocketClose = (): void => {
        // Display error message
        const banner = setBanner();
        this.setState({ banner });
        // Reconnect to server
        this.websocket = websocketByTag({
            tag: "socket",
            onData: this.handleWebsocketData,
            onOpen: this.handleWebsocketOpen,
            onClose: this.handleWebsocketClose,
        }) as WebSocket;
    };

    handleWebsocketData = (event: MessageEvent): void => {
        const { adc, geophone, scale } = this.state;
        const jsonData = JSON.parse(event.data);
        const banner = setBanner(jsonData, this.prevTs, scale);
        const areas = setAreas(
            this.state.areas,
            jsonData,
            this.prevTs,
            QUENE_LENGTH,
            adc,
            geophone,
            scale
        );

        this.prevTs = jsonData.ts;
        this.setState({ banner, areas });
    };

    async componentDidMount(): Promise<void> {
        // Get scale standard from localStorage or fallback
        const scale = getLocalStorage(
            "scale",
            fallbackScale.property(),
            true
        ) as IntensityStandardProperty;
        // Get ADC & Geophone parameters from redux
        let { adc } = this.props.adc;
        let { geophone } = this.props.geophone;
        // Query ADC & Geophone parameters from server
        const { resolution } = adc;
        const { ehz, ehe, ehn } = geophone;
        if (resolution === -1 || ehz * ehe * ehn === 0) {
            const res = await restfulApiByTag({
                tag: "station",
            });
            if (res.data) {
                // Get new state
                adc = setADC(res);
                geophone = setGeophone(res);
                // Update redux
                const { updateADC, updateGeophone } = this.props;
                updateGeophone(geophone);
                updateADC(adc);
            } else {
                const error = "取得测站资讯时发生错误，功能无法使用";
                toast.error(error);
                return Promise.reject(error);
            }
        }
        // Update state
        this.setState({ adc, geophone, scale });
        // Dail WebSocket
        this.websocket = websocketByTag({
            tag: "socket",
            onData: this.handleWebsocketData,
            onOpen: this.handleWebsocketOpen,
            onClose: this.handleWebsocketClose,
        }) as WebSocket;
    }

    componentWillUnmount(): void {
        // Close websocket
        if (this.websocket) {
            this.websocket.close();
            this.websocket = null;
        }
    }

    render() {
        const { areas, banner } = this.state;
        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />
                    <Banner {...banner} />
                    <Container layout="none">
                        {areas.map(({ area, chart }, index) => (
                            <Area key={index} {...area}>
                                <Chart
                                    {...chart}
                                    tooltip={true}
                                    zooming={true}
                                    animation={false}
                                    tickPrecision={0.4}
                                    tickInterval={0.0001}
                                />
                            </Area>
                        ))}
                    </Container>
                </Content>

                <Scroller />
                <Footer />
                <Toaster position="top-center" />
            </View>
        );
    }
}

const mapStateToProps = (state: ReduxStore) => {
    return { ...state };
};

export default connect(mapStateToProps, {
    updateGeophone,
    updateADC,
})(Realtime);
