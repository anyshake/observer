import { mdiCheck, mdiClipboardOutline, mdiContentSave } from "@mdi/js";
import Icon from "@mdi/react";
import { saveAs } from "file-saver";
import { useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import theme from "react-syntax-highlighter/dist/esm/styles/prism/atom-dark";

import { setClipboardText } from "../helpers/utils/setClipboardText";

interface CodeProps {
	readonly language?: string;
	readonly fileName?: string;
	readonly children: string;
}

export const Code = (props: CodeProps) => {
	const { fileName, language, children } = props;

	const [isCopied, setIsCopied] = useState(false);

	const handleCopy = async (text: string) => {
		await setClipboardText(text);
		setIsCopied(true);
		setTimeout(() => {
			setIsCopied(false);
		}, 2000);
	};

	const handleDownload = (text: string) => {
		const blob = new Blob([text], {
			type: "text/plain;charset=utf-8"
		});
		saveAs(blob, fileName);
	};

	return (
		<div className="rounded-lg bg-gray-700 p-2">
			<div className="flex justify-end space-x-3 px-4 text-white">
				<div
					className="opacity-60 hover:opacity-100 transition-all cursor-pointer"
					onClick={() => {
						void handleCopy(children);
					}}
				>
					<Icon path={isCopied ? mdiCheck : mdiClipboardOutline} size={0.8} />
				</div>
				{!!fileName?.length && (
					<div
						className="opacity-60 hover:opacity-100 transition-all cursor-pointer"
						onClick={() => {
							handleDownload(children);
						}}
					>
						<Icon path={mdiContentSave} size={0.8} />
					</div>
				)}
			</div>
			<SyntaxHighlighter language={language} style={theme}>
				{children}
			</SyntaxHighlighter>
		</div>
	);
};
