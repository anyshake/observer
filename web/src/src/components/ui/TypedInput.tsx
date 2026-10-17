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
                    className="select join-item app-select"
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
                    className="input join-item app-input app-pointer"
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
                    className="select join-item app-select"
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
                className="input join-item app-input"
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
                    className="btn join-item app-button-apply"
                    onClick={() => onSubmit?.(value)}
                    disabled={!hasChanged}
                >
                    <Icon className="app-icon" path={mdiCheck} size={0.9} />
                </button>
            </div>

            {isArrayType && (
                <DialogModal
                    open={isInputModalOpen}
                    onClose={() => setIsInputModalOpen(false)}
                    heading={
                        <h2 className="app-dialog-heading-colored">
                            <Icon className="app-icon" path={mdiTableEdit} size={1} />
                            <span>{t('components.TypedInput.editing_list', { fieldName })}</span>
                        </h2>
                    }
                >
                    <div className="join app-full-width">
                        <input
                            className="join-item input app-input app-full-width"
                            placeholder={t('components.TypedInput.add_new_item')}
                            value={tempArrayValue}
                            type={inputTypeMap[dataType as keyof typeof inputTypeMap]}
                            onChange={({ currentTarget }) => setTempArrayValue(currentTarget.value)}
                        />
                        <button className="btn join-item" onClick={handleAddArrayItem}>
                            <Icon className="app-icon" path={mdiPlus} size={0.9} />
                        </button>
                    </div>
                    <div className="app-array-table-container">
                        <table className="app-array-table">
                            <thead className="app-array-table-head">
                                <tr>
                                    <th className="app-array-table-index-heading">#</th>
                                    <th className="app-array-table-heading">
                                        {t('components.TypedInput.table_value')}
                                    </th>
                                    <th className="app-array-table-actions-heading">
                                        {t('components.TypedInput.table_actions')}
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="app-array-table-body">
                                {value.map((item: unknown, idx: number) => (
                                    <tr key={idx} className="app-array-table-row">
                                        <td className="app-array-table-index">{idx + 1}</td>
                                        <td className="app-array-table-value">
                                            <span className="app-array-table-text">
                                                {String(item)}
                                            </span>
                                        </td>
                                        <td className="app-array-table-cell">
                                            <div className="app-array-table-actions">
                                                <button
                                                    className="app-icon-action"
                                                    onClick={() => handleUpdateArrayItem(idx, 'up')}
                                                    disabled={idx === 0}
                                                >
                                                    <Icon
                                                        className="app-icon"
                                                        path={mdiMenuUp}
                                                        size={0.9}
                                                    />
                                                </button>
                                                <button
                                                    className="app-icon-action"
                                                    onClick={() =>
                                                        handleUpdateArrayItem(idx, 'down')
                                                    }
                                                    disabled={idx === value.length - 1}
                                                >
                                                    <Icon
                                                        className="app-icon"
                                                        path={mdiMenuDown}
                                                        size={0.9}
                                                    />
                                                </button>
                                                <button
                                                    className="app-icon-action-enabled"
                                                    onClick={() => handleRemoveArrayItem(idx)}
                                                >
                                                    <Icon
                                                        className="app-icon"
                                                        path={mdiDelete}
                                                        size={0.9}
                                                    />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </DialogModal>
            )}
        </>
    );
};
