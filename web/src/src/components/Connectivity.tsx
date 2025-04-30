import { mdiAlert, mdiAlertCircle, mdiCheckCircle } from '@mdi/js';
import { useMemo } from 'react';

interface IConnectivity {
    readonly status: 'ok' | 'error' | 'warning';
    readonly message: string;
}

export const Connectivity = ({ status, message }: IConnectivity) => {
    const statusConfig = useMemo(
        () => ({
            ok: {
                alertClass: 'alert-success',
                statusClass: 'status-success',
                icon: mdiCheckCircle,
                iconColor: 'text-success'
            },
            error: {
                alertClass: 'alert-error',
                statusClass: 'status-error',
                icon: mdiAlertCircle,
                iconColor: 'text-error'
            },
            warning: {
                alertClass: 'alert-warning',
                statusClass: 'status-warning',
                icon: mdiAlert,
                iconColor: 'text-warning'
            }
        }),
        []
    );

    return (
        <div
            role="alert"
            className={`alert ${statusConfig[status].alertClass} alert-outline flex max-w-2xl`}
        >
            <div className={`status ${statusConfig[status].statusClass} animate-bounce`} />
            <span>{message}</span>
        </div>
    );
};
