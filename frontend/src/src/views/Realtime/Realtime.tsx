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
import { fallbackScale } from "../../config/global";
import { ReduxStoreProps } from "../../config/store";
import { connect } from "react-redux";
import { update as updateADC } from "../../store/adc";
import { update as updateGeophone } from "../../store/geophone";
import mapStateToProps from "../../helpers/utils/mapStateToProps";
import { WithTranslation, withTranslation } from "react-i18next";

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

class Realtime extends Component<
    ReduxStoreProps & WithTranslation,
    RealtimeState
> {
    prevTs: number;
    websocket: WebSocket | null | {};
    constructor(props: ReduxStoreProps & WithTranslation) {
        super(props);
        this.state = {
            banner: {
                type: "warning",
                label: { id: "views.realtime.banner.warning.label" },
                text: { id: "views.realtime.banner.warning.text" },
            },
            areas: [
                {
                    tag: "ehz",
                    area: {
                        label: { id: "views.realtime.areas.ehz.label" },
                        text: {
                            id: "views.realtime.areas.ehz.text",
                            format: {
                                pga: "0.00000",
                                pgv: "0.00000",
                                intensity: "Unknown",
                            },
                        },
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
                        label: { id: "views.realtime.areas.ehe.label" },
                        text: {
                            id: "views.realtime.areas.ehe.text",
                            format: {
                                pga: "0.00000",
                                pgv: "0.00000",
                                intensity: "Unknown",
                            },
                        },
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
                        label: { id: "views.realtime.areas.ehn.label" },
                        text: {
                            id: "views.realtime.areas.ehn.text",
                            format: {
                                pga: "0.00000",
                                pgv: "0.00000",
                                intensity: "Unknown",
                            },
                        },
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
                ehz: 1,
                ehe: 1,
                ehn: 1,
            },
            adc: {
                fullscale: 1,
                resolution: 1,
            },
            scale: fallbackScale.property(),
        };
        // Some initializations
        this.websocket = null;
        this.prevTs = 0;
    }

    // WebSocket OnOpen handler
    handleWebsocketOpen = (): void => {
        // Display success message after connection established
        setTimeout(() => {
            const { t } = this.props;
            toast.success(t("views.realtime.toasts.websocket_connected"));
        }, 500);
    };

    // WebSocket OnClose handler
    handleWebsocketClose = (): void => {
        // Reconnect to server if websocket closed
        // Unless leaving the component, the ref will be {}
        // So we can use instanceof to check if still in component
        if (this.websocket && this.websocket instanceof WebSocket) {
            // setBanner returns error when no arguments passed
            const banner = setBanner();
            this.setState({ banner });
            // Reconnect to server
            this.websocket = websocketByTag({
                tag: "socket",
                onData: this.handleWebsocketData,
                onOpen: this.handleWebsocketOpen,
                onClose: this.handleWebsocketClose,
            }) as WebSocket;
        }
    };

    // WebSocket OnData handler
    handleWebsocketData = (event: MessageEvent): void => {
        const jsonData = JSON.parse(event.data);
        const { adc, geophone, scale } = this.state;
        const banner = setBanner(jsonData, this.prevTs, scale);
        // Get waveform retention time from global config
        const { retention } = this.props.retention;
        const areas = setAreas(
            this.state.areas,
            jsonData,
            this.prevTs,
            retention,
            adc,
            geophone,
            scale
        );

        this.prevTs = jsonData.ts;
        this.setState({ banner, areas });
    };

    async componentDidMount(): Promise<void> {
        // Get ADC, Geophone, scale standard from Redux
        let { adc } = this.props.adc;
        const { scale } = this.props.scale;
        let { geophone } = this.props.geophone;
        const { ehz, ehe, ehn } = geophone;
        const { resolution } = adc;

        // Query again from server if value is not set
        if (resolution === -1 || ehz * ehe * ehn === 0) {
            const res = await restfulApiByTag({
                tag: "station",
            });
            if (res.data) {
                // Get new formatted state
                geophone = setGeophone(res);
                adc = setADC(res);
                // Apply to Redux store
                const { updateADC, updateGeophone } = this.props;
                updateGeophone && updateGeophone(geophone);
                updateADC && updateADC(adc);
            } else {
                // Show error and return if failed
                const { t } = this.props;
                toast.error(t("views.realtime.toasts.fetch_metadata_error"));
                return;
            }
        }

        // Apply to component state
        this.setState({ adc, geophone, scale });
        // Establish websocket connection
        this.websocket = websocketByTag({
            tag: "socket",
            onData: this.handleWebsocketData,
            onOpen: this.handleWebsocketOpen,
            onClose: this.handleWebsocketClose,
        }) as WebSocket;
    }

    componentWillUnmount(): void {
        // Close websocket connection when leaving
        if (this.websocket instanceof WebSocket) {
            this.websocket?.close();
            this.websocket = {};
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
                                    tickPrecision={1}
                                    tickInterval={10}
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

export default connect(mapStateToProps, {
    updateGeophone,
    updateADC,
})(withTranslation()(Realtime));
