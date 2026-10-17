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

const levelIcons = {
    info: mdiInformationOutline,
    success: mdiCheckCircleOutline,
    warning: mdiAlertOutline,
    error: mdiAlertCircleOutline
} satisfies Record<ServiceNotificationLevel, string>;

export const sendServiceNotification = (
    level: ServiceNotificationLevel,
    serviceId: string,
    message: string,
    duration = 6000
) =>
    toast.custom(
        ({ id, visible }) => (
            <div
                role="status"
                data-level={level}
                className={`app-notification ${
                    visible ? 'app-notification-visible' : 'app-notification-hidden'
                }`}
            >
                <div className="app-notification-accent" />
                <div className="app-notification-content">
                    <div className="app-notification-icon">
                        <Icon path={levelIcons[level]} size={0.9} />
                    </div>
                    <div className="app-notification-text">
                        <p className="app-notification-source">{serviceId}</p>
                        <p className="app-notification-message">{message}</p>
                    </div>
                    <button
                        type="button"
                        aria-label="Dismiss notification"
                        className="app-notification-close"
                        onClick={() => toast.dismiss(id)}
                    >
                        <Icon path={mdiClose} size={0.75} />
                    </button>
                </div>
            </div>
        ),
        { duration, position: 'top-right' }
    );
