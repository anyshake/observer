import type { LayerSpecification } from 'maplibre-gl';

import { getMapTilesUrl } from '../helpers/app/getMapTilesUrl';

export const HomeConstraints = {
    lineChartRetention: 100,
    pollInterval: 2000,
    maxGapSeconds: 10,
    mapMinZoom: 1,
    mapMaxZoom: 6,
    mapDefaultZoom: 2,
    mapTileUrl: getMapTilesUrl(),
    mapTileLayers: [
        {
            id: 'background',
            type: 'background',
            paint: { 'background-color': '#f8f9fc' }
        },
        {
            id: 'ocean',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_ocean',
            paint: { 'fill-color': '#d1d1f0' }
        },
        {
            id: 'countries',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_admin_0_countries',
            paint: {
                'fill-color': '#f2f0eb'
            }
        },
        {
            id: 'glaciated',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_glaciated_areas',
            paint: { 'fill-color': '#ffffff', 'fill-opacity': 0.8 }
        },
        {
            id: 'ice-shelves',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_antarctic_ice_shelves_polys',
            paint: { 'fill-color': '#eef2ff', 'fill-opacity': 0.8 }
        },
        {
            id: 'lakes',
            type: 'fill',
            source: 'naturalearth',
            'source-layer': 'ne_10m_lakes',
            paint: { 'fill-color': '#d1d1f0' }
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
            paint: {
                'line-color': '#b9b0db',
                'line-width': 0.1
            }
        },
        {
            id: 'coastline',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_coastline',
            paint: {
                'line-color': '#9b90c2',
                'line-width': 0.8
            }
        },
        {
            id: 'admin-0-outline',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_admin_0_countries',
            paint: {
                'line-color': '#9285c5',
                'line-width': 0.6,
                'line-dasharray': [4, 2]
            }
        },
        {
            id: 'admin-1-lines',
            type: 'line',
            source: 'naturalearth',
            'source-layer': 'ne_10m_admin_1_states_provinces_lines',
            paint: {
                'line-color': '#b9b0db',
                'line-width': 0.8,
                'line-dasharray': [1, 1]
            }
        },
        {
            id: 'populated-places',
            type: 'circle',
            source: 'naturalearth',
            'source-layer': 'ne_10m_populated_places',
            paint: {
                'circle-radius': 2.5,
                'circle-color': '#6a5acd',
                'circle-stroke-width': 1,
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
                'text-size': 10,
                'text-offset': [0, -0.6],
                'text-anchor': 'bottom',
                'text-font': ['monospace'],
                'text-letter-spacing': 0.05
            },
            paint: {
                'text-color': '#4a4458',
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
