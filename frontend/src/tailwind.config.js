/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
        "./node_modules/react-tailwindcss-datepicker/dist/index.esm.js",
    ],
    theme: {
        extend: {
            backgroundImage: {},
            keyframes: {
                appear: {
                    "0%": { opacity: 0 },
                    "100%": { opacity: 1 },
                },
                disappear: {
                    "0%": { opacity: 1 },
                    "100%": { opacity: 0 },
                },
            },
            animation: {
                appear: "appear 1s ease 0s 1 forwards",
                disappear: "disappear 1s ease 0s 1 forwards",
            },
        },
    },
    darkMode: "class",
    plugins: [],
};
