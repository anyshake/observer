import type { LayerSpecification } from 'maplibre-gl';

import { getMapTilesUrl } from '../helpers/app/getMapTilesUrl';

export const HomeConstraints = {
    lineChartRetention: 100,
    pollInterval: 2000,
    maxGapSeconds: 10,
    mapMinZoom: 3,
    mapMaxZoom: 7,
    mapDefaultZoom: 4,
    mapTileUrl: getMapTilesUrl(),
    mapTileLayers: [
        {
            id: 'background',
            type: 'background',
            paint: { 'background-color': '#dceaf6' }
        },
        {
            id: 'ocean',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_ocean',
            paint: { 'fill-color': '#a3c4e7' }
        },
        {
            id: 'countries',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_admin_0_countries',
            paint: {
                'fill-color': '#e8e4d8',
                'fill-outline-color': '#b0a890'
            }
        },
        {
            id: 'glaciated',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_glaciated_areas',
            paint: { 'fill-color': '#ffffff', 'fill-opacity': 0.6 }
        },
        {
            id: 'ice-shelves',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_antarctic_ice_shelves_polys',
            paint: { 'fill-color': '#f0f8ff', 'fill-opacity': 0.7 }
        },
        {
            id: 'lakes',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_lakes',
            paint: { 'fill-color': '#a3c4e7' }
        },
        {
            id: 'playas',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_playas',
            paint: { 'fill-color': '#d4c9a8', 'fill-opacity': 0.5 }
        },
        {
            id: 'reefs',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_reefs',
            paint: { 'line-color': '#72b8d4', 'line-width': 1 }
        },
        {
            id: 'minor-islands',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_minor_islands',
            paint: { 'fill-color': '#e8e4d8' }
        },
        {
            id: 'rivers',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_rivers_lake_centerlines',
            paint: { 'line-color': '#a3c4e7', 'line-width': 1 }
        },
        {
            id: 'coastline',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_coastline',
            paint: { 'line-color': '#7eb8d0', 'line-width': 0.8 }
        },
        {
            id: 'admin-1-lines',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_admin_1_states_provinces_lines',
            paint: {
                'line-color': '#c0b898',
                'line-width': 0.5,
                'line-dasharray': [3, 2]
            }
        },
        {
            id: 'populated-places',
            type: 'circle',
            source: 'naturalearth',
            'source-layer': 'ne_10m_populated_places',
            paint: {
                'circle-radius': 2,
                'circle-color': '#666666',
                'circle-stroke-width': 0.5,
                'circle-stroke-color': '#ffffff'
            }
        },
        {
            id: 'populated-places-label',
            type: 'symbol',
            source: 'naturalearth',
            'source-layer': 'ne_10m_populated_places',
            layout: {
                'text-field': ['coalesce', ['get', 'NAME'], ['get', 'name']],
                'text-size': 11,
                'text-offset': [0, -1.2],
                'text-anchor': 'bottom',
                'text-max-width': 6,
                'text-allow-overlap': false
            },
            paint: {
                'text-color': '#333333',
                'text-halo-color': '#ffffff',
                'text-halo-width': 1.5
            }
        }
    ] satisfies LayerSpecification[]
};

export const RealTimeConstraints = {
    id: 'realtime',
    minWidth: 200,
    minHeight: 150,
    maxWidth: 800,
    maxHeight: 600,
    // waveform default options
    minSpanValue: 100,
    lineColor: '#8A3EED',
    // spectrogram default options
    fftSize: 1024,
    windowSize: 512,
    overlap: Math.floor(512 * 0.86),
    freqRange: [0, 25] as [number, number],
    getDynamicDB: (index: number) => {
        if (index > 2 && index < 6) {
            return { minDB: 20, maxDB: 120 };
        }
        return { minDB: 110, maxDB: 170 };
    }
};

export const HistoryConstraints = {
    id: 'history',
    minWidth: 200,
    minHeight: 150,
    maxWidth: 800,
    maxHeight: 600,
    // waveform default options
    minSpanValue: 100,
    lineColor: '#8A3EED',
    // spectrogram default options
    fftSize: 1024,
    windowSize: 512,
    overlap: Math.floor(512 * 0.86),
    freqRange: [0, 25] as [number, number],
    getDynamicDB: (index: number) => {
        if (index > 2 && index < 6) {
            return { minDB: 20, maxDB: 120 };
        }
        return { minDB: 110, maxDB: 170 };
    }
};

export const DownloadConstraints = {
    pollInterval: 5000
};

export const SettingsConstraints = {
    pollInterval: 5000
};
