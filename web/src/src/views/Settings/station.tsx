import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { ErrorPage } from '../../components/ui/ErrorPage';
import { InputType, TypedInput } from '../../components/ui/TypedInput';
import { useGetStationConfigQuery, useUpdateStationConfigMutation } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';

export const Station = () => {
    const { t } = useTranslation();

    const {
        data: getStationConfigData,
        refetch: getStationConfigRefetch,
        error: getStationConfigError,
        loading: getStationConfigLoading
    } = useGetStationConfigQuery();
    const [stationConfig, setStationConfig] = useState<
        Array<{
            name: string;
            key: string;
            description: string;
            namespace: string;
            configType: InputType;
            isRequired: boolean;
            currentValue: unknown;
            options?: Record<string, unknown>;
        }>
    >();
    useEffect(() => {
        if (getStationConfigData?.getStationConfigConstraint) {
            setStationConfig(
                getStationConfigData.getStationConfigConstraint.map((i) => ({
                    name: i.name,
                    description: i.description,
                    namespace: i.namespace,
                    key: i.key,
                    configType: i.configType as InputType,
                    isRequired: i.isRequired,
                    currentValue: i.currentValue,
                    options: i.options
                }))
            );
        }
    }, [getStationConfigData]);

    const [updateStationConfig] = useUpdateStationConfigMutation();
    const handleUpdateStationConfig = useCallback(
        async (key: string, val: unknown, isRequired: boolean) => {
            if (isRequired && !val && typeof val !== 'boolean') {
                if (val !== 0 && typeof val === 'number') {
                    sendUserAlert(t('views.Settings.Station.update_station.required'), true);
                    return;
                }
            }
            const requestFn = async () => {
                await updateStationConfig({ variables: { key, value: val } });
                await getStationConfigRefetch();
            };
            await sendPromiseAlert(
                requestFn(),
                t('views.Settings.Station.update_station.updating'),
                t('views.Settings.Station.update_station.success'),
                (error) => t('views.Settings.Station.update_station.error', { error })
            );
        },
        [t, getStationConfigRefetch, updateStationConfig]
    );

    return getStationConfigError ? (
        <ErrorPage
            content={getStationConfigError.message}
            debug={JSON.stringify(getStationConfigError)}
        />
    ) : getStationConfigLoading ? (
        <div className="flex items-center justify-center">
            <span className="loading loading-spinner text-primary mt-6" />
        </div>
    ) : (
        <div className="flex flex-wrap justify-center gap-4 lg:justify-start">
            {stationConfig?.map((item, index) => (
                <fieldset className="fieldset w-md" key={`${item.namespace}-${item.key}-${index}`}>
                    <legend className="fieldset-legend">
                        <span className="font-bold text-gray-700">{item.name}</span>
                        {item.isRequired && <span className="text-red-500">*</span>}
                    </legend>
                    <TypedInput
                        dataType={item.configType}
                        placeholder={item.name}
                        defaultValue={item.currentValue as string}
                        onSubmit={(value) =>
                            handleUpdateStationConfig(item.key, value, item.isRequired)
                        }
                        fieldName={item.name}
                        options={item.options}
                    />
                    <p className="w-[80%] text-xs text-gray-500">{item.description}</p>
                </fieldset>
            ))}
        </div>
    );
};
