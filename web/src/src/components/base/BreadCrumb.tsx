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
        <div className="bg-base-100 rounded-lg px-5 py-3 shadow-md">
            <ol className="flex space-x-2 text-sm font-medium text-gray-700">
                <li className="cursor-pointer hover:text-gray-900">
                    <Link className="flex" to={'/'}>
                        <Icon className="mr-2 flex-shrink-0 self-center" path={mdiHome} size={1} />
                        <span className="my-2">/</span>
                    </Link>
                </li>
                {pathname !== basename && (
                    <li className="flex">
                        <Icon
                            className="mr-2 flex-shrink-0 self-center"
                            path={mdiChevronRight}
                            size={1}
                        />
                        <Link className="my-2 cursor-pointer hover:text-gray-900" to={pathname}>
                            {title}
                        </Link>
                    </li>
                )}
            </ol>
        </div>
    );
};
