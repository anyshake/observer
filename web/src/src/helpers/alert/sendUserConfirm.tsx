import toast from 'react-hot-toast';

interface Options {
    title: string;
    cancelBtnText: string;
    confirmBtnText: string;
    onConfirmed: () => void;
    onCancelled?: () => void;
    timeout?: number;
}

export const sendUserConfirm = (
    message: string,
    { title, cancelBtnText, confirmBtnText, onConfirmed, onCancelled, timeout = 30000 }: Options
) =>
    toast.custom(
        ({ visible, id }) => (
            <div className={`app-confirm ${visible ? 'app-block' : 'app-hidden'}`}>
                <div className="app-confirm-content">
                    <div className="app-confirm-row">
                        <div className="app-confirm-text">
                            <p className="app-confirm-title">{title}</p>
                            <p className="app-confirm-message">{message}</p>
                        </div>
                    </div>
                </div>
                <div className="app-confirm-actions">
                    <button
                        onClick={() => {
                            toast.dismiss(id);
                            onConfirmed();
                        }}
                        className="btn btn-md app-confirm-danger"
                    >
                        {confirmBtnText}
                    </button>
                    <button
                        onClick={() => {
                            toast.dismiss(id);
                            onCancelled?.();
                        }}
                        className="btn btn-md app-confirm-cancel"
                    >
                        {cancelBtnText}
                    </button>
                </div>
            </div>
        ),
        { duration: timeout }
    );
