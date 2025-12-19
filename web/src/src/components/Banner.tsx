import { useMemo } from 'react';

interface IBanner {
    readonly status: 'ok' | 'error' | 'warning';
    readonly message: string;
}

export const Banner = ({ status, message }: IBanner) => {
    const statusConfig = useMemo(
        () => ({
            ok: {
                alertClass: 'alert-success',
                statusClass: 'status-success'
            },
            error: {
                alertClass: 'alert-error',
                statusClass: 'status-error'
            },
            warning: {
                alertClass: 'alert-warning',
                statusClass: 'status-warning'
            }
        }),
        []
    );

    return (
        <div
            role="alert"
            className={`alert ${statusConfig[status].alertClass} alert-outline flex w-fit max-w-full`}
        >
            <div className={`status ${statusConfig[status].statusClass} animate-bounce`} />
            <span>{message}</span>
        </div>
    );
};
