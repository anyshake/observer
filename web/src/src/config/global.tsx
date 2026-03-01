import logo from '/anyshake.svg';

interface IGlobalConfig {
    readonly repository: string;
    readonly logo: string;
    readonly name: string;
    readonly copyright: string;
    readonly homepage: string;
    readonly update: number;
    readonly footer: string;
}

export const globalConfig: IGlobalConfig = {
    logo,
    name: 'config.global.name',
    footer: 'config.global.footer',
    update: 10 * 60 * 1000,
    copyright: 'SensePlex Limited',
    homepage: 'https://anyshake.org',
    repository: 'https://github.com/anyshake/observer'
};
