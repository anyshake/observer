import { mdiDelete } from "@mdi/js";
import Icon from "@mdi/react";
import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import Typography from "@mui/material/Typography";

export interface ProgressProps {
	readonly precision?: number;
	readonly label?: string;
	readonly value: number;
	readonly onCancel?: () => void;
}

export const Progress = (props: ProgressProps) => {
	const { value, label, precision, onCancel } = props;

	return (
		<Box sx={{ display: "flex", alignItems: "center" }}>
			<div
				className="cursor-pointer px-1 rounded-lg bg-gray-100 hover:bg-gray-300"
				onClick={onCancel}
			>
				<Icon path={mdiDelete} size={0.8} />
			</div>
			<Box sx={{ width: "100%", mx: 1 }}>
				<LinearProgress
					className="rounded-full"
					variant="determinate"
					color="secondary"
					value={value}
				/>
			</Box>
			<Box sx={{ minWidth: 100 }}>
				<Typography className="overflow-scroll p-2" color="text.secondary" variant="body2">
					{`[${value.toFixed(precision ?? 2)}%] ${label}`}
				</Typography>
			</Box>
		</Box>
	);
};
