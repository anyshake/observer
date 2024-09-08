import { mdiLanConnect, mdiLanDisconnect, mdiLanPending } from "@mdi/js";
import Icon from "@mdi/react";

export interface BannerProps {
	readonly title: string;
	readonly content: string;
	readonly className?: string;
	readonly type: "success" | "warning" | "error";
}

export const Banner = (props: BannerProps) => {
	const { content, title, type, className } = props;

	let theme = "";
	switch (type) {
		case "success":
			theme = "from-green-400 to-blue-500";
			break;
		case "warning":
			theme = "from-orange-400 to-orange-600";
			break;
		case "error":
			theme = "from-red-400 to-red-600";
			break;
	}

	return (
		<div
			className={`my-2 shadow-xl p-6 text-sm text-white rounded-lg bg-gradient-to-r ${theme} ${
				className ?? ""
			}`}
		>
			<div className="flex flex-col space-y-2">
				<div className="flex gap-2 font-bold text-lg">
					<Icon
						className={type === "success" ? "size-6" : "hidden"}
						path={mdiLanConnect}
						size={1}
					/>
					<Icon
						className={type === "warning" ? "size-6" : "hidden"}
						path={mdiLanPending}
						size={1}
					/>
					<Icon
						className={type === "error" ? "size-6" : "hidden"}
						path={mdiLanDisconnect}
						size={1}
					/>
					<span>{title}</span>
				</div>

				<span className="pl-3 text-md font-medium">
					{content.split("\n").map((item) => (
						<p key={item}>
							{item}
							<br />
						</p>
					))}
				</span>
			</div>
		</div>
	);
};
