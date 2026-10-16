import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { type GetHomeDataQuery, useGetUpgradeStatusLazyQuery } from '../../graphql';
import { isOlderThanRelease } from '../../helpers/app/isOlderThanRelease';

export const useUpgradeMessage = (getHomeDataData: GetHomeDataQuery | undefined) => {
    const { t } = useTranslation();
    const [getUpgradeStatus] = useGetUpgradeStatusLazyQuery();
    const [upgradeMessage, setUpgradeMessage] = useState('');

    const generateUpgradeMessage = useCallback(
        (
            current: string,
            latest: string,
            required: string,
            eligible: boolean,
            applied: boolean,
            checkFailed: boolean
        ) => {
            if (checkFailed) {
                return t('views.Home.upgrade.check_failed');
            }
            if (eligible) {
                return t('views.Home.upgrade.update_available', { latest });
            }
            if (applied) {
                return t('views.Home.upgrade.restart_needed', { latest, current });
            }
            if (isOlderThanRelease(current, latest) && isOlderThanRelease(current, required)) {
                return t('views.Home.upgrade.manual_upgrade_needed', { latest, current, required });
            }

            return '';
        },
        [t]
    );

    useEffect(() => {
        if (getHomeDataData?.getCurrentUser.admin) {
            getUpgradeStatus()
                .then(({ data, error }) => {
                    if (error) {
                        setUpgradeMessage(t('views.Home.upgrade.check_failed'));
                        return;
                    }
                    let msg = '';
                    if (data?.getUpgradeStatus) {
                        const { current, latest, required, eligible, applied, checkFailed } =
                            data.getUpgradeStatus;
                        msg = generateUpgradeMessage(
                            current,
                            latest,
                            required,
                            eligible,
                            applied,
                            checkFailed
                        );
                    }
                    setUpgradeMessage(msg);
                })
                .catch(() => setUpgradeMessage(t('views.Home.upgrade.check_failed')));
        }
    }, [generateUpgradeMessage, getHomeDataData?.getCurrentUser.admin, getUpgradeStatus, t]);

    return { upgradeMessage };
};
