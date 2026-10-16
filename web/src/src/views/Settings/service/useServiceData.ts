import { useEffect, useState } from 'react';

import { SettingsConstraints } from '../../../config/constraints';
import { useGetServiceDataQuery } from '../../../graphql';

export const useServiceData = () => {
    const {
        data: getServiceDataData,
        refetch: getServiceDataRefetch,
        error: getServiceDataError,
        loading: getServiceDataLoading
    } = useGetServiceDataQuery({ pollInterval: SettingsConstraints.pollInterval });

    const [serviceConfig, setServiceConfig] = useState<
        Record<
            string,
            {
                serviceName: string;
                constraints: Array<{
                    key: string;
                    name: string;
                    description: string;
                    configType: string;
                    isRequired: boolean;
                    currentValue: unknown;
                    options?: Record<string, unknown>;
                }>;
            }
        >
    >({});
    useEffect(() => {
        if (getServiceDataData?.getServiceConfigConstraint) {
            const configMap: typeof serviceConfig = {};
            getServiceDataData.getServiceConfigConstraint.forEach((service) => {
                const validConstraints = service.constraints.filter(
                    (constraint): constraint is NonNullable<typeof constraint> =>
                        constraint !== null
                );
                configMap[service.serviceId] = {
                    serviceName: service.serviceName,
                    constraints: validConstraints.map((constraint) => ({
                        key: constraint.key,
                        name: constraint.name,
                        description: constraint.description,
                        configType: constraint.configType,
                        isRequired: constraint.isRequired,
                        currentValue: constraint.currentValue,
                        options: constraint.options
                    }))
                };
            });
            setServiceConfig(configMap);
        }
    }, [getServiceDataData?.getServiceConfigConstraint]);

    const [serviceStatus, setServiceStatus] = useState<
        Record<
            string,
            {
                description: string;
                serviceName: string;
                restarts: number;
                startedAt: number;
                stoppedAt: number;
                updatedAt: number;
                isRunning: boolean;
            }
        >
    >({});
    useEffect(() => {
        if (getServiceDataData?.getServiceStatus) {
            const statusMap: typeof serviceStatus = {};
            getServiceDataData.getServiceStatus.forEach((service) => {
                statusMap[service.serviceId] = {
                    description: service.description,
                    serviceName: service.name,
                    restarts: service.restarts,
                    startedAt: service.startedAt,
                    stoppedAt: service.stoppedAt,
                    updatedAt: service.updatedAt,
                    isRunning: service.isRunning
                };
            });
            setServiceStatus(statusMap);
        }
    }, [getServiceDataData?.getServiceStatus]);

    return {
        getServiceDataRefetch,
        getServiceDataError,
        getServiceDataLoading,
        serviceConfig,
        serviceStatus
    };
};
