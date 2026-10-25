import { mdiGithub } from '@mdi/js';
import Icon from '@mdi/react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

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
        <footer className="app-footer">
            <span className="app-footer-text">{t(text)}</span>
            <div className="app-footer-links">
                <Link className="app-footer-copyright" to={homepage} target="_blank">
                    {`© ${currentYear} ${copyright}`}
                </Link>
                <Link to={repository} target="_blank">
                    <Icon className="app-footer-icon" path={mdiGithub} size={1} />
                </Link>
            </div>
        </footer>
    );
};
