import { mdiAlertCircle } from "@mdi/js";
import Icon from "@mdi/react";

const NotFound = () => {
	return (
		<div className="w-full min-h-screen flex flex-col items-center justify-center">
			<div className="animate-bounce text-amber-600">
				<Icon className="py-2" path={mdiAlertCircle} size={4} />
			</div>

			<h2 className="py-2 text-2xl font-bold text-gray-600">404 Not Found</h2>
		</div>
	);
};

export default NotFound;
