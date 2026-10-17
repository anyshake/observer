import type { Map } from 'maplibre-gl';

import type { ThemeColors } from './useThemeColors';

export const applyMapTheme = (map: Map, colors: ThemeColors['map']) => {
    const paint = {
        background: { 'background-color': colors.background },
        ocean: { 'fill-color': colors.water },
        countries: { 'fill-color': colors.land },
        glaciated: { 'fill-color': colors.ice },
        'ice-shelves': { 'fill-color': colors.iceShelf },
        lakes: { 'fill-color': colors.water },
        playas: { 'fill-color': colors.playas },
        reefs: { 'line-color': colors.reef },
        'minor-islands': { 'fill-color': colors.island },
        rivers: { 'line-color': colors.river },
        coastline: { 'line-color': colors.coast },
        'admin-0-outline': { 'line-color': colors.boundary },
        'admin-1-lines': { 'line-color': colors.boundary },
        'populated-places': {
            'circle-color': colors.city,
            'circle-stroke-color': colors.halo
        },
        'populated-places-label': {
            'text-color': colors.label,
            'text-halo-color': colors.halo
        }
    };

    for (const [id, properties] of Object.entries(paint)) {
        if (map.getLayer(id)) {
            for (const [property, value] of Object.entries(properties)) {
                map.setPaintProperty(id, property, value);
            }
        }
    }
};
