export interface Geophone {
    readonly ehz: number;
    readonly ehe: number;
    readonly ehn: number;
    readonly [key: string]: number;
}
