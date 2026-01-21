import { mdiBugPause } from '@mdi/js';
import Icon from '@mdi/react';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { hideLoaderAnimation } from '../../helpers/app/hideLoaderAnimation';
import { CodeBlock } from '../widget/CodeBlock';

interface IErrorPage {
    readonly code?: number;
    readonly heading?: string;
    readonly content?: string;
    readonly action?: {
        readonly label: string;
        readonly onClick: () => void;
    };
    readonly debug?: string;
}

export const ErrorPage = ({ code, heading, content, action, debug }: IErrorPage) => {
    const { t } = useTranslation();

    const [currentTime] = useState(Date.now());
    const [isDebug, setIsDebug] = useState(false);

    useEffect(() => {
        hideLoaderAnimation();
    }, []);

    return (
        <div className="flex min-h-screen flex-col items-center justify-center space-y-4 p-5">
            <h1 className="text-6xl font-bold tracking-tight text-gray-800">{code ?? ':-('}</h1>
            <p className="text-lg font-medium text-gray-600">
                {heading ?? t('components.ErrorPage.something_went_wrong')}
            </p>
            <p className="text-gray-500">{content ?? t('components.ErrorPage.try_again_later')}</p>

            {action && (
                <button
                    onClick={action.onClick}
                    className="btn btn-md rounded-md bg-gray-500 px-4 py-2 text-gray-200 transition-all hover:bg-gray-600"
                >
                    {action.label}
                </button>
            )}

            {debug && (
                <button
                    className="cursor-pointer rounded-full bg-gray-100 p-2 transition-all hover:scale-110"
                    onClick={() => {
                        setIsDebug(!isDebug);
                    }}
                >
                    <Icon className="flex-shrink-0 text-gray-400" path={mdiBugPause} size={1} />
                </button>
            )}

            {debug && isDebug && (
                <div className="container w-[350px] md:w-[600px]">
                    <CodeBlock language="javascript" fileName={`error_stack_${currentTime}.log`}>
                        {debug}
                    </CodeBlock>
                </div>
            )}
        </div>
    );
};
