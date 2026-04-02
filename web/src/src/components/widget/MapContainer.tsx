import 'maplibre-gl/dist/maplibre-gl.css';

import { mdiMapMarker } from '@mdi/js';
import maplibregl from 'maplibre-gl';
import { useCallback, useEffect, useMemo, useRef } from 'react';

export interface IMapContainer {
    readonly height: number;
    readonly minZoom: number;
    readonly maxZoom: number;
    readonly zoom: number;
    readonly tileUrl: string;
    readonly layers: maplibregl.LayerSpecification[];
    readonly coordinates: number[];
    readonly scrollWheelZoom?: boolean;
    readonly zoomControl?: boolean;
    readonly borderRadius?: string;
    readonly dragging?: boolean;
    readonly onClick?: (coordinates: [number, number]) => void;
}

const createMapStyle = (
    minZoom: number,
    maxZoom: number,
    tileUrl: string,
    layers: maplibregl.LayerSpecification[]
): maplibregl.StyleSpecification => ({
    version: 8,
    layers,
    sources: {
        naturalearth: {
            type: 'vector',
            tiles: [tileUrl],
            minzoom: minZoom,
            maxzoom: maxZoom
        }
    }
});

const createMarkerElement = (): HTMLDivElement => {
    const el = document.createElement('div');
    el.innerHTML = `<svg viewBox="0 0 24 24" style="width: 32px; height: 32px; fill: #364153; filter: drop-shadow(2px 4px 6px rgba(0, 0, 0, 0.5));">
        <path d="${mdiMapMarker}" stroke="white" stroke-width="0.7" />
    </svg>`;
    el.style.cursor = 'pointer';
    return el;
};

const animateMapTo = (
    map: maplibregl.Map,
    center: [number, number],
    zoom: number,
    minZoom: number,
    animationId: number,
    getAnimationId: () => number,
    onComplete?: () => void
) => {
    const currentCenter = map.getCenter();
    const currentZoom = map.getZoom();
    const targetTravelZoom = Math.max(minZoom, 1);
    const travelZoom = Math.min(currentZoom, zoom, targetTravelZoom);
    const hasCenterChanged = currentCenter.lng !== center[0] || currentCenter.lat !== center[1];
    const hasZoomChanged = currentZoom !== zoom;

    if (!hasCenterChanged && !hasZoomChanged) {
        onComplete?.();
        return;
    }

    const runStage = (action: () => void, next?: () => void) => {
        if (getAnimationId() !== animationId) {
            return;
        }

        map.once('moveend', () => {
            if (getAnimationId() !== animationId) {
                return;
            }

            if (next) {
                next();
                return;
            }

            onComplete?.();
        });
        action();
    };

    map.stop();

    if (currentZoom > travelZoom) {
        runStage(
            () => {
                map.easeTo({
                    duration: 500,
                    easing: (t) => t,
                    essential: true,
                    zoom: travelZoom
                });
            },
            () => {
                runStage(
                    () => {
                        map.easeTo({
                            center,
                            duration: 1000,
                            easing: (t) => t,
                            essential: true,
                            zoom: travelZoom
                        });
                    },
                    () => {
                        if (zoom !== travelZoom) {
                            runStage(() => {
                                map.easeTo({
                                    center,
                                    duration: 500,
                                    easing: (t) => t,
                                    essential: true,
                                    zoom
                                });
                            });
                            return;
                        }

                        onComplete?.();
                    }
                );
            }
        );
        return;
    }

    runStage(
        () => {
            map.easeTo({
                center,
                duration: 1000,
                easing: (t) => t,
                essential: true,
                zoom: travelZoom
            });
        },
        () => {
            if (zoom !== travelZoom) {
                runStage(() => {
                    map.easeTo({
                        center,
                        duration: 500,
                        easing: (t) => t,
                        essential: true,
                        zoom
                    });
                });
                return;
            }

            onComplete?.();
        }
    );
};

