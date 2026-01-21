import { mdiCheck, mdiDelete, mdiMenuDown, mdiMenuUp, mdiPlus, mdiTableEdit } from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from './DialogModal';

export type InputType = 'string' | 'string[]' | 'int' | 'int[]' | 'float' | 'float[]' | 'bool';

interface ITypedInput {
    readonly fieldName: string;
    readonly placeholder?: string;
    readonly dataType: InputType;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    readonly defaultValue?: any;
    readonly options?: Record<string, unknown>;
    readonly onSubmit?: (value: unknown) => void;
}

export const TypedInput = ({
    fieldName,
    options,
    placeholder,
    dataType,
    defaultValue,
    onSubmit
}: ITypedInput) => {
    const { t } = useTranslation();

    const inputTypeMap = useMemo(
        () => ({
            string: 'text',
            'string[]': 'text',
            int: 'number',
            'int[]': 'number',
            float: 'number',
            'float[]': 'number'
        }),
        []
    );

    const isArrayType = dataType.endsWith('[]');
    const baseType = dataType.replace('[]', '') as InputType;

    const [isInputModalOpen, setIsInputModalOpen] = useState(false);
    const [value, setValue] = useState(defaultValue ?? (isArrayType ? [] : ''));
    const [tempArrayValue, setTempArrayValue] = useState('');
    const [hasChanged, setHasChanged] = useState(false);

    const normalizeValue = useCallback(
        (val: unknown) => {
            if (val === null && dataType.endsWith('[]')) {
                return [];
            }
            if (dataType === 'string') {
                return (val as string).toLowerCase?.() ?? val;
            }
            if (dataType === 'string[]' && Array.isArray(val)) {
                return val.map((item: string) => item.toLowerCase?.() ?? item);
            }
            return val;
        },
        [dataType]
    );
    useEffect(() => {
        setHasChanged(
            JSON.stringify(normalizeValue(value)) !== JSON.stringify(normalizeValue(defaultValue))
        );
    }, [defaultValue, normalizeValue, value]);

    const parseValue = useCallback((input: string, type: InputType) => {
        switch (type) {
            case 'int':
                return parseInt(input);
            case 'float':
                return parseFloat(input);
            case 'bool':
                return input === 'true';
            default:
                return input;
        }
    }, []);

    const handleAddArrayItem = useCallback(() => {
        if (tempArrayValue === '') {
            return;
        }
        const parsed = parseValue(tempArrayValue, baseType);
        setValue([...value, parsed]);
        setTempArrayValue('');
    }, [baseType, parseValue, tempArrayValue, value]);

    const handleRemoveArrayItem = useCallback(
        (index: number) => {
            const newArray = [...value];
            newArray.splice(index, 1);
            setValue(newArray);
        },
        [value]
    );

    const handleUpdateArrayItem = useCallback(
        (index: number, direction: 'up' | 'down') => {
            const newArray = [...value];
            const temp = newArray[index];
            if (direction === 'up') {
                newArray[index] = newArray[index - 1];
                newArray[index - 1] = temp;
            } else {
                newArray[index] = newArray[index + 1];
                newArray[index + 1] = temp;
            }
            setValue(newArray);
        },
        [value]
    );

    const renderInput = useCallback(() => {
        if (options) {
            return (
                <select
                    className="select join-item border border-gray-300 shadow-sm focus:outline-none"
                    value={value}
                    onChange={({ currentTarget }) =>
                        setValue(parseValue(currentTarget.value, dataType))
                    }
                >
                    {Object.entries(options).map(([label, val]) => (
                        <option key={label} value={String(val)}>
                            {label}
                        </option>
                    ))}
                </select>
            );
        }

        if (isArrayType) {
            return (
                <input
                    className="input join-item cursor-pointer border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                    readOnly
                    placeholder={placeholder}
                    value={value.join(', ')}
                    onClick={() => setIsInputModalOpen(true)}
                />
            );
        }

        if (dataType === 'bool') {
            return (
                <select
                    className="select join-item border border-gray-300 shadow-sm focus:outline-none"
                    value={value}
                    onChange={({ currentTarget }) => setValue(currentTarget.value === 'true')}
                >
                    <option value="true">{t('components.TypedInput.yes')}</option>
                    <option value="false">{t('components.TypedInput.no')}</option>
                </select>
            );
        }

        return (
            <input
                className="input join-item border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                placeholder={placeholder}
                type={inputTypeMap[dataType]}
                value={inputTypeMap[dataType] === 'number' ? (isNaN(value) ? 0 : value) : value}
                onChange={({ currentTarget }) =>
                    setValue(parseValue(currentTarget.value, dataType))
                }
            />
        );
    }, [dataType, inputTypeMap, isArrayType, options, parseValue, placeholder, t, value]);

    return (
        <>
            <div className="join">
                {renderInput()}
                <button
                    className="btn join-item text-base-200 bg-purple-500 hover:bg-purple-600 disabled:bg-gray-300 disabled:text-gray-500"
                    onClick={() => onSubmit?.(value)}
                    disabled={!hasChanged}
                >
                    <Icon className="flex-shrink-0" path={mdiCheck} size={0.9} />
                </button>
            </div>

            {isArrayType && (
                <DialogModal
                    open={isInputModalOpen}
                    onClose={() => setIsInputModalOpen(false)}
                    heading={
                        <h2 className="flex items-center space-x-2 text-lg font-extrabold text-gray-800">
                            <Icon className="flex-shrink-0" path={mdiTableEdit} size={1} />
                            <span>{t('components.TypedInput.editing_list', { fieldName })}</span>
                        </h2>
                    }
                >
                    <div className="join w-full">
                        <input
                            className="join-item input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                            placeholder={t('components.TypedInput.add_new_item')}
                            value={tempArrayValue}
                            type={inputTypeMap[dataType as keyof typeof inputTypeMap]}
                            onChange={({ currentTarget }) => setTempArrayValue(currentTarget.value)}
                        />
                        <button className="btn join-item" onClick={handleAddArrayItem}>
                            <Icon className="flex-shrink-0" path={mdiPlus} size={0.9} />
                        </button>
                    </div>
                    <ul
                        className={`mt-6 ${dataType.startsWith('string') ? 'list-decimal' : 'list-disc'} space-y-2 px-4`}
                    >
                        {value.map((item: unknown, idx: number) => (
                            <li key={idx}>
                                <div className="flex flex-col items-center gap-2 sm:flex-row sm:justify-between">
                                    <span className="w-full break-words select-none sm:max-w-2/3">
                                        {String(item)}
                                    </span>
                                    <div className="flex gap-2">
                                        <button
                                            className="cursor-pointer text-gray-500 hover:text-gray-700 disabled:cursor-not-allowed disabled:text-gray-300"
                                            onClick={() => handleUpdateArrayItem(idx, 'up')}
                                            disabled={idx === 0}
                                        >
                                            <Icon
                                                className="flex-shrink-0"
                                                path={mdiMenuUp}
                                                size={0.9}
                                            />
                                        </button>
                                        <button
                                            className="cursor-pointer text-gray-500 hover:text-gray-700 disabled:cursor-not-allowed disabled:text-gray-300"
                                            onClick={() => handleUpdateArrayItem(idx, 'down')}
                                            disabled={idx === value.length - 1}
                                        >
                                            <Icon
                                                className="flex-shrink-0"
                                                path={mdiMenuDown}
                                                size={0.9}
                                            />
                                        </button>
                                        <button
                                            className="cursor-pointer text-gray-500 hover:text-gray-700"
                                            onClick={() => handleRemoveArrayItem(idx)}
                                        >
                                            <Icon
                                                className="flex-shrink-0"
                                                path={mdiDelete}
                                                size={0.9}
                                            />
                                        </button>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </DialogModal>
            )}
        </>
    );
};
