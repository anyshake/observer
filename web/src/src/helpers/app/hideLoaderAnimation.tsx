export const hideLoaderAnimation = async () => {
    const loadingScreen = document.querySelector('.public-loading') as HTMLElement;
    if (loadingScreen) {
        loadingScreen.style.opacity = '0';
        await new Promise((resolve) => setTimeout(resolve, 200));
        document.querySelector('.public-loading')?.remove();
    }
};