export const MapContainer = ({
    height,
    minZoom,
    maxZoom,
    zoom,
    tileUrl,
    layers,
    borderRadius = '8px',
    coordinates,
    scrollWheelZoom,
    zoomControl,
    dragging,
    onClick
}: IMapContainer) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const mapRef = useRef<maplibregl.Map | null>(null);
    const markerRef = useRef<maplibregl.Marker | null>(null);
    const isMapLoadedRef = useRef(false);
    const hasAnimatedToTargetRef = useRef(false);
    const animationIdRef = useRef(0);
    const [latitude, longitude] = useMemo(() => coordinates, [coordinates]);
    const latestViewRef = useRef({ latitude, longitude, zoom });

    useEffect(() => {
        latestViewRef.current = { latitude, longitude, zoom };
    }, [latitude, longitude, zoom]);

    const showMarkerAt = useCallback((target: [number, number]) => {
        const marker = markerRef.current;
        const map = mapRef.current;

        if (!marker || !map) {
            return;
        }

        marker.setLngLat(target).addTo(map);
    }, []);

    useEffect(() => {
        if (!containerRef.current) {
            return;
        }

        const { latitude: currentLatitude, longitude: currentLongitude } = latestViewRef.current;

        const map = new maplibregl.Map({
            container: containerRef.current,
            style: createMapStyle(minZoom, maxZoom, tileUrl, layers),
            center: [0, 0],
            zoom,
            minZoom,
            maxZoom,
            attributionControl: false,
            doubleClickZoom: false,
            dragPan: dragging ?? true,
            scrollZoom: scrollWheelZoom ?? true
        });

        if (zoomControl) {
            map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-right');
        }

        const marker = new maplibregl.Marker({
            element: createMarkerElement(),
            anchor: 'bottom'
        }).setLngLat([currentLongitude, currentLatitude]);

        mapRef.current = map;
        markerRef.current = marker;
        isMapLoadedRef.current = false;
        hasAnimatedToTargetRef.current = false;

        map.once('load', () => {
            const {
                latitude: currentLatitude,
                longitude: currentLongitude,
                zoom: currentZoom
            } = latestViewRef.current;
            isMapLoadedRef.current = true;
            animationIdRef.current += 1;
            animateMapTo(
                map,
                [currentLongitude, currentLatitude],
                currentZoom,
                minZoom,
                animationIdRef.current,
                () => animationIdRef.current,
                () => showMarkerAt([currentLongitude, currentLatitude])
            );
            hasAnimatedToTargetRef.current = true;
        });

        return () => {
            animationIdRef.current += 1;
            isMapLoadedRef.current = false;
            hasAnimatedToTargetRef.current = false;
            marker.remove();
            map.remove();
            mapRef.current = null;
            markerRef.current = null;
        };
    }, [
        dragging,
        layers,
        maxZoom,
        minZoom,
        scrollWheelZoom,
        tileUrl,
        showMarkerAt,
        zoom,
        zoomControl
    ]);

    useEffect(() => {
        const map = mapRef.current;
        if (!map || !isMapLoadedRef.current) {
            return;
        }

        if (hasAnimatedToTargetRef.current) {
            animationIdRef.current += 1;
            animateMapTo(
                map,
                [longitude, latitude],
                zoom,
                minZoom,
                animationIdRef.current,
                () => animationIdRef.current,
                () => showMarkerAt([longitude, latitude])
            );
        }
    }, [latitude, longitude, minZoom, showMarkerAt, zoom]);

    const handleClick = useCallback(
        (e: maplibregl.MapMouseEvent) => {
            onClick?.([e.lngLat.lat, e.lngLat.lng]);
        },
        [onClick]
    );
    useEffect(() => {
        const map = mapRef.current;
        if (!map) {
            return;
        }
        map.on('click', handleClick);
        return () => {
            map.off('click', handleClick);
        };
    }, [handleClick]);

    return (
        <div
            ref={containerRef}
            className="z-0"
            style={{ cursor: 'default', borderRadius, height, overflow: 'hidden' }}
        />
    );
};
