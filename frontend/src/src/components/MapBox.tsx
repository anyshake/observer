import { useEffect, useRef } from "react";
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
    readonly scrollWheelZoom?: boolean;
    readonly zoomControl?: boolean;
    readonly flyTo?: boolean;
    readonly dragging?: boolean;
}

export const MapBox = (props: MapBoxProps) => {
    const {
        className,
        minZoom,
        flyTo,
        maxZoom,
        zoom,
        tile,
        center,
        marker,
        scrollWheelZoom,
        zoomControl,
        dragging,
    } = props;
    const icon = new L.Icon({
        iconUrl: LocationIcon,
        iconAnchor: [9, 24],
        iconSize: [18, 25],
    });

    const mapRef = useRef<L.Map>(null);

    useEffect(() => {
        const map = mapRef.current;
        if (map) {
            map.flyTo(center, zoom);
        }
    }, [center, zoom, flyTo]);

    return (
        <MapContainer
            ref={mapRef}
            className={`z-0 w-full ${className ?? ""}`}
            scrollWheelZoom={scrollWheelZoom}
            zoomControl={zoomControl}
            attributionControl={false}
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
};
