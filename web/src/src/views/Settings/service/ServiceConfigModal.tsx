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
                        <div className="text-gray-800">
                            <h2 className="flex items-center space-x-2 py-4 text-lg font-extrabold">
                                <Icon
                                    className="flex-shrink-0"
                                    path={mdiSquareEditOutline}
                                    size={1}
                                />
                                <span>
                                    {t('views.Settings.Service.edit_service_modal.title', {
                                        name: serviceConfig[selectedServiceId].serviceName
                                    })}
                                </span>
                            </h2>
                            <p className="text-sm">
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
                        <div className="my-2 flex items-center space-x-2 text-sm text-gray-500">
                            <span className="loading loading-spinner loading-xs" />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'error' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-red-500">
                            <Icon className="flex-shrink-0" path={mdiAlertCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'success' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-green-500">
                            <Icon className="flex-shrink-0" path={mdiCheckCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    <div className="flex flex-wrap justify-center gap-4 lg:justify-start">
                        {serviceConfig[selectedServiceId].constraints.map((item, index) => (
                            <fieldset
                                className="fieldset w-md"
                                key={`${selectedServiceId}-${item.key}-${index}`}
                            >
                                <legend className="fieldset-legend">
                                    <span className="font-bold text-gray-700">{item.name}</span>
                                    {item.isRequired && <span className="text-red-500">*</span>}
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
                                <p className="w-[80%] text-xs text-gray-500">{item.description}</p>
                            </fieldset>
                        ))}
                    </div>
                </DialogModal>
            )}
        </>
    );
};
