import { mdiGithub } from '@mdi/js';
import Icon from '@mdi/react';
import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

interface IFooter {
    readonly copyright: string;
    readonly homepage: string;
    readonly repository: string;
    readonly text: string;
}

export const Footer = ({ text, homepage, copyright, repository }: IFooter) => {
    const [t] = useTranslation();
    const currentYear = useMemo(() => new Date().getFullYear(), []);

    return (
        <footer className="bg-base-300 flex w-full flex-col justify-between px-8 py-2 text-gray-500 sm:flex-row">
            <span className="self-center text-center text-xs italic sm:ml-12">{t(text)}</span>
            <div className="flex items-center justify-center text-center">
                <Link
                    className="self-center text-sm hover:underline"
                    to={homepage}
                    target="_blank"
                >{`© ${currentYear} ${copyright}`}</Link>
                <Link to={repository} target="_blank">
                    <Icon
                        className="mx-2 flex-shrink-0 self-center hover:opacity-80"
                        path={mdiGithub}
                        size={1}
                    />
                </Link>
            </div>
        </footer>
    );
};
