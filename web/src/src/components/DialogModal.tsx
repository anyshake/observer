import { mdiClose } from '@mdi/js';
import Icon from '@mdi/react';
import { ReactNode, useEffect, useRef } from 'react';

interface IDialogModal {
    readonly open: boolean;
    readonly fullScreen?: boolean;
    readonly enlarge?: boolean;
    readonly heading?: ReactNode | ReactNode[];
    readonly children: ReactNode | ReactNode[];
    readonly onClose?: () => void;
}

export const DialogModal = ({
    open,
    fullScreen,
    enlarge,
    onClose,
    heading,
    children
}: IDialogModal) => {
    const dialogRef = useRef<HTMLDialogElement>(null);

    useEffect(() => {
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
    }, [open]);

    useEffect(() => {
        const dialog = dialogRef.current;
        if (!dialog || !onClose) {
            return;
        }

        const handleClose = () => {
            onClose();
        };

        dialog.addEventListener('close', handleClose);
        return () => dialog.removeEventListener('close', handleClose);
    }, [onClose]);

    return (
        <dialog ref={dialogRef} className="modal">
            <div
                className={`modal-box ${fullScreen ? 'h-screen w-full max-w-none' : enlarge ? 'w-[90%] max-w-none sm:w-[80%] md:w-[60%]' : ''}`}
            >
                <form method="dialog">
                    <button className="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">
                        <Icon className="flex-shrink-0" path={mdiClose} size={0.8} />
                    </button>
                </form>
                {heading}
                <div className="py-4">{children}</div>
            </div>
        </dialog>
    );
};
