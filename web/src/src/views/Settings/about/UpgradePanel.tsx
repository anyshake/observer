import { mdiDownload, mdiRefresh } from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { useGetUpgradeStatusQuery, useTriggerUpgradeMutation } from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';

export const UpgradePanel = () => {
    const { t } = useTranslation();
    const { data, loading, error, refetch } = useGetUpgradeStatusQuery({
        fetchPolicy: 'cache-and-network',
        notifyOnNetworkStatusChange: true
    });
    const [triggerUpgrade, { loading: triggering }] = useTriggerUpgradeMutation();

    const status = data?.getUpgradeStatus;
    const supported = status !== null && status !== undefined;

    const handleCheck = useCallback(() => {
        refetch({ fresh: true }).catch(() => undefined);
    }, [refetch]);

    const handleUpgrade = useCallback(() => {
        sendUserConfirm(t('views.Settings.About.upgrade.confirm_message'), {
            title: t('views.Settings.About.upgrade.confirm_title'),
            cancelBtnText: t('views.Settings.About.upgrade.cancel_button'),
            confirmBtnText: t('views.Settings.About.upgrade.confirm_button'),
            onConfirmed: () => {
                sendPromiseAlert(
                    triggerUpgrade(),
                    t('views.Settings.About.upgrade.queued'),
                    t('views.Settings.About.upgrade.queued'),
                    (err) => (err instanceof Error ? err.message : String(err))
                );
            }
        });
    }, [t, triggerUpgrade]);

    const summary = (() => {
        if (!status) {
            return '';
        }
        if (status.applied) {
            return t('views.Settings.About.upgrade.applied_yes');
        }
        return status.eligible
            ? t('views.Settings.About.upgrade.eligible_yes')
            : t('views.Settings.About.upgrade.eligible_no');
    })();

    const upgradeEnabled = Boolean(status?.eligible) && !status?.applied && !triggering;

    return (
        <section className="app-software-panel" aria-labelledby="software-upgrade-title">
            <h3 id="software-upgrade-title" className="app-subheading">
                {t('views.Settings.About.upgrade.title')}
            </h3>

            {error ? (
                <p role="alert" className="app-text-error">
                    {error.message}
                </p>
            ) : !supported && !loading ? (
                <p className="app-help-text">{t('views.Settings.About.upgrade.unsupported')}</p>
            ) : (
                <>
                    <dl className="app-software-facts">
                        <div>
                            <dt>{t('views.Settings.About.upgrade.current')}</dt>
                            <dd className="app-break-all">{status?.current || '—'}</dd>
                        </div>
                        <div>
                            <dt>{t('views.Settings.About.upgrade.latest')}</dt>
                            <dd className="app-break-all">{status?.latest || '—'}</dd>
                        </div>
                        <div>
                            <dt>{t('views.Settings.About.upgrade.required')}</dt>
                            <dd className="app-break-all">{status?.required || '—'}</dd>
                        </div>
                        <div>
                            <dt>{t('views.Settings.About.upgrade.eligible')}</dt>
                            <dd>
                                {status?.eligible
                                    ? t('views.Settings.About.upgrade.yes')
                                    : t('views.Settings.About.upgrade.no')}
                            </dd>
                        </div>
                        <div>
                            <dt>{t('views.Settings.About.upgrade.applied')}</dt>
                            <dd>
                                {status?.applied
                                    ? t('views.Settings.About.upgrade.yes')
                                    : t('views.Settings.About.upgrade.no')}
                            </dd>
                        </div>
                    </dl>
                    {status?.checkFailed ? (
                        <p role="alert" className="app-text-error">
                            {t('views.Settings.About.upgrade.check_failed')}
                        </p>
                    ) : (
                        summary && <p className="app-help-text">{summary}</p>
                    )}
                </>
            )}

            <div className="app-software-footer">
                <div className="app-toolbar">
                    <button
                        type="button"
                        className="btn btn-sm app-button-toolbar"
                        disabled={loading}
                        onClick={handleCheck}
                    >
                        <Icon
                            className="app-icon"
                            path={mdiRefresh}
                            size={0.7}
                            aria-hidden="true"
                        />
                        <span>
                            {loading
                                ? t('views.Settings.About.upgrade.checking')
                                : t('views.Settings.About.upgrade.check_button')}
                        </span>
                    </button>
                    <button
                        type="button"
                        className="btn btn-sm app-button-toolbar"
                        disabled={!upgradeEnabled}
                        onClick={handleUpgrade}
                    >
                        <Icon
                            className="app-icon"
                            path={mdiDownload}
                            size={0.7}
                            aria-hidden="true"
                        />
                        <span>{t('views.Settings.About.upgrade.upgrade_button')}</span>
                    </button>
                </div>
            </div>
        </section>
    );
};
