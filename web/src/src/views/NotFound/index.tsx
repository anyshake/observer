import { useTranslation } from 'react-i18next';

import { ErrorPage } from '../../components/ui/ErrorPage';

const NotFound = () => {
    const { t } = useTranslation();

    const handleGoBack = () => {
        window.history.back();
    };

    return (
        <div className="p-8">
            <ErrorPage
                code={404}
                heading={t('views.NotFound.page_not_found')}
                content={t('views.NotFound.check_your_url')}
                action={{
                    onClick: handleGoBack,
                    label: t('views.NotFound.go_back_button')
                }}
            />
        </div>
    );
};

export default NotFound;
