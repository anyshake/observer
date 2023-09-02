import { Component, RefObject, createRef } from "react";
import { MapContainer, Marker, TileLayer } from "react-leaflet";
import LocationIcon from "../assets/icons/location-dot-solid.svg";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

export interface MapBoxProps {
    readonly className?: string;
    readonly minZoom: number;
    readonly maxZoom: number;
    readonly zoom: number;
    readonly tile: string;
    readonly center: [number, number];
    readonly marker?: [number, number];
    readonly zoomControl?: boolean;
    readonly flyTo?: boolean;
    readonly dragging?: boolean;
}

export interface MapBoxState {
    readonly map: RefObject<L.Map> | undefined;
}

export default class MapBox extends Component<MapBoxProps, MapBoxState> {
    constructor(props: MapBoxProps) {
        super(props);
        this.state = {
            map: createRef(),
        };
    }

    componentDidUpdate(): void {
        const { center, flyTo, zoom } = this.props;
        const map = this.state.map?.current;
        if (map && flyTo) {
            map?.flyTo(center, zoom);
        }
    }

    render() {
        const {
            className,
            minZoom,
            maxZoom,
            center,
            zoom,
            tile,
            marker,
            dragging,
            zoomControl,
        } = this.props;
        const { map } = this.state;
        const icon = new L.Icon({
            iconUrl: LocationIcon,
            iconAnchor: [13, 28],
            iconSize: [18, 25],
        });

        return (
            <MapContainer
                ref={map}
                className={`z-0 w-full ${className}`}
                zoomControl={zoomControl}
                attributionControl={false}
                scrollWheelZoom={false}
                doubleClickZoom={false}
                dragging={dragging}
                maxZoom={maxZoom}
                minZoom={minZoom}
                center={center}
                zoom={zoom}
                style={{
                    cursor: "default",
                }}
            >
                <TileLayer url={tile} />
                {marker && <Marker position={marker} icon={icon} />}
            </MapContainer>
        );
    }
}
