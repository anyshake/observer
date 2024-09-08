import { mdiEarth } from "@mdi/js";
import Icon from "@mdi/react";
import { ChangeEvent, useState } from "react";

interface HeaderProps {
	readonly title: string;
	readonly currentLocale: string;
	readonly locales: Record<string, string>;
	readonly onSwitchLocale: (locale: string) => void;
}

export const Header = (props: HeaderProps) => {
	const { title, currentLocale, locales, onSwitchLocale } = props;
	const [newLocale, setNewLocale] = useState<string | null>(null);

	const handleSelectChange = ({ target }: ChangeEvent<HTMLSelectElement>) => {
		setNewLocale(target.value);
		onSwitchLocale(target.value);
	};

	return (
		<header className="fixed w-full z-10 flex justify-between bg-gray-200 items-center h-16 px-5">
			<h1 className="ml-14 text-gray-800 text-xl font-bold">{title}</h1>
			<div className="flex text-gray-500 space-x-1">
				<Icon path={mdiEarth} size={0.8} />
				<select
					className="text-xs bg-transparent focus:outline-none max-w-[4.9rem] truncate"
					onChange={handleSelectChange}
					value={newLocale ?? currentLocale}
				>
					<option disabled>Choose Language</option>
					{Object.entries(locales).map(([key, value]) => (
						<option key={key} value={key} className="text-gray-800">
							{value}
						</option>
					))}
				</select>
			</div>
		</header>
	);
};
