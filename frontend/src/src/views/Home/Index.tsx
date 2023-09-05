import { Component } from "react";
import Header from "../../components/Header";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Content from "../../components/Content";
import Banner, { BannerProps } from "../../components/Banner";
import Label, { LabelProps } from "../../components/Label";
import View from "../../components/View";
import Area, { AreaProps } from "../../components/Area";
import Chart, { ChartProps } from "../../components/Chart";
import Container from "../../components/Container";
import Footer from "../../components/Footer";
import MapBox, { MapBoxProps } from "../../components/MapBox";
import Scroller from "../../components/Scroller";
import CheckIcon from "../../assets/icons/circle-check-solid.svg";
import BugIcon from "../../assets/icons/bug-solid.svg";
import PlaneIcon from "../../assets/icons/paper-plane-solid.svg";
import ErrorIcon from "../../assets/icons/circle-xmark-solid.svg";
import TimerIcon from "../../assets/icons/hourglass-half-solid.svg";
import ClockIcon from "../../assets/icons/clock-solid.svg";
import restfulApiByTag, {
    ApiResponse,
} from "../../helpers/request/restfulApiByTag";
import Polling from "../../components/Polling";
import setBanner from "./setBanner";
import setLabels from "./setLabels";
import setMap from "./setMap";
import setAreas from "./setAreas";
import { connect } from "react-redux";
import { ReduxStoreProps } from "../../config/store";
import { update as updateADC } from "../../store/adc";
import { update as updateGeophone } from "../../store/geophone";
import mapStateToProps from "../../helpers/utils/mapStateToProps";
import { WithTranslation, withTranslation } from "react-i18next";

// 120s by default
const QUENE_LENGTH = 120;

export interface HomeArea {
    readonly tag: string;
    readonly area: AreaProps;
    readonly chart: ChartProps;
}

export interface HomeMap {
    readonly area: AreaProps;
    readonly instance: MapBoxProps;
}

export interface HomeState {
    readonly banner: BannerProps;
    readonly labels: LabelProps[];
    readonly areas: HomeArea[];
    readonly map: HomeMap;
}

class Home extends Component<ReduxStoreProps & WithTranslation, HomeState> {
    constructor(props: ReduxStoreProps & WithTranslation) {
        super(props);
        this.state = {
            banner: {
                type: "warning",
                label: { id: "views.home.banner.warning.label" },
                text: { id: "views.home.banner.warning.text" },
            },
            labels: [
                {
                    tag: "messages",
                    label: { id: "views.home.labels.messages.label" },
                    unit: { id: "views.home.labels.messages.unit" },
                    value: "0",
                    icon: CheckIcon,
                    color: true,
                },
                {
                    tag: "errors",
                    label: { id: "views.home.labels.errors.label" },
                    unit: { id: "views.home.labels.errors.unit" },
                    value: "0",
                    icon: BugIcon,
                    color: true,
                },
                {
                    tag: "pushed",
                    label: { id: "views.home.labels.pushed.label" },
                    unit: { id: "views.home.labels.pushed.unit" },
                    value: "0",
                    icon: PlaneIcon,
                    color: true,
                },
                {
                    tag: "failures",
                    label: { id: "views.home.labels.failures.label" },
                    unit: { id: "views.home.labels.failures.unit" },
                    value: "0",
                    icon: ErrorIcon,
                    color: true,
                },
                {
                    tag: "queued",
                    label: { id: "views.home.labels.queued.label" },
                    unit: { id: "views.home.labels.queued.unit" },
                    value: "0",
                    icon: TimerIcon,
                    color: true,
                },
                {
                    tag: "offset",
                    label: { id: "views.home.labels.offset.label" },
                    unit: { id: "views.home.labels.offset.unit" },
                    value: "0",
                    icon: ClockIcon,
                    color: true,
                },
            ],
            areas: [
                {
                    tag: "cpu",
                    area: {
                        label: { id: "views.home.areas.cpu.label" },
                        text: {
                            id: "views.home.areas.cpu.text",
                            format: { usage: "0.00%" },
                        },
                    },
                    chart: {
                        backgroundColor: "#22c55e",
                        lineWidth: 5,
                        height: 250,
                        series: {
                            type: "line",
                            color: "#fff",
                            data: [],
                        },
                    },
                },
                {
                    tag: "memory",
                    area: {
                        label: { id: "views.home.areas.memory.label" },
                        text: {
                            id: "views.home.areas.memory.text",
                            format: { usage: "0.00%" },
                        },
                    },
                    chart: {
                        backgroundColor: "#06b6d4",
                        lineWidth: 5,
                        height: 250,
                        series: {
                            type: "line",
                            color: "#fff",
                            data: [],
                        },
                    },
                },
            ],
            map: {
                area: {
                    label: { id: "views.home.map.area.label" },
                    text: {
                        id: "views.home.map.area.text",
                        format: {
                            altitude: "0.00",
                            latitude: "0.00",
                            longitude: "0.00",
                        },
                    },
                },
                instance: {
                    zoom: 6,
                    minZoom: 3,
                    maxZoom: 7,
                    flyTo: true,
                    center: [0, 0],
                    tile: "/tiles/{z}/{x}/{y}/tile.webp",
                },
            },
        };
    }

    // Handler for Polling component to fetch status from server
    handleFetch = async (tag: string): Promise<ApiResponse> => {
        const res = await restfulApiByTag({ tag });
        return res;
    };

    // Handler for Polling component to handle error
    handleError = (): void => {
        // setBanner returns error when no arguments passed
        const banner = setBanner();
        this.setState({ banner });
    };

    // Handler for Polling component to process server response
    handleData = (res: ApiResponse): void => {
        const { error } = res;
        const banner = setBanner(res);

        if (!error) {
            // Update ADC & Geophone
            const { adc, geophone } = res.data;
            const { updateADC, updateGeophone } = this.props;
            // Update status labels
            const map = setMap(this.state.map, res);
            const labels = setLabels(this.state.labels, res);
            const areas = setAreas(this.state.areas, res, QUENE_LENGTH);
            // Update component state
            this.setState({ labels, areas, map });
            // Apply ADC & Geophone parameters to Redux store
            updateGeophone && updateGeophone(geophone);
            updateADC && updateADC(adc);
        }

        // Update banner
        this.setState({ banner });
    };

    render() {
        const { banner, labels, areas } = this.state;
        const { area, instance } = this.state.map;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />
                    <Polling
                        timer={1000}
                        tag={"station"}
                        onData={this.handleData}
                        onFetch={this.handleFetch}
                        onError={this.handleError}
                    >
                        <Banner {...banner} />

                        <Container layout={"flex"}>
                            {labels.map((label, index) => (
                                <Label
                                    key={index}
                                    {...label}
                                    className="md:w-1/2 lg:w-1/3"
                                />
                            ))}
                        </Container>

                        <Container layout={"grid"}>
                            {areas.map(({ area, chart }, index) => (
                                <Area key={index} {...area}>
                                    <Chart {...chart} />
                                </Area>
                            ))}
                        </Container>

                        <Container layout={"none"}>
                            <Area label={area.label} text={area.text}>
                                <MapBox className="h-[400px]" {...instance} />
                            </Area>
                        </Container>
                    </Polling>
                </Content>

                <Scroller />
                <Footer />
            </View>
        );
    }
}

export default connect(mapStateToProps, {
    updateGeophone,
    updateADC,
})(withTranslation()(Home));
