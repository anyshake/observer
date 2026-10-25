import { mdiLicense, mdiRefresh } from '@mdi/js';
import Icon from '@mdi/react';
import { lazy, Suspense, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';

const Markdown = lazy(() =>
    import('../../../components/widget/Markdown').then(({ Markdown }) => ({ default: Markdown }))
);

const LicenseContent = ({ fileName }: { fileName: string }) => {
    const { t } = useTranslation();
    const [text, setText] = useState<string>();
    const [failed, setFailed] = useState(false);
    const [attempt, setAttempt] = useState(0);

    useEffect(() => {
        const controller = new AbortController();
        setFailed(false);
        setText(undefined);

        const load = async () => {
            try {
                const response = await fetch(new URL(fileName, document.baseURI), {
                    signal: controller.signal,
                    cache: 'no-cache'
                });
                if (!response.ok || response.headers.get('Content-Type')?.includes('text/html')) {
                    throw new Error('License report unavailable');
                }
                const content = await response.text();
                if (!content.trim()) {
                    throw new Error('License report is empty');
                }
                if (!controller.signal.aborted) {
                    setText(content);
                }
            } catch {
                if (!controller.signal.aborted) {
                    setFailed(true);
                }
            }
        };

        void load();
        return () => controller.abort();
    }, [attempt, fileName]);

    return (
        <div className="app-stack-sm">
            {failed ? (
                <div className="app-stack-sm">
                    <p role="alert" className="app-text-error">
                        {t('views.Settings.About.third_party.load_error')}
                    </p>
                    <button
                        type="button"
                        className="btn btn-sm app-button-toolbar"
                        onClick={() => setAttempt((value) => value + 1)}
                    >
                        <Icon
                            className="app-icon"
                            path={mdiRefresh}
                            size={0.7}
                            aria-hidden="true"
                        />
                        <span>{t('views.Settings.About.third_party.retry')}</span>
                    </button>
                </div>
            ) : text === undefined ? (
                <p role="status" className="app-help-text">
                    {t('views.Settings.About.third_party.loading')}
                </p>
            ) : (
                <Suspense
                    fallback={
                        <p className="app-help-text">
                            {t('views.Settings.About.third_party.loading')}
                        </p>
                    }
                >
                    <Markdown className="prose app-license-markdown">{text}</Markdown>
                </Suspense>
            )}
        </div>
    );
};

const LicenseSection = ({ fileName, title }: { fileName: string; title: string }) => {
    const [expanded, setExpanded] = useState(false);

    return (
        <details
            className="collapse-arrow app-license-section collapse"
            onToggle={({ currentTarget }) => setExpanded(currentTarget.open)}
        >
            <summary className="collapse-title app-license-heading">{title}</summary>
            <div className="collapse-content">
                {expanded && <LicenseContent fileName={fileName} />}
            </div>
        </details>
    );
};

export const ThirdPartyLicenses = () => {
    const { t } = useTranslation();
    const [open, setOpen] = useState(false);

    return (
        <>
            <button
                type="button"
                className="btn btn-sm app-button-toolbar"
                aria-haspopup="dialog"
                onClick={() => setOpen(true)}
            >
                <Icon className="app-icon" path={mdiLicense} size={0.7} aria-hidden="true" />
                <span>{t('views.Settings.About.third_party.title')}</span>
            </button>
            <DialogModal
                open={open}
                onClose={() => setOpen(false)}
                heading={
                    <div className="app-stack-sm">
                        <h2 className="app-row app-subheading">
                            <Icon
                                className="app-icon"
                                path={mdiLicense}
                                size={1}
                                aria-hidden="true"
                            />
                            <span>{t('views.Settings.About.third_party.title')}</span>
                        </h2>
                        <p className="app-help-text">
                            {t('views.Settings.About.third_party.description')}
                        </p>
                    </div>
                }
            >
                {open && (
                    <div className="app-license-list">
                        <LicenseSection
                            fileName="third-party-go-licenses.md"
                            title={t('views.Settings.About.third_party.backend')}
                        />
                        <LicenseSection
                            fileName="third-party-web-licenses.md"
                            title={t('views.Settings.About.third_party.web')}
                        />
                    </div>
                )}
            </DialogModal>
        </>
    );
};
