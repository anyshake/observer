import { useCallback, useEffect, useState } from "react";
import ScrollIcon from "../assets/icons/arrow-up-solid.svg";

interface ScrollerProps {
    readonly threshold: number;
}

export const Scroller = (props: ScrollerProps) => {
    const { threshold = 100 } = props;
    const [showButton, setShowButton] = useState(false);

    const scrollToTop = () => window.scrollTo({ top: 0, behavior: "smooth" });

    const toggleButton = useCallback(
        () => setShowButton(window.scrollY > threshold),
        [threshold]
    );

    useEffect(() => {
        document.addEventListener("scroll", toggleButton);
        return () => document.removeEventListener("scroll", toggleButton);
    }, [toggleButton]);

    return (
        <button
            className={`bg-purple-500 hover:bg-purple-600 duration-300 size-10 rounded-full bottom-16 right-3 flex justify-center items-center ${
                showButton
                    ? "fixed animate-fade-left animate-duration-300"
                    : "hidden"
            }`}
            onClick={scrollToTop}
        >
            <img className="size-4" src={ScrollIcon} alt="" />
        </button>
    );
};
