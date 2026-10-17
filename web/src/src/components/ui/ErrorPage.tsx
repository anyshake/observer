import { mdiBugPause } from '@mdi/js';
import Icon from '@mdi/react';
import { useEffect, useMemo, useState } from 'react';
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

    const currentTime = useMemo(() => Date.now(), []);
    const [isDebug, setIsDebug] = useState(false);

    useEffect(() => {
        hideLoaderAnimation();
    }, []);

    return (
        <div className="app-error-page">
            <h1 className="app-error-code">{code ?? ':-('}</h1>
            <p className="app-error-title">
                {heading ?? t('components.ErrorPage.something_went_wrong')}
            </p>
            <p className="app-text-muted">{content ?? t('components.ErrorPage.try_again_later')}</p>

            {action && (
                <button onClick={action.onClick} className="btn btn-md app-error-back">
                    {action.label}
                </button>
            )}

            {debug && (
                <button
                    className="app-error-details-button"
                    onClick={() => {
                        setIsDebug(!isDebug);
                    }}
                >
                    <Icon className="app-error-icon" path={mdiBugPause} size={1} />
                </button>
            )}

            {debug && isDebug && (
                <div className="app-error-details">
                    <CodeBlock language="javascript" fileName={`error_stack_${currentTime}.log`}>
                        {debug}
                    </CodeBlock>
                </div>
            )}
        </div>
    );
};
