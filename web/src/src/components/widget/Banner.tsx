import { useMemo } from 'react';

interface IBanner {
    readonly status: 'ok' | 'error' | 'warning';
    readonly message: string;
    readonly fullWidth?: boolean;
}

export const Banner = ({ status, message, fullWidth = false }: IBanner) => {
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
            className={`alert ${statusConfig[status].alertClass} alert-outline app-banner ${fullWidth ? 'app-full-width' : 'app-fit-width'} app-max-width`}
        >
            <div className={`status ${statusConfig[status].statusClass} app-bounce`} />
            <span>{message}</span>
        </div>
    );
};
