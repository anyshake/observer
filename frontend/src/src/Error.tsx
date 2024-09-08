import { mdiBugPause } from "@mdi/js";
import Icon from "@mdi/react";
import { useEffect, useState } from "react";

import { Code } from "./components/Code";
import { hideLoading } from "./helpers/app/hideLoading";

interface ErrorProps {
	readonly code?: number;
	readonly heading?: string;
	readonly content?: string;
	readonly action?: {
		readonly label: string;
		readonly onClick: () => void;
	};
	readonly debug?: string;
}

export const Error = ({ code, heading, content, action, debug }: ErrorProps) => {
	const [isDebug, setIsDebug] = useState(false);

	useEffect(() => {
		hideLoading();
	}, []);

	return (
		<div className="p-5 flex min-h-screen flex-col items-center justify-center space-y-4">
			<h1 className="text-6xl font-bold tracking-tight text-gray-800 dark:text-gray-200">
				{code ?? ":-("}
			</h1>
			<p className="text-lg font-medium text-gray-600 dark:text-gray-500">
				{heading ?? "Oops! Something went wrong."}
			</p>
			<p className="text-gray-500 dark:text-gray-600">
				{content ?? "Please try to refresh the page or come back later."}
			</p>

			{action && (
				<button
					onClick={action.onClick}
					className="px-4 py-2 text-gray-200 bg-gray-500 rounded-md hover:bg-gray-600 transition-all"
				>
					{action.label}
				</button>
			)}

			{debug && (
				<button
					className="p-2 rounded-full bg-gray-100 dark:bg-gray-700 hover:scale-110 transition-all"
					onClick={() => {
						setIsDebug(!isDebug);
					}}
				>
					<Icon
						className="text-gray-400 dark:text-gray-500"
						path={mdiBugPause}
						size={1}
					/>
				</button>
			)}

			{debug && isDebug && (
				<div className="container w-[350px] md:w-[600px]">
					<Code language="javascript" fileName={`error_stack_${Date.now()}.log`}>
						{debug}
					</Code>
				</div>
			)}
		</div>
	);
};
