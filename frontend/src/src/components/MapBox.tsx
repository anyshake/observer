import "leaflet/dist/leaflet.css";

import { mdiMapMarker } from "@mdi/js";
import { divIcon, Map } from "leaflet";
import { useEffect, useRef } from "react";
import { MapContainer, Marker, TileLayer } from "react-leaflet";

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
		dragging
	} = props;
	const icon = divIcon({
		className: "leaflet-data-marker",
		html: `<svg viewBox="0 0 24 24" style="width: 32px; height: 32px; fill: #6565f1; filter: drop-shadow(2px 4px 6px rgba(0, 0, 0, 0.5));">
            <path d="${mdiMapMarker}" stroke="black" stroke-width="0.7" />
        </svg>`,
		iconAnchor: [16, 32]
	});

	const mapRef = useRef<Map>(null);

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
			style={{ cursor: "default" }}
		>
			<TileLayer url={tile} />
			{marker && <Marker position={marker} icon={icon} />}
		</MapContainer>
	);
};
