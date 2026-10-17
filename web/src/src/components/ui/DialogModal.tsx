import { mdiClose } from '@mdi/js';
import Icon from '@mdi/react';
import { ReactNode, useEffect, useRef } from 'react';

interface IDialogModal {
    readonly open: boolean;
    readonly nativeModal?: boolean;
    readonly fullScreen?: boolean;
    readonly heading?: ReactNode | ReactNode[];
    readonly children: ReactNode | ReactNode[];
    readonly onClose?: () => void;
}

export const DialogModal = ({
    open,
    nativeModal = true,
    fullScreen,
    onClose,
    heading,
    children
}: IDialogModal) => {
    const dialogRef = useRef<HTMLDialogElement>(null);

    useEffect(() => {
        if (!nativeModal) {
            return;
        }

        const dialog = dialogRef.current;
        if (!dialog) {
            return;
        }

        if (open) {
            if (!dialog.open) {
                dialog.showModal();
            }
        } else {
            if (dialog.open) {
                dialog.close();
            }
        }
    }, [nativeModal, open]);

    useEffect(() => {
        if (!nativeModal) {
            return;
        }

        const dialog = dialogRef.current;
        if (!dialog || !onClose) {
            return;
        }

        const handleClose = () => {
            onClose();
        };

        dialog.addEventListener('close', handleClose);
        return () => dialog.removeEventListener('close', handleClose);
    }, [nativeModal, onClose]);

    useEffect(() => {
        if (nativeModal || !open || !onClose) {
            return;
        }

        const handleKeyDown = ({ key }: KeyboardEvent) => {
            if (key === 'Escape') {
                onClose();
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [nativeModal, onClose, open]);

    if (!nativeModal) {
        if (!open) {
            return null;
        }

        return (
            <div
                className="app-dialog-overlay"
                role="dialog"
                aria-modal="true"
                onMouseDown={({ currentTarget, target }) => {
                    if (currentTarget === target) {
                        onClose?.();
                    }
                }}
            >
                <div
                    className={`app-dialog-panel ${fullScreen ? 'app-dialog-fullscreen' : 'app-dialog-window'}`}
                >
                    <button
                        className="btn btn-sm btn-circle btn-ghost app-dialog-close"
                        onClick={onClose}
                    >
                        <Icon className="app-icon" path={mdiClose} size={0.8} />
                    </button>
                    {heading}
                    <div className="app-dialog-body">{children}</div>
                </div>
            </div>
        );
    }

    return (
        <dialog ref={dialogRef} className="modal app-native-dialog">
            <div
                className={`modal-box app-native-dialog-panel ${fullScreen ? 'app-dialog-fullscreen' : 'app-dialog-window'}`}
            >
                <form method="dialog">
                    <button className="btn btn-sm btn-circle btn-ghost app-dialog-close">
                        <Icon className="app-icon" path={mdiClose} size={0.8} />
                    </button>
                </form>
                {heading}
                <div className="app-dialog-body">{children}</div>
            </div>
        </dialog>
    );
};
