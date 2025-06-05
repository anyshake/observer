import { useMemo, useState } from 'react';

import { CodeBlock } from '../../components/CodeBlock';
import { ErrorPage } from '../../components/ErrorPage';
import { useGetStationMetadataQuery } from '../../graphql';

export const Metadata = () => {
    const availableFormats = useMemo(
        () => [
            {
                value: 'seiscomp_xml',
                label: 'SeisComP XML',
                fileName: 'seiscomp_xml.xml',
                language: 'xml'
            },
            {
                value: 'station_xml',
                label: 'FDSN StationXML',
                fileName: 'fdsn_station_xml.xml',
                language: 'xml'
            }
        ],
        []
    );
    const [formatIndex, setFormatIndex] = useState(0);
    const {
        data: getStationMetadataData,
        loading: getStationMetadataLoading,
        error: getStationMetadataError
    } = useGetStationMetadataQuery({ variables: { format: availableFormats[formatIndex].value } });

    return getStationMetadataError ? (
        <ErrorPage
            content={getStationMetadataError.message}
            debug={JSON.stringify(getStationMetadataError)}
        />
    ) : (
        <div className="mx-auto max-w-3xl space-y-4">
            <div className="flex flex-wrap gap-4">
                {availableFormats.map(({ value, label }, index) => {
                    const inputId = `format-radio-${value}`;
                    return (
                        <div key={`${index}-${value}`} className="flex items-center gap-2">
                            <input
                                id={inputId}
                                type="radio"
                                value={value}
                                checked={formatIndex === index}
                                onChange={() => setFormatIndex(index)}
                                className="radio radio-primary radio-xs"
                            />
                            <label htmlFor={inputId} className="label cursor-pointer">
                                {label}
                            </label>
                        </div>
                    );
                })}
            </div>

            {getStationMetadataLoading ? (
                <span className="loading loading-spinner text-primary" />
            ) : (
                <CodeBlock
                    fileName={availableFormats[formatIndex].fileName}
                    language={availableFormats[formatIndex].language}
                >
                    {getStationMetadataData!.getStationMetadata}
                </CodeBlock>
            )}
        </div>
    );
};
