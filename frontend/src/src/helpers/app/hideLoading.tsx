export const hideLoading = () => {
	const loadingScreen = document.querySelector(".public-loading") as HTMLElement;
	if (loadingScreen) {
		loadingScreen.style.opacity = "0";
		setTimeout(() => {
			document.querySelector(".public-loading")?.remove();
		}, 500);
	}
};
