import { mdiAlertCircle, mdiCheckCircle, mdiSquareEditOutline } from '@mdi/js';
import Icon from '@mdi/react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import { InputType, TypedInput } from '../../../components/ui/TypedInput';
import type { useServiceConfig } from './useServiceConfig';
import type { useServiceData } from './useServiceData';

export const ServiceConfigModal = ({
    serviceConfig,
    serviceStatus,
    configure
}: {
    serviceConfig: ReturnType<typeof useServiceData>['serviceConfig'];
    serviceStatus: ReturnType<typeof useServiceData>['serviceStatus'];
    configure: ReturnType<typeof useServiceConfig>;
}) => {
    const { t } = useTranslation();
    const {
        selectedServiceId,
        isConfigureModalOpen,
        setIsConfigureModalOpen,
        configureModalStatus,
        setConfigureModalStatus,
        handleUpdateServiceConfig
    } = configure;

    return (
        <>
            {serviceConfig[selectedServiceId] && (
                <DialogModal
                    fullScreen
                    heading={
                        <div className="app-text-heading">
                            <h2 className="app-service-config-heading">
                                <Icon className="app-icon" path={mdiSquareEditOutline} size={1} />
                                <span>
                                    {t('views.Settings.Service.edit_service_modal.title', {
                                        name: serviceConfig[selectedServiceId].serviceName
                                    })}
                                </span>
                            </h2>
                            <p className="app-text-small">
                                {serviceStatus[selectedServiceId].description}
                            </p>
                        </div>
                    }
                    open={isConfigureModalOpen}
                    onClose={() => {
                        setIsConfigureModalOpen(false);
                        setConfigureModalStatus({ message: '' });
                    }}
                >
                    {configureModalStatus.status === 'updating' && (
                        <div className="app-form-status-muted">
                            <span className="loading loading-spinner loading-xs" />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'error' && (
                        <div className="app-form-status-error">
                            <Icon className="app-icon" path={mdiAlertCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'success' && (
                        <div className="app-form-status-success">
                            <Icon className="app-icon" path={mdiCheckCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    <div className="app-field-grid">
                        {serviceConfig[selectedServiceId].constraints.map((item, index) => (
                            <fieldset
                                className="fieldset app-fieldset"
                                key={`${selectedServiceId}-${item.key}-${index}`}
                            >
                                <legend className="fieldset-legend">
                                    <span className="app-field-name">{item.name}</span>
                                    {item.isRequired && <span className="app-text-error">*</span>}
                                </legend>
                                <TypedInput
                                    key={`${item.key}-${item.currentValue}`}
                                    dataType={item.configType as InputType}
                                    placeholder={item.name}
                                    defaultValue={item.currentValue as string}
                                    onSubmit={(value) =>
                                        handleUpdateServiceConfig(
                                            selectedServiceId,
                                            item.key,
                                            value,
                                            item.isRequired
                                        )
                                    }
                                    fieldName={item.name}
                                    options={item.options}
                                />
                                <p className="app-field-help">{item.description}</p>
                            </fieldset>
                        ))}
                    </div>
                </DialogModal>
            )}
        </>
    );
};
