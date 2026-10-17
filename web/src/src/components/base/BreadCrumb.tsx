import { mdiChevronRight, mdiHome } from '@mdi/js';
import Icon from '@mdi/react';
import { Link } from 'react-router-dom';

interface IBreadCrumb {
    readonly title: string;
    readonly basename: string;
    readonly pathname: string;
}

export const BreadCrumb = ({ title, basename, pathname }: IBreadCrumb) => {
    return (
        <div className="app-breadcrumb">
            <ol className="app-breadcrumb-list">
                <li className="app-breadcrumb-home">
                    <Link className="app-flex" to={'/'}>
                        <Icon className="app-breadcrumb-icon" path={mdiHome} size={1} />
                        <span className="app-breadcrumb-separator">/</span>
                    </Link>
                </li>
                {pathname !== basename && (
                    <li className="app-flex">
                        <Icon className="app-breadcrumb-icon" path={mdiChevronRight} size={1} />
                        <Link className="app-breadcrumb-link" to={pathname}>
                            {title}
                        </Link>
                    </li>
                )}
            </ol>
        </div>
    );
};
