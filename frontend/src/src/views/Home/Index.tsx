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
import restfulApiByTag, { ApiResponse } from "../../helpers/request/restfulApiByTag";
import Polling from "../../components/Polling";
import setBanner from "./setBanner";
import setLabels from "./setLabels";
import setMap from "./setMap";
import setAreas from "./setAreas";
import { connect } from "react-redux";
import { ReduxStore, ReduxStoreProps } from "../../config/store";
import { update as updateADC } from "../../store/adc";
import { update as updateGeophone } from "../../store/geophone";

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

class Home extends Component<ReduxStoreProps, HomeState> {
    constructor(props: ReduxStoreProps) {
        super(props);
        this.state = {
            banner: {
                type: "warning",
                label: "正在连接服务器",
                text: "请稍等...",
            },
            labels: [
                {
                    tag: "messages",
                    label: "已解码讯息量",
                    value: "0",
                    unit: "条",
                    icon: CheckIcon,
                    color: true,
                },
                {
                    tag: "errors",
                    label: "帧错误讯息量",
                    value: "0",
                    unit: "条",
                    icon: BugIcon,
                    color: true,
                },
                {
                    tag: "pushed",
                    label: "已推送讯息量",
                    value: "0",
                    unit: "条",
                    icon: PlaneIcon,
                    color: true,
                },
                {
                    tag: "failures",
                    label: "推送失败讯息量",
                    value: "0",
                    unit: "条",
                    icon: ErrorIcon,
                    color: true,
                },
                {
                    tag: "queued",
                    label: "等待推送讯息量",
                    value: "0",
                    unit: "条",
                    icon: TimerIcon,
                    color: true,
                },
                {
                    tag: "offset",
                    label: "系统时间偏移量",
                    value: "0",
                    unit: "秒",
                    icon: ClockIcon,
                    color: true,
                },
            ],
            areas: [
                {
                    tag: "cpu",
                    area: {
                        label: "CPU 占用率",
                        text: "当前占用率：正在获取中",
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
                        label: "RAM 占用率",
                        text: "当前占用率：正在获取中",
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
                    label: "测站所在位置",
                    text: "正在加载位置资讯",
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

    handleFetch = async (tag: string): Promise<ApiResponse> => {
        const res = await restfulApiByTag({ tag });
        return res;
    };

    handleError = (): void => {
        const banner = setBanner();
        this.setState({ banner });
    };

    handleData = (res: ApiResponse): void => {
        const { error } = res;
        const banner = setBanner(res);

        if (!error) {
            // Update ADC & Geophone
            const { adc, geophone } = res.data;
            const { updateADC, updateGeophone } = this.props;
            // Update redux store
            updateGeophone(geophone);
            updateADC(adc);
            // Update labels state
            const map = setMap(this.state.map, res);
            const labels = setLabels(this.state.labels, res);
            const areas = setAreas(this.state.areas, res, QUENE_LENGTH);
            // Update this.state
            this.setState({ labels, areas, map });
        }

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
                        retry={3}
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

const mapStateToProps = (state: ReduxStore) => {
    return { ...state };
};

export default connect(mapStateToProps, {
    updateGeophone,
    updateADC,
})(Home);
