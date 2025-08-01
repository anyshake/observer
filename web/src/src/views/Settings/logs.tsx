import { useEffect, useState } from 'react';

import { CodeBlock } from '../../components/CodeBlock';
import { ErrorPage } from '../../components/ErrorPage';
import { useGetApplicationLogsQuery } from '../../graphql';

export const Logs = () => {
    const {
        data: getApplicationLogsData,
        loading: getApplicationLogsLoading,
        error: getApplicationLogsError
    } = useGetApplicationLogsQuery({ pollInterval: 5000 });

    const [applicationLogs, setApplicationLogs] = useState<string>('');
    useEffect(() => {
        if (getApplicationLogsData?.getApplicationLogs) {
            setApplicationLogs(
                getApplicationLogsData?.getApplicationLogs
                    ?.map((line) => {
                        const logObj = JSON.parse(line);
                        return `${logObj.time} - [${logObj.level}] - [${logObj.module}] - ${logObj.msg}`;
                    })
                    .reverse()
                    .join('\n')
            );
        }
    }, [getApplicationLogsData?.getApplicationLogs]);

    return getApplicationLogsError ? (
        <ErrorPage
            content={getApplicationLogsError.message}
            debug={JSON.stringify(getApplicationLogsError)}
        />
    ) : (
        <div className="mx-auto max-w-3xl space-y-4">
            {getApplicationLogsLoading || !applicationLogs.length ? (
                <span className="loading loading-spinner text-primary" />
            ) : (
                <CodeBlock fileName="logs.txt" language="go" showLineNumbers>
                    {applicationLogs}
                </CodeBlock>
            )}
        </div>
    );
};
