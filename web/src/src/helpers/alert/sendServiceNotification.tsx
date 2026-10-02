import {
    mdiAlertCircleOutline,
    mdiAlertOutline,
    mdiCheckCircleOutline,
    mdiClose,
    mdiInformationOutline
} from '@mdi/js';
import Icon from '@mdi/react';
import toast from 'react-hot-toast';

export type ServiceNotificationLevel = 'info' | 'success' | 'warning' | 'error';

const levelStyles = {
    info: {
        accent: 'bg-blue-500',
        icon: mdiInformationOutline,
        iconBackground: 'bg-blue-50',
        iconColor: 'text-blue-600'
    },
    success: {
        accent: 'bg-emerald-500',
        icon: mdiCheckCircleOutline,
        iconBackground: 'bg-emerald-50',
        iconColor: 'text-emerald-600'
    },
    warning: {
        accent: 'bg-amber-500',
        icon: mdiAlertOutline,
        iconBackground: 'bg-amber-50',
        iconColor: 'text-amber-600'
    },
    error: {
        accent: 'bg-red-500',
        icon: mdiAlertCircleOutline,
        iconBackground: 'bg-red-50',
        iconColor: 'text-red-600'
    }
} satisfies Record<
    ServiceNotificationLevel,
    { accent: string; icon: string; iconBackground: string; iconColor: string }
>;

export const sendServiceNotification = (
    level: ServiceNotificationLevel,
    serviceId: string,
    message: string,
    duration = 6000
) =>
    toast.custom(
        ({ id, visible }) => {
            const style = levelStyles[level];

            return (
                <div
                    role="status"
                    className={`flex w-[calc(100vw-2rem)] max-w-sm overflow-hidden rounded-lg border border-gray-200 bg-white shadow-xl transition-all duration-200 ${
                        visible ? 'translate-y-0 opacity-100' : '-translate-y-2 opacity-0'
                    }`}
                >
                    <div className={`w-1.5 shrink-0 ${style.accent}`} />
                    <div className="flex min-w-0 flex-1 items-start gap-3 p-4">
                        <div
                            className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full ${style.iconBackground}`}
                        >
                            <Icon className={style.iconColor} path={style.icon} size={0.9} />
                        </div>
                        <div className="min-w-0 flex-1">
                            <p className="truncate text-xs font-semibold tracking-wide text-gray-500 uppercase">
                                {serviceId}
                            </p>
                            <p className="mt-1 text-sm leading-5 break-words text-gray-800">
                                {message}
                            </p>
                        </div>
                        <button
                            type="button"
                            aria-label="Dismiss notification"
                            className="shrink-0 rounded-md p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700"
                            onClick={() => toast.dismiss(id)}
                        >
                            <Icon path={mdiClose} size={0.75} />
                        </button>
                    </div>
                </div>
            );
        },
        { duration, position: 'top-right' }
    );
