import { mdiGithub } from "@mdi/js";
import Icon from "@mdi/react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { i18nConfig } from "../config/i18n";

interface FooterProps {
	readonly author: string;
	readonly homepage: string;
	readonly repository: string;
	readonly currentLocale: string;
	readonly text: Record<keyof typeof i18nConfig.resources, string>;
}

export const Footer = (props: FooterProps) => {
	const { text, homepage, currentLocale, author, repository } = props;
	const [currentYear] = useState(new Date().getFullYear());

	return (
		<footer className="w-full bg-gray-200 text-gray-500 flex flex-col px-6 py-2 sm:flex-row justify-between">
			<span className="text-xs text-center md:ml-12 self-center">{text[currentLocale]}</span>
			<div className="flex text-center justify-center">
				<Link
					className="text-sm hover:underline self-center"
					to={homepage}
					target="_blank"
				>{`© ${currentYear} ${author}`}</Link>
				<Link to={repository} target="_blank">
					<Icon className="mx-2 self-center" path={mdiGithub} size={1} />
				</Link>
			</div>
		</footer>
	);
};
