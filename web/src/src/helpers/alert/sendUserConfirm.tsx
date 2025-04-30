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
            <div
                className={`animate-fade animate-duration-300 flex w-full max-w-md rounded-lg bg-white shadow-xl ${
                    visible ? 'block' : 'hidden'
                }`}
            >
                <div className="flex-1 p-4">
                    <div className="flex items-start">
                        <div className="ml-3 flex-1">
                            <p className="text-sm font-medium text-gray-900">{title}</p>
                            <p className="mt-1 text-sm text-gray-500">{message}</p>
                        </div>
                    </div>
                </div>
                <div className="m-2 flex items-center justify-center gap-2">
                    <button
                        onClick={() => {
                            toast.dismiss(id);
                            onConfirmed();
                        }}
                        className="btn btn-md rounded-lg bg-red-500 px-3 py-2 font-medium text-white transition-all hover:bg-red-700"
                    >
                        {confirmBtnText}
                    </button>
                    <button
                        onClick={() => {
                            toast.dismiss(id);
                            onCancelled?.();
                        }}
                        className="btn btn-md rounded-lg px-3 py-2 font-medium transition-all hover:bg-gray-300"
                    >
                        {cancelBtnText}
                    </button>
                </div>
            </div>
        ),
        { duration: timeout }
    );
