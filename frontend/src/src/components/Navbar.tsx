import { Link } from "react-router-dom";

import ArrowIcon from "../assets/icons/angle-right-solid.svg";
import HomeIcon from "../assets/icons/house-solid.svg";

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
						<img className="my-2 w-5 h-4 mr-2" src={HomeIcon} alt="" />
						<span className="my-2">/</span>
					</Link>
				</li>
				{pathname !== basename && (
					<li className="flex">
						<img className="self-center w-4 h-4 mr-2" src={ArrowIcon} alt="" />
						<Link className="my-2 cursor-pointer hover:text-gray-900" to={pathname}>
							{title}
						</Link>
					</li>
				)}
			</ol>
		</nav>
	);
};
