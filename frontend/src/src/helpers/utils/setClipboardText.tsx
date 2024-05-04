export const setClipboardText = async (text: string) => {
	const clipboard = navigator.clipboard ?? {
		writeText: (text) => {
			const copyInput = document.createElement("input");
			copyInput.value = text;
			document.body.appendChild(copyInput);
			copyInput.select();
			document.execCommand("copy");
			document.body.removeChild(copyInput);
		}
	};

	if (clipboard) {
		await clipboard.writeText(text);
		return true;
	}

	return false;
};
