import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useUpdateServiceConfigMutation } from '../../../graphql';

export const useServiceConfig = (getServiceDataRefetch: () => Promise<unknown>) => {
    const { t } = useTranslation();
    const [isConfigureModalOpen, setIsConfigureModalOpen] = useState(false);
    const [configureModalStatus, setConfigureModalStatus] = useState<{
        message: string;
        status?: 'updating' | 'success' | 'error';
    }>({ message: '' });
    const [selectedServiceId, setSelectedServiceId] = useState('');
    const [updateServiceConfigMutation] = useUpdateServiceConfigMutation();
    const handleUpdateServiceConfig = useCallback(
        async (serviceId: string, key: string, val: unknown, isRequired: boolean) => {
            if (isRequired && !val && typeof val !== 'boolean') {
                if (val !== 0 && typeof val === 'number') {
                    setConfigureModalStatus({ status: 'error', message: 'This field is required' });
                    return;
                }
            }
            try {
                setConfigureModalStatus({
                    status: 'updating',
                    message: t('views.Settings.Service.edit_service.updating')
                });
                await updateServiceConfigMutation({ variables: { serviceId, key, val } });
                await getServiceDataRefetch();
                setConfigureModalStatus({
                    status: 'success',
                    message: t('views.Settings.Service.edit_service.success')
                });
            } catch (error) {
                setConfigureModalStatus({
                    status: 'error',
                    message: t('views.Settings.Service.edit_service.error', { error })
                });
            }
            await new Promise((resolve) => setTimeout(resolve, 5000));
            setConfigureModalStatus({ message: '' });
        },
        [t, getServiceDataRefetch, updateServiceConfigMutation]
    );

    return {
        isConfigureModalOpen,
        setIsConfigureModalOpen,
        configureModalStatus,
        setConfigureModalStatus,
        selectedServiceId,
        setSelectedServiceId,
        handleUpdateServiceConfig
    };
};
