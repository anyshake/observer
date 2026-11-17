const text = (value: unknown) => (typeof value === 'string' && value.length > 0 ? value : '');

export interface IStationIdentity {
    readonly name: string;
    readonly code: string;
    readonly place: string;
}

// Derives display name, NET.STA.LOC code and place from the getStationConfig map.
export const getStationIdentity = (config: Record<string, unknown> | null | undefined) => {
    const values = config ?? {};
    return {
        name: text(values.station_name),
        code: [values.network_code, values.station_code, values.location_code]
            .map(text)
            .filter(Boolean)
            .join('.'),
        place: [text(values.station_place), text(values.station_country)].filter(Boolean).join(', ')
    } satisfies IStationIdentity;
};
