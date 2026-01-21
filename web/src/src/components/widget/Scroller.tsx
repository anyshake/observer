import { mdiArrowUp } from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useState } from 'react';

interface IScroller {
    readonly threshold: number;
}

export const Scroller = ({ threshold }: IScroller) => {
    const [showButton, setShowButton] = useState(false);

    const scrollToTop = useCallback(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }, []);

    const toggleButton = useCallback(() => {
        setShowButton(window.scrollY > threshold);
    }, [threshold]);

    useEffect(() => {
        document.addEventListener('scroll', toggleButton);
        return () => {
            document.removeEventListener('scroll', toggleButton);
        };
    }, [toggleButton]);

    return (
        <button
            className={`right-3 bottom-16 flex size-10 items-center justify-center rounded-full bg-purple-500 text-white duration-300 hover:cursor-pointer hover:bg-purple-600 ${
                showButton ? 'animate-fade-left animate-duration-300 fixed' : 'hidden'
            }`}
            onClick={scrollToTop}
        >
            <Icon className="flex-shrink-0" path={mdiArrowUp} size={1} />
        </button>
    );
};
