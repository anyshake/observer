import { Component } from "react";
import Header from "../../components/Header";
import Sidebar from "../../components/Sidebar";
import Navbar from "../../components/Navbar";
import Content from "../../components/Content";
import Banner, { BannerProps } from "../../components/Banner";
import Badges, { BadgesProps } from "../../components/Badges";
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
import requestByTag, { ApiResponse } from "../../helpers/requestByTag";
import Polling from "../../components/Polling";
import setBanner from "./setBanner";
import setBadges from "./setBadges";
import setMap from "./setMap";
import setAreas from "./setAreas";

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
    readonly badges: BadgesProps;
    readonly areas: HomeArea[];
    readonly map: HomeMap;
}

export default class Home extends Component<{}, HomeState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            banner: {
                type: "warning",
                label: "正在连接服务器",
                text: "请稍等...",
            },
            badges: {
                list: [
                    {
                        tag: "messages",
                        label: "已解码讯息量",
                        value: "0",
                        unit: "条",
                        icon: CheckIcon,
                    },
                    {
                        tag: "errors",
                        label: "帧错误讯息量",
                        value: "0",
                        unit: "条",
                        icon: BugIcon,
                    },
                    {
                        tag: "pushed",
                        label: "已推送讯息量",
                        value: "0",
                        unit: "条",
                        icon: PlaneIcon,
                    },
                    {
                        tag: "failures",
                        label: "推送失败讯息量",
                        value: "0",
                        unit: "条",
                        icon: ErrorIcon,
                    },
                    {
                        tag: "queued",
                        label: "等待推送讯息量",
                        value: "0",
                        unit: "条",
                        icon: TimerIcon,
                    },
                    {
                        tag: "offset",
                        label: "系统时间偏移量",
                        value: "0",
                        unit: "秒",
                        icon: ClockIcon,
                    },
                ],
            },
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
                    className: "h-[400px]",
                    zoom: 3,
                    minZoom: 3,
                    maxZoom: 7,
                    flyTo: true,
                    marker: [0, 0],
                    center: [33, 106],
                    tile: "/tiles/{z}/{x}/{y}/tile.webp",
                },
            },
        };
    }

    handleFetch = async (tag: string): Promise<ApiResponse> => {
        return await requestByTag({ tag });
    };

    handleData = (res: ApiResponse): boolean => {
        const { error } = res;
        if (!res.data) {
            return false;
        }

        const banner = setBanner(res);
        if (!error) {
            const map = setMap(this.state.map, res);
            const badges = setBadges(this.state.badges, res);
            const areas = setAreas(this.state.areas, res, QUENE_LENGTH);
            this.setState({ banner, badges, areas, map });
            return true;
        }

        this.setState({ banner });
        return false;
    };

    render() {
        const { banner, badges, areas } = this.state;
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
                    >
                        <Banner {...banner} />
                        <Badges {...badges} />

                        <Container layout={"grid"}>
                            {areas.map(({ area, chart }, index) => (
                                <Area key={index} {...area}>
                                    <Chart {...chart} />
                                </Area>
                            ))}
                        </Container>

                        <Container layout={"none"}>
                            <Area label={area.label} text={area.text}>
                                <MapBox {...instance} />
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
