export const setClipboardText = async (text: string): Promise<boolean> => {
    if (navigator.clipboard?.writeText) {
        try {
            await navigator.clipboard.writeText(text);
            return true;
        } catch {
            // Some browsers expose the API but reject writes; try the selection fallback.
        }
    }

    const previousFocus = document.activeElement;
    const copyInput = document.createElement('textarea');
    copyInput.value = text;
    copyInput.readOnly = true;
    copyInput.style.position = 'fixed';
    copyInput.style.opacity = '0';
    document.body.appendChild(copyInput);

    try {
        copyInput.focus();
        copyInput.select();
        return document.execCommand('copy');
    } catch {
        return false;
    } finally {
        copyInput.remove();
        if (previousFocus instanceof HTMLElement) {
            previousFocus.focus();
        }
    }
};
