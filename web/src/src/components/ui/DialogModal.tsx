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
                className="fixed inset-0 z-40 flex h-dvh items-center justify-center bg-black/40 p-4"
                role="dialog"
                aria-modal="true"
                onMouseDown={({ currentTarget, target }) => {
                    if (currentTarget === target) {
                        onClose?.();
                    }
                }}
            >
                <div
                    className={`relative min-h-0 overflow-y-auto rounded-lg bg-white p-6 pb-[calc(1.5rem+env(safe-area-inset-bottom))] shadow-2xl ${fullScreen ? 'h-full max-h-full w-full max-w-none' : 'max-h-[90dvh] w-[90%] max-w-2xl sm:w-[80%] md:w-[60%]'}`}
                >
                    <button
                        className="btn btn-sm btn-circle btn-ghost absolute top-2 right-2"
                        onClick={onClose}
                    >
                        <Icon className="flex-shrink-0" path={mdiClose} size={0.8} />
                    </button>
                    {heading}
                    <div className="pt-4">{children}</div>
                </div>
            </div>
        );
    }

    return (
        <dialog ref={dialogRef} className="modal h-dvh">
            <div
                className={`modal-box min-h-0 overflow-y-auto pb-[calc(1.5rem+env(safe-area-inset-bottom))] ${fullScreen ? 'h-full max-h-full w-full max-w-none' : 'max-h-[90dvh] w-[90%] max-w-2xl sm:w-[80%] md:w-[60%]'}`}
            >
                <form method="dialog">
                    <button className="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">
                        <Icon className="flex-shrink-0" path={mdiClose} size={0.8} />
                    </button>
                </form>
                {heading}
                <div className="pt-4">{children}</div>
            </div>
        </dialog>
    );
};
