import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import type { GetHomeDataQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { healthStatusClass } from './healthStatus';

interface IServiceTable {
    readonly services: GetHomeDataQuery['getServiceStatus'];
    readonly stale: boolean;
}

// Renders a table on wide containers and stacked rows on narrow ones (see home.css).
export const ServiceTable = ({ services, stale }: IServiceTable) => {
    const { t } = useTranslation();
    const sorted = useMemo(
        () => [...(services ?? [])].sort((a, b) => a.name.localeCompare(b.name)),
        [services]
    );
    const timestamp = (value: number) => (value > 0 ? getTimeString(value) : '—');

    if (sorted.length === 0) {
        return <p className="app-help-text">{t('views.Home.service_status.empty')}</p>;
    }

    const labels = {
        restarts: t('views.Home.service_status.restarts'),
        startedAt: t('views.Home.service_status.started_at'),
        stoppedAt: t('views.Home.service_status.stopped_at')
    };

    return (
        <div className="app-service-list-container">
            <div className="app-service-list">
                <div className="app-service-head" aria-hidden>
                    <span>{t('views.Home.service_status.name')}</span>
                    <span>{t('views.Home.service_status.state')}</span>
                    <span>{labels.restarts}</span>
                    <span>{labels.startedAt}</span>
                    <span>{labels.stoppedAt}</span>
                </div>
                <ul>
                    {sorted.map((service) => {
                        const level = stale ? 'unknown' : service.isRunning ? 'ok' : 'error';
                        return (
                            <li key={service.serviceId} className="app-service-row">
                                <div className="app-service-name">
                                    <div className="app-text-medium">{service.name}</div>
                                    <div className="app-status-description">
                                        {service.description}
                                    </div>
                                </div>
                                <span className="app-service-state">
                                    <span
                                        className={`status ${healthStatusClass[level]}`}
                                        aria-hidden
                                    />
                                    {t(
                                        stale
                                            ? 'views.Home.health_level.unknown'
                                            : service.isRunning
                                              ? 'views.Home.service_status.running'
                                              : 'views.Home.service_status.stopped'
                                    )}
                                </span>
                                <dl className="app-service-meta">
                                    <div>
                                        <dt>{labels.restarts}</dt>
                                        <dd>{service.restarts}</dd>
                                    </div>
                                    <div>
                                        <dt>{labels.startedAt}</dt>
                                        <dd>{timestamp(service.startedAt)}</dd>
                                    </div>
                                    <div>
                                        <dt>{labels.stoppedAt}</dt>
                                        <dd>{timestamp(service.stoppedAt)}</dd>
                                    </div>
                                </dl>
                            </li>
                        );
                    })}
                </ul>
            </div>
        </div>
    );
};
