import { Map, Marker } from 'pigeon-maps';

interface IMapContainer {
    readonly height: number;
    readonly zoom?: number;
    readonly borderRadius?: string;
    readonly coordinates: number[];
}

export const MapContainer = ({
    height,
    coordinates,
    zoom,
    borderRadius = '8px'
}: IMapContainer) => {
    const [latitude, longitude] = coordinates;

    return (
        <div
            style={{
                borderRadius: borderRadius,
                overflow: 'hidden',
                width: '100%',
                height: `${height}px`,
                position: 'relative'
            }}
        >
            <Map
                animate={true}
                height={height}
                center={[latitude, longitude]}
                defaultZoom={zoom ?? 6}
                provider={(x, y, z) => `/tiles/${z}/${x}/${y}.webp`}
                maxZoom={7}
                minZoom={3}
                attribution={false}
            >
                <Marker width={30} anchor={[latitude, longitude]} />
            </Map>
        </div>
    );
};
