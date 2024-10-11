import { i18nConfig } from "./i18n";

interface GlobalConfig {
	readonly name: string;
	readonly title: string;
	readonly author: string;
	readonly homepage: string;
	readonly repository: string;
	readonly duration: {
		readonly default: number;
		readonly maximum: number;
		readonly minimum: number;
	};
	readonly retention: {
		readonly default: number;
		readonly maximum: number;
		readonly minimum: number;
	};
	readonly footer: Record<keyof typeof i18nConfig.resources, string>;
}

export const globalConfig: GlobalConfig = {
	name: "Observer",
	author: "AnyShake",
	title: "AnyShake Observer",
	homepage: "https://anyshake.org",
	repository: "https://github.com/anyshake/observer",
	duration: { default: 300, maximum: 3600, minimum: 10 },
	retention: { default: 120, maximum: 600, minimum: 10 },
	footer: {
		"en-US": "Listen to the whispering earth.",
		"zh-CN": "听见地球",
		"zh-TW": "聽見地球",
	}
};
