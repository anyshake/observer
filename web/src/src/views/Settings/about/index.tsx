import { mdiGithub, mdiOpenInNew } from '@mdi/js';
import Icon from '@mdi/react';
import { useTranslation } from 'react-i18next';

import { globalConfig } from '../../../config/global';
import { useGetSoftwareVersionQuery } from '../../../graphql';
import { ThirdPartyLicenses } from './ThirdPartyLicenses';
import { UpgradePanel } from './UpgradePanel';

export const About = ({ isAdmin }: { isAdmin: boolean }) => {
    const { t } = useTranslation();
    const { data, loading, error } = useGetSoftwareVersionQuery({ fetchPolicy: 'cache-first' });

    return (
        <div className="app-settings-content">
            <section className="app-software-panel" aria-labelledby="software-about-title">
                <div className="app-software-heading">
                    <img src={globalConfig.logo} alt="" className="app-software-logo" />
                    <div className="app-stack-sm">
                        <h2 id="software-about-title" className="app-subheading">
                            {t(globalConfig.name)}
                        </h2>
                        <p className="app-help-text">{t('views.Settings.About.description')}</p>
                    </div>
                </div>

                <dl className="app-software-facts">
                    <div>
                        <dt>{t('views.Settings.About.version')}</dt>
                        <dd className="app-break-all">
                            {loading
                                ? t('views.Settings.About.upgrade.loading')
                                : data?.getSoftwareVersion || '—'}
                        </dd>
                    </div>
                    <div>
                        <dt>{t('views.Settings.About.copyright')}</dt>
                        <dd>© 2025 {globalConfig.copyright}</dd>
                    </div>
                    <div>
                        <dt>{t('views.Settings.About.license')}</dt>
                        <dd>
                            <a
                                className="app-link-action"
                                href={`${globalConfig.repository}/blob/master/LICENSE`}
                                target="_blank"
                                rel="noopener noreferrer"
                            >
                                AGPL-3.0-or-later
                            </a>
                        </dd>
                    </div>
                    <div>
                        <dt>{t('views.Settings.About.inquiries')}</dt>
                        <dd>
                            <a className="app-link-action" href="mailto:anyshake@senseplex.net">
                                anyshake@senseplex.net
                            </a>
                        </dd>
                    </div>
                </dl>
                {error && (
                    <p role="alert" className="app-text-error">
                        {error.message}
                    </p>
                )}

                <div className="app-software-footer">
                    <p className="app-help-text">{t('views.Settings.About.license_note')}</p>
                    <div className="app-toolbar">
                        <a
                            className="btn btn-sm app-button-toolbar"
                            href={globalConfig.homepage}
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            <Icon
                                className="app-icon"
                                path={mdiOpenInNew}
                                size={0.7}
                                aria-hidden="true"
                            />
                            <span>{t('views.Settings.About.website')}</span>
                        </a>
                        <a
                            className="btn btn-sm app-button-toolbar"
                            href={globalConfig.repository}
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            <Icon
                                className="app-icon"
                                path={mdiGithub}
                                size={0.7}
                                aria-hidden="true"
                            />
                            <span>{t('views.Settings.About.source_code')}</span>
                        </a>
                        <ThirdPartyLicenses />
                    </div>
                </div>
            </section>

            {isAdmin && <UpgradePanel />}
        </div>
    );
};
