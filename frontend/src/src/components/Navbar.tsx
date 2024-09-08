import { mdiChevronRight, mdiHome } from "@mdi/js";
import Icon from "@mdi/react";
import { Link } from "react-router-dom";

interface NavbarProps {
	title: string;
	basename: string;
	pathname: string;
}

export const Navbar = (props: NavbarProps) => {
	const { title, basename, pathname } = props;

	return (
		<nav className="px-5 py-3 rounded-lg bg-gray-100">
			<ol className="text-sm font-medium text-gray-700 flex space-x-2">
				<li className="cursor-pointer hover:text-gray-900">
					<Link className="flex" to={"/"}>
						<Icon className="self-center mr-2" path={mdiHome} size={1} />
						<span className="my-2">/</span>
					</Link>
				</li>
				{pathname !== basename && (
					<li className="flex">
						<Icon className="self-center mr-2" path={mdiChevronRight} size={1} />
						<Link className="my-2 cursor-pointer hover:text-gray-900" to={pathname}>
							{title}
						</Link>
					</li>
				)}
			</ol>
		</nav>
	);
};
