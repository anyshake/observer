import { mdiEarth, mdiExitToApp } from "@mdi/js";
import Icon from "@mdi/react";
import { ChangeEvent, useState } from "react";

interface HeaderProps {
	readonly title: string;
	readonly currentLocale: string;
	readonly locales: Record<string, string>;
	readonly onSwitchLocale: (locale: string) => void;
	readonly onLogout?: () => void;
}

export const Header = ({
	title,
	currentLocale,
	locales,
	onSwitchLocale,
	onLogout
}: HeaderProps) => {
	const [newLocale, setNewLocale] = useState<string | null>(null);

	const handleSelectChange = ({ target }: ChangeEvent<HTMLSelectElement>) => {
		setNewLocale(target.value);
		onSwitchLocale(target.value);
	};

	return (
		<header className="fixed w-full z-10 flex justify-between bg-gray-200 items-center h-16 px-6">
			<h1 className="ml-14 text-gray-800 text-xl font-bold">{title}</h1>

			<div className="flex space-x-3 -mr-2 sm:mr-2">
				<div className="flex text-gray-500 space-x-1 hover:opacity-80">
					<Icon className="self-center" path={mdiEarth} size={0.8} />
					<select
						className="text-xs text-center bg-transparent focus:outline-none max-w-[4.8rem] truncate cursor-pointer"
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

				{onLogout && (
					<>
						<hr className="h-5 w-[1px] bg-gray-400" />
						<button
							className="flex text-gray-500 space-x-1 hover:opacity-80"
							onClick={onLogout}
						>
							<Icon className="self-center" path={mdiExitToApp} size={0.8} />
						</button>
					</>
				)}
			</div>
		</header>
	);
};
