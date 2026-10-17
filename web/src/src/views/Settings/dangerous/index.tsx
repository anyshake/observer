import { useTranslation } from 'react-i18next';

import { Banner } from '../../../components/widget/Banner';
import { DangerousAction } from './action';
import { CleanupModal } from './CleanupModal';
import { ExecutableBackups } from './executableBackups';
import { useCleanup } from './useCleanup';
import { useDangerousActions } from './useDangerousActions';

interface IDangerous {
    currentLocale: string;
}

export const Dangerous = ({ currentLocale }: IDangerous) => {
    const { t } = useTranslation();
    const cleanup = useCleanup();
    const {
        cleanupTargets,
        isCleanupRunning,
        runningTarget,
        cleanupStatusMeta,
        handleCleanupTargetClick
    } = cleanup;
    const { otherActions, handleOtherActionClick } = useDangerousActions();

    return (
        <div className="app-settings-content-padded">
            <Banner
                fullWidth={true}
                status={cleanupStatusMeta.status}
                message={cleanupStatusMeta.message}
            />

            <div>
                <ExecutableBackups currentLocale={currentLocale} />
                {cleanupTargets.map(({ id, title, description, buttonText }) => {
                    const isActiveTarget = isCleanupRunning && runningTarget === id;

                    return (
                        <DangerousAction
                            key={id}
                            title={title}
                            description={description}
                            buttonText={buttonText}
                            value={id}
                            disabled={isCleanupRunning}
                            activeLabel={
                                isActiveTarget
                                    ? t('views.Settings.Dangerous.purge_data.mode_running')
                                    : undefined
                            }
                            onClick={handleCleanupTargetClick}
                        />
                    );
                })}
                {otherActions.map(({ title, description, buttonText }, index) => (
                    <DangerousAction
                        key={`${index}-${title}`}
                        title={title}
                        description={description}
                        buttonText={buttonText}
                        value={index}
                        onClick={handleOtherActionClick}
                    />
                ))}
            </div>

            <CleanupModal currentLocale={currentLocale} cleanup={cleanup} />
        </div>
    );
};
