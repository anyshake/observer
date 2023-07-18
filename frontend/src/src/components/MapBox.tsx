import { Component } from "react";
import { MapContainer, Marker, TileLayer } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L from "leaflet";
import LocationIcon from "../assets/icons/location-dot-solid.svg";

export interface MapBoxProps {
    readonly className?: string;
    readonly minZoom: number;
    readonly maxZoom: number;
    readonly zoom: number;
    readonly tile: string;
    readonly center: [number, number];
    readonly marker?: [number, number];
    readonly flyTo?: boolean;
}

export default class MapBox extends Component<MapBoxProps> {
    render() {
        const { className, minZoom, maxZoom, center, zoom, tile, marker } =
            this.props;
        const icon = new L.Icon({
            iconUrl: LocationIcon,
            iconAnchor: [13, 28],
            iconSize: [18, 25],
        });

        return (
            <MapContainer
                className={`z-0 w-full ${className}`}
                attributionControl={false}
                scrollWheelZoom={false}
                zoomControl={true}
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
