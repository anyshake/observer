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
            className={`app-scroller ${showButton ? 'app-scroller-visible' : 'app-hidden'}`}
            onClick={scrollToTop}
        >
            <Icon className="app-icon" path={mdiArrowUp} size={1} />
        </button>
    );
};
