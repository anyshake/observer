const getCopyright = (name: string, author: string) => {
    return `© ${new Date().getFullYear()} ${name} | ${author}`;
};

export default getCopyright;
