import { useEffect, useState } from 'react';

import { DownloadConstraints } from '../../config/constraints';
import { useGetFileListDataQuery } from '../../graphql';

export const useFileLists = () => {
    const { loading: getFileListDataLoading, data: getFileListDataData } = useGetFileListDataQuery({
        pollInterval: DownloadConstraints.pollInterval
    });

    const parseMseedFileName = (filename: string) => {
        const match = filename.match(
            /^([A-Z0-9]+)\.([A-Z0-9]+)\.([A-Z0-9]+)\.([A-Z0-9]+)\.D\.(\d{4})\.(\d{3})\.mseed$/i
        );
        if (!match) {
            return null;
        }
        return {
            channelCode: match[4],
            year: parseInt(match[5]),
            dayOfYear: parseInt(match[6])
        };
    };
    const [miniSeedList, setMiniSeedList] = useState<
        Array<{
            id: number;
            namespace: string;
            filePath: string;
            fileName: string;
            size: number;
            timestamp: number;
        }>
    >([]);
    useEffect(() => {
        if (getFileListDataData?.getMiniSeedFiles) {
            const rawList = [...getFileListDataData.getMiniSeedFiles]
                .sort((a, b) => {
                    const fa = parseMseedFileName(a!.fileName);
                    const fb = parseMseedFileName(b!.fileName);
                    if (!fa || !fb) {
                        return 0;
                    }
                    return (
                        fb.year - fa.year ||
                        fb.dayOfYear - fa.dayOfYear ||
                        fa.channelCode.localeCompare(fb.channelCode)
                    );
                })
                .map((file, index) => ({
                    id: index + 1,
                    namespace: file!.namespace,
                    filePath: file!.filePath,
                    fileName: file!.fileName,
                    size: file!.size,
                    timestamp: file!.modifiedAt
                }));
            setMiniSeedList(rawList);
        }
    }, [getFileListDataData]);

    const parseHelicorderFilename = (filename: string) => {
        const match = filename.match(
            /^([A-Z0-9]+)\.([A-Z0-9]+)\.([A-Z0-9]+)\.([A-Z0-9]+)\.(\d{4})\.(\d{3})\.svg$/i
        );
        if (!match) {
            return null;
        }
        return {
            channelCode: match[4],
            year: parseInt(match[5]),
            dayOfYear: parseInt(match[6])
        };
    };
    const [helicorderList, setHelicorderList] = useState<
        Array<{
            id: number;
            namespace: string;
            filePath: string;
            fileName: string;
            size: number;
            timestamp: number;
        }>
    >([]);
    useEffect(() => {
        if (getFileListDataData?.getHelicorderFiles) {
            const rawList = [...getFileListDataData.getHelicorderFiles]
                .sort((a, b) => {
                    const fa = parseHelicorderFilename(a!.fileName);
                    const fb = parseHelicorderFilename(b!.fileName);
                    if (!fa || !fb) {
                        return 0;
                    }
                    return (
                        fb.year - fa.year ||
                        fb.dayOfYear - fa.dayOfYear ||
                        fa.channelCode.localeCompare(fb.channelCode)
                    );
                })
                .map((file, index) => ({
                    id: index + 1,
                    namespace: file!.namespace,
                    filePath: file!.filePath,
                    fileName: file!.fileName,
                    size: file!.size,
                    timestamp: file!.modifiedAt
                }));
            setHelicorderList(rawList);
        }
    }, [getFileListDataData]);

    return { getFileListDataLoading, miniSeedList, helicorderList };
};
